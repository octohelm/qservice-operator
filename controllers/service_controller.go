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

	"github.com/go-logr/logr"
	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	"github.com/octohelm/qservice-operator/pkg/converter"
	istiov1beta1 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Service{}).
		Owns(&istiov1beta1.VirtualService{}).
		Complete(r)
}

func (r *ServiceReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	ok, err := r.istioInjectionEnabledInNamespace(ctx, request.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}
	if !ok {
		return reconcile.Result{}, nil
	}

	s := &corev1.Service{}
	if err := r.Client.Get(ctx, request.NamespacedName, s); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	ctx = controllerutil.ContextWithControllerClient(ctx, r.Client)

	vs := converter.ToClusterVirtualServiceFromService(s)
	r.setControllerReference(vs, s)

	if err := applyVirtualService(ctx, s.Namespace, vs); err != nil {
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}

func (r *ServiceReconciler) setControllerReference(obj metav1.Object, owner metav1.Object) {
	if err := controllerutil.SetControllerReference(owner, obj, r.Scheme); err != nil {
		r.Log.Error(err, "")
	}
	obj.SetAnnotations(controllerutil.AnnotateControllerGeneration(obj.GetAnnotations(), owner.GetGeneration()))
}

func (r *ServiceReconciler) istioInjectionEnabledInNamespace(ctx context.Context, namespace string) (bool, error) {
	n := &corev1.Namespace{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: namespace, Namespace: ""}, n)
	if err != nil {
		return false, err
	}
	return n.Labels["istio-injection"] == "enabled", nil
}
