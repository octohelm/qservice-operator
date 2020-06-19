/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package deployment

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-courier/ptr"
	"github.com/octohelm/qservice-operator/pkg/apiutil"
	"github.com/octohelm/qservice-operator/pkg/strfmt"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_deployment")

// Add creates a new Deployment Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDeployment{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("deployment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Deployment
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileDeployment{}

type ReconcileDeployment struct {
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileDeployment) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Deployment")

	ctx := context.Background()

	dep := &appsv1.Deployment{}
	err := r.client.Get(ctx, request.NamespacedName, dep)
	if err != nil {
		if errors.IsNotFound(err) {
			if err := deleteHorizontalPodAutoscaler(ctx, r.client, request.NamespacedName); err != nil {
				return reconcile.Result{}, nil
			}
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	hpa := toHorizontalPodAutoscaler(dep)
	if hpa != nil {
		if err := applyHorizontalPodAutoscaler(ctx, r.client, dep.Namespace, hpa); err != nil {
			return reconcile.Result{}, err
		}
	} else {
		if err := deleteHorizontalPodAutoscaler(ctx, r.client, request.NamespacedName); err != nil {
			return reconcile.Result{}, nil
		}
	}

	return reconcile.Result{}, nil
}

func deleteHorizontalPodAutoscaler(ctx context.Context, c client.Client, namespacedName types.NamespacedName) error {
	hpa := &autoscalingv2beta1.HorizontalPodAutoscaler{}
	hpa.Name = namespacedName.Name
	hpa.Namespace = namespacedName.Namespace

	if err := c.Delete(ctx, hpa); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

func applyHorizontalPodAutoscaler(ctx context.Context, c client.Client, namespace string, hpa *autoscalingv2beta1.HorizontalPodAutoscaler) error {
	hpa.Namespace = namespace

	current := &autoscalingv2beta1.HorizontalPodAutoscaler{}

	err := c.Get(ctx, types.NamespacedName{Name: hpa.Name, Namespace: namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		return c.Create(ctx, hpa)
	}

	return c.Patch(ctx, hpa, apiutil.JSONPatch(types.MergePatchType))
}

var AnnotationKeyAutoScalingGroup = "autoscaling.octohelm.tech"

func toHorizontalPodAutoscaler(dep *appsv1.Deployment) *autoscalingv2beta1.HorizontalPodAutoscaler {
	annotations := dep.GetAnnotations()

	hasAuthScaling := false

	for key := range annotations {
		if strings.HasPrefix(key, AnnotationKeyAutoScalingGroup) {
			hasAuthScaling = true
			break
		}
	}

	if hasAuthScaling {
		hpa := &autoscalingv2beta1.HorizontalPodAutoscaler{}
		hpa.Name = dep.Name
		hpa.Labels = dep.Labels

		hpa.Spec.ScaleTargetRef.APIVersion = dep.APIVersion
		hpa.Spec.ScaleTargetRef.Kind = dep.Kind
		hpa.Spec.ScaleTargetRef.Name = dep.Name

		if v, ok := annotations[AnnotationKeyAutoScalingGroup+"/minScale"]; ok {
			i, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				hpa.Spec.MinReplicas = ptr.Int32(int32(i))
			}
		}

		if v, ok := annotations[AnnotationKeyAutoScalingGroup+"/maxScale"]; ok {
			i, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				hpa.Spec.MaxReplicas = int32(i)
			}
		}

		if hpa.Spec.MaxReplicas == 0 {
			hpa.Spec.MaxReplicas = 5
		}

		if v, ok := annotations[AnnotationKeyAutoScalingGroup+"/metrics"]; ok {
			metrics, err := strfmt.ParseMetrics(v)
			if err != nil {
				return nil
			}
			hpa.Spec.Metrics = metrics
		}

		return hpa
	}

	return nil
}
