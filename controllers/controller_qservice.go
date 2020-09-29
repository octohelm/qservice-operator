package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/go-logr/logr"
	servingv1alpha1 "github.com/octohelm/qservice-operator/apis/serving/v1alpha1"
	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	"github.com/octohelm/qservice-operator/pkg/converter"
	"github.com/octohelm/qservice-operator/pkg/strfmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// QServiceReconciler reconciles a QService object
type QServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *QServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&servingv1alpha1.QService{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Watches(&source.Kind{Type: &servingv1alpha1.QIngress{}}, &handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(object handler.MapObject) []reconcile.Request {
				// to trigger sync qingresses as QService status
				if labelSet := object.Meta.GetLabels(); labelSet != nil {
					if app, ok := labelSet[LabelServiceName]; ok {
						return []reconcile.Request{{NamespacedName: types.NamespacedName{
							Name:      app,
							Namespace: object.Meta.GetNamespace(),
						}}}
					}
				}
				return []reconcile.Request{}
			}),
		}).
		Complete(r)
}

// Reconcile reads that state of the cluster for a QService object and makes changes based on the state read
// and what is in the QService.Spec
func (r *QServiceReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	log := r.Log.WithValues("namespace", request.Namespace, "name", request.Name)

	qsvc := &servingv1alpha1.QService{}

	err := r.Client.Get(ctx, request.NamespacedName, qsvc)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		log.Error(err, "Failed to get QService")
		return reconcile.Result{}, err
	}

	if err := r.applyQService(ctx, qsvc); err != nil {
		log.Error(err, "apply failed")
		return reconcile.Result{}, err
	}

	_ = r.updateStatusFromDeployment(ctx, qsvc)
	_ = r.updateStatusFromQIngresses(ctx, qsvc)

	if err := r.Client.Status().Update(ctx, qsvc); err != nil {
		log.Error(err, "update status failed")
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}

func (r *QServiceReconciler) updateStatusFromQIngresses(ctx context.Context, qsvc *servingv1alpha1.QService) error {
	list := &servingv1alpha1.QIngressList{}

	s, _ := labels.NewRequirement(LabelServiceName, selection.Equals, []string{qsvc.Name})

	if err := r.Client.List(ctx,
		list,
		client.InNamespace(qsvc.Namespace),
		client.MatchingLabelsSelector{Selector: labels.NewSelector().Add(*s)},
	); err != nil {
		return err
	}

	if len(list.Items) == 0 {
		return nil
	}

	qsvc.Status.Ingresses = map[string][]strfmt.Ingress{}

	for i := range list.Items {
		item := list.Items[i]

		host, ok := IngressGateways.IngressGatewayHost(item.Spec.Ingress.Host)
		if ok {
			item.Spec.Ingress.Host = host
		}

		qsvc.Status.Ingresses[item.Name] = append(qsvc.Status.Ingresses[item.Name], item.Spec.Ingress)
	}

	return nil
}

func (r *QServiceReconciler) updateStatusFromDeployment(ctx context.Context, qsvc *servingv1alpha1.QService) error {
	deployment := &appsv1.Deployment{}

	if err := r.Client.Get(ctx, types.NamespacedName{Name: qsvc.Name, Namespace: qsvc.Namespace}, deployment); err != nil {
		return err
	}

	if reflect.DeepEqual(deployment.Status, qsvc.Status.DeploymentStage) {
		return nil
	}

	podList := &corev1.PodList{}

	if err := r.Client.List(
		ctx, podList,
		client.InNamespace(qsvc.Namespace),
		client.MatchingLabels(map[string]string{
			"app": qsvc.Name,
		}),
	); err != nil {
		return err
	}

	qsvc.Status.DeploymentStatus = deployment.Status
	qsvc.Status.DeploymentStage, qsvc.Status.DeploymentComments = toDeploymentStage(&deployment.Status, podList.Items)

	// update deployment status to trigger lifecycle for getting container status from pod list
	if qsvc.Status.DeploymentStage == "PROCESSING" {
		for i := range deployment.Status.Conditions {
			c := deployment.Status.Conditions[i]

			idx := i

			if c.Type == "Progressing" && c.Reason != "NewReplicaSetAvailable" {
				go func() {
					interval := 5 * time.Second

					time.Sleep(interval)

					deployment.Status.Conditions[idx].LastUpdateTime = metav1.Time{
						Time: deployment.Status.Conditions[idx].LastUpdateTime.Add(interval),
					}

					err := r.Client.Status().Update(ctx, deployment)
					if err != nil {
						if !apierrors.IsConflict(err) {
							r.Log.Error(err, "update deployment status failed")
						}
					}
				}()
			}
		}
	}

	return nil
}

func toDeploymentStage(status *appsv1.DeploymentStatus, pods []corev1.Pod) (string, string) {
	if status.UnavailableReplicas == 0 && status.AvailableReplicas == status.Replicas {
		return "DONE", ""
	}

	for _, c := range status.Conditions {
		if c.Type == appsv1.DeploymentReplicaFailure {
			return "FAILED", ""
		}
	}

	b := bytes.NewBuffer(nil)
	stage := "PROCESSING"

	for _, pod := range pods {
		podStatus := pod.Status

		if podStatus.Phase == corev1.PodPending || podStatus.Phase == corev1.PodFailed {
			for i := range podStatus.ContainerStatuses {
				containerStatus := podStatus.ContainerStatuses[i]
				if !containerStatus.Ready {
					if containerStatus.State.Waiting != nil && containerStatus.State.Waiting.Message != "" {

						if containerStatus.State.Waiting.Reason != "ContainerCreating" {
							stage = "FAILED"
						}

						_, _ = io.WriteString(b, fmt.Sprintf("[%s] %s", containerStatus.State.Waiting.Reason, containerStatus.State.Waiting.Message))
					}
				}
			}
		}
	}

	return stage, b.String()
}

func (r *QServiceReconciler) setControllerReference(obj metav1.Object, owner metav1.Object) {
	if err := controllerutil.SetControllerReference(owner, obj, r.Scheme); err != nil {
		r.Log.Error(err, "")
	}
	obj.SetAnnotations(controllerutil.AnnotateControllerGeneration(obj.GetAnnotations(), owner.GetGeneration()))
}

func (r *QServiceReconciler) applyQService(ctx context.Context, qsvc *servingv1alpha1.QService) error {
	qsvc.Labels["app"] = qsvc.Name

	ctx = controllerutil.ContextWithControllerClient(ctx, r.Client)

	qsvc.Status.ServiceConditions = []servingv1alpha1.QServiceCondition{}

	appendServiceConditions := func(tpe string, err error) {
		condition := servingv1alpha1.QServiceCondition{Type: tpe}
		if err == nil {
			condition.Status = corev1.ConditionTrue
		} else {
			condition.Status = corev1.ConditionFalse
			condition.Reason = err.Error()
		}
		qsvc.Status.ServiceConditions = append(qsvc.Status.ServiceConditions, condition)
	}

	if err := r.applyImagePullSecret(ctx, qsvc); err != nil {
		appendServiceConditions("ImagePullSecret", err)
		return err
	}
	appendServiceConditions("ImagePullSecret", nil)

	if err := r.applyDeployment(ctx, qsvc); err != nil {
		appendServiceConditions("Deployment", err)
		return err
	}
	appendServiceConditions("Deployment", nil)

	if err := r.applyService(ctx, qsvc); err != nil {
		appendServiceConditions("Service", err)
		return err
	}
	appendServiceConditions("Service", nil)

	return nil
}

func (r *QServiceReconciler) applyService(ctx context.Context, qsvc *servingv1alpha1.QService) error {
	s := converter.ToService(qsvc)
	r.setControllerReference(s, qsvc)

	if len(s.Spec.Ports) > 0 {
		if err := applyService(ctx, s); err != nil {
			return err
		}
	}

	return nil
}

const AnnotationImageKeyPullSecret = "serving.octohelm.tech/imagePullSecret"

func (r *QServiceReconciler) applyImagePullSecret(ctx context.Context, qsvc *servingv1alpha1.QService) error {
	if pullSecret, ok := qsvc.Annotations[AnnotationImageKeyPullSecret]; ok {
		ips, err := strfmt.ParseImagePullSecret(pullSecret)
		if err != nil {
			return err
		}

		// service scoped pull secret
		ips.Name = qsvc.Name + "--pull-secret"
		qsvc.Spec.Image = ips.PrefixTag(qsvc.Spec.Image)
		qsvc.Spec.ImagePullSecret = ips.SecretName()

		secret := converter.ToImagePullSecret(ips, qsvc.Namespace)
		r.setControllerReference(secret, qsvc)

		if err := applySecret(ctx, secret); err != nil {
			return err
		}
	}

	return nil
}

func (r *QServiceReconciler) applyDeployment(ctx context.Context, qsvc *servingv1alpha1.QService) error {
	deployment := converter.ToDeployment(qsvc)
	r.setControllerReference(deployment, qsvc)

	if err := applyDeployment(ctx, deployment); err != nil {
		return err
	}

	return nil
}
