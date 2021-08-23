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

package controllers

import (
	"context"
	"strconv"
	"strings"

	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"

	"github.com/go-courier/x/ptr"
	"github.com/go-logr/logr"
	"github.com/octohelm/qservice-operator/pkg/strfmt"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Owns(&autoscalingv2beta1.HorizontalPodAutoscaler{}).
		Complete(r)
}

func (r *DeploymentReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	ok, err := r.autoScalingEnabledInNamespace(ctx, request.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}
	if !ok {
		return reconcile.Result{}, nil
	}

	dep := &appsv1.Deployment{}
	if err := r.Client.Get(ctx, request.NamespacedName, dep); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	hpa := toHorizontalPodAutoscaler(dep)
	if hpa != nil {
		r.setControllerReference(hpa, dep)

		if err := applyHorizontalPodAutoscaler(ctx, r.Client, dep.Namespace, hpa); err != nil {
			return reconcile.Result{}, err
		}
	} else {
		if err := deleteHorizontalPodAutoscaler(ctx, r.Client, request.NamespacedName); err != nil {
			return reconcile.Result{}, nil
		}
	}

	return reconcile.Result{}, nil
}

func (r *DeploymentReconciler) autoScalingEnabledInNamespace(ctx context.Context, namespace string) (bool, error) {
	n := &corev1.Namespace{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: namespace, Namespace: ""}, n)
	if err != nil {
		return false, err
	}
	return n.Labels["autoscaling"] == "enabled", nil
}

func (r *DeploymentReconciler) setControllerReference(obj metav1.Object, owner metav1.Object) {
	if err := controllerutil.SetControllerReference(owner, obj, r.Scheme); err != nil {
		r.Log.Error(err, "")
	}
	obj.SetAnnotations(controllerutil.AnnotateControllerGeneration(obj.GetAnnotations(), owner.GetGeneration()))
}

func deleteHorizontalPodAutoscaler(ctx context.Context, c client.Client, namespacedName types.NamespacedName) error {
	hpa := &autoscalingv2beta1.HorizontalPodAutoscaler{}
	hpa.Namespace = namespacedName.Namespace
	hpa.Name = namespacedName.Name

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

	return c.Patch(ctx, hpa, client.Merge)
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
		hpa.Namespace = dep.Namespace
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
