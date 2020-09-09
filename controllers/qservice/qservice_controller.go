package qservice

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-logr/logr"
	servingv1alpha1 "github.com/octohelm/qservice-operator/apis/serving/v1alpha1"
	"github.com/octohelm/qservice-operator/pkg/converter"
	"github.com/octohelm/qservice-operator/pkg/strfmt"
	istiov1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// QServiceReconciler reconciles a QService object
type QServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *QServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	cm := ctrl.NewControllerManagedBy(mgr).
		For(&servingv1alpha1.QService{})

	objects := []runtime.Object{
		&appsv1.Deployment{},
		&autoscalingv2beta1.HorizontalPodAutoscaler{},
		&corev1.Service{},
		&extensionsv1beta1.Ingress{},
		&corev1.Secret{},
		&istiov1alpha3.VirtualService{},
	}

	for i := range objects {
		cm = cm.Watches(&source.Kind{Type: objects[i]}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &servingv1alpha1.QService{},
		})
	}

	return cm.Complete(r)
}

// Reconcile reads that state of the cluster for a QService object and makes changes based on the state read
// and what is in the QService.Spec
func (r *QServiceReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	qsvc := &servingv1alpha1.QService{}

	err := r.Client.Get(ctx, request.NamespacedName, qsvc)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info("QService resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		r.Log.Error(err, "Failed to get QService")
		return reconcile.Result{}, err
	}

	if err := r.applyQService(ctx, qsvc); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.updateDeploymentStage(ctx, qsvc); err != nil {
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}

func (r *QServiceReconciler) updateDeploymentStage(ctx context.Context, qsvc *servingv1alpha1.QService) error {
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

	err := r.Client.Status().Update(ctx, qsvc)
	if err != nil {
		return err
	}

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
	_ = controllerutil.SetControllerReference(owner, obj, r.Scheme)
	obj.SetAnnotations(annotateControllerGeneration(obj.GetAnnotations(), owner.GetGeneration()))
}

func (r *QServiceReconciler) getFlagsFromNamespace(ctx context.Context, namespace string) (Flags, error) {
	n := &corev1.Namespace{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: namespace, Namespace: ""}, n)
	if err != nil {
		return Flags{}, err
	}
	return FlagsFromNamespaceLabels(n.Labels), nil
}

func (r *QServiceReconciler) applyQService(ctx context.Context, qsvc *servingv1alpha1.QService) error {
	flags, err := r.getFlagsFromNamespace(ctx, qsvc.Namespace)
	if err != nil {
		return err
	}

	qsvc.Labels["app"] = qsvc.Name

	ctx = ContextWithControllerClient(ctx, r.Client)

	return with(
		r.applyImagePullSecret,
		r.applyDeployment,
		r.applyService,
		r.applyIngress,
	)(ctx, qsvc, &flags)
}

type process = func(ctx context.Context, qsvc *servingv1alpha1.QService, flags *Flags) error

func with(processes ...process) process {
	return func(ctx context.Context, qsvc *servingv1alpha1.QService, flags *Flags) error {
		for i := range processes {
			p := processes[i]
			if err := p(ctx, qsvc, flags); err != nil {
				return fmt.Errorf("process(%d) %s", i, err)
			}
		}
		return nil
	}
}

var autoIngressHosts = toAutoIngressHosts(os.Getenv("AUTO_INGRESS_HOSTS"))

func toAutoIngressHosts(v string) map[string]bool {
	if v == "" {
		return map[string]bool{}
	}

	m := map[string]bool{}

	for _, h := range strings.Split(v, ",") {
		m[h] = true
	}

	return m
}

func (r *QServiceReconciler) applyIngress(ctx context.Context, qsvc *servingv1alpha1.QService, flags *Flags) error {
	if len(autoIngressHosts) > 0 {
		m := map[string]bool{}

		for i := range qsvc.Spec.Ingresses {
			ingress := qsvc.Spec.Ingresses[i]
			m[ingress.Host] = true
		}

		for autoIngressHost := range autoIngressHosts {
			h := fmt.Sprintf("%s---%s.%s", qsvc.Name, qsvc.Namespace, autoIngressHost)

			if m[h] {
				continue
			}

			port := uint16(80)

			if len(qsvc.Spec.Ports) > 0 {
				port = qsvc.Spec.Ports[0].Port
			}

			qsvc.Spec.Ingresses = append(qsvc.Spec.Ingresses, strfmt.Ingress{
				Scheme: "http",
				Host:   h,
				Port:   port,
			})
		}
	}

	if len(qsvc.Spec.Ingresses) > 0 {
		ingress := converter.ToIngress(qsvc)
		r.setControllerReference(ingress, qsvc)

		if err := applyIngress(ctx, qsvc.Namespace, ingress); err != nil {
			return err
		}

		groupedIngresses := map[string][]strfmt.Ingress{}

		for _, h := range qsvc.Spec.Ingresses {
			groupedIngresses[h.Host] = append(groupedIngresses[h.Host], h)
		}

		hostIDs := make([]string, 0)

		for host := range groupedIngresses {
			vs := converter.ToExportedVirtualService(qsvc, host, groupedIngresses[host])
			r.setControllerReference(vs, qsvc)

			hostIDs = append(hostIDs, vs.Name)

			if err := applyVirtualService(ctx, qsvc.Namespace, vs); err != nil {
				return err
			}
		}

		go func() {
			ls, _ := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
				MatchLabels: map[string]string{
					"controlled-by": ingress.Name,
				},
				MatchExpressions: []metav1.LabelSelectorRequirement{{
					Key:      "host",
					Operator: metav1.LabelSelectorOpNotIn,
					Values:   hostIDs,
				}},
			})

			err := r.Client.DeleteAllOf(ctx, &istiov1alpha3.VirtualService{},
				client.InNamespace(qsvc.Namespace),
				client.MatchingLabelsSelector{
					Selector: ls,
				},
			)
			if err != nil {
				r.Log.Error(err, "cleanup failed")
			}
		}()
	}

	return nil
}

func (r *QServiceReconciler) applyService(ctx context.Context, qsvc *servingv1alpha1.QService, flags *Flags) error {
	s := converter.ToService(qsvc)
	r.setControllerReference(s, qsvc)

	if err := applyService(ctx, qsvc.Namespace, s); err != nil {
		return err
	}

	if flags.IstioEnabled {
		vs := converter.ToClusterVirtualService(qsvc)
		r.setControllerReference(vs, qsvc)

		if err := applyVirtualService(ctx, qsvc.Namespace, vs); err != nil {
			return err
		}
	}

	return nil
}

func (r *QServiceReconciler) applyImagePullSecret(ctx context.Context, qsvc *servingv1alpha1.QService, flags *Flags) error {
	if pullSecret, ok := qsvc.Annotations[AnnotationImageKeyPullSecret]; ok {
		ips, err := strfmt.ParseImagePullSecret(pullSecret)
		if err != nil {
			return err
		}

		// service scoped pull secret
		ips.Name = qsvc.Name + "--pull-secret"
		qsvc.Spec.Image = ips.PrefixTag(qsvc.Spec.Image)
		qsvc.Spec.ImagePullSecret = ips.SecretName()

		secret := converter.ToImagePullSecret(ips)
		r.setControllerReference(secret, qsvc)

		if err := applySecret(ctx, qsvc.Namespace, secret); err != nil {
			return err
		}
	}

	return nil
}

func (r *QServiceReconciler) applyDeployment(ctx context.Context, qsvc *servingv1alpha1.QService, flags *Flags) error {
	deployment := converter.ToDeployment(qsvc)
	r.setControllerReference(deployment, qsvc)

	if err := applyDeployment(ctx, qsvc.Namespace, deployment); err != nil {
		return err
	}

	return nil
}
