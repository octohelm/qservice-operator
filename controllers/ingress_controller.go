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
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// IngressReconciler reconciles a Ingress object
type IngressReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionsv1beta1.Ingress{}).
		Owns(&istiov1beta1.VirtualService{}).
		Complete(r)
}

func (r *IngressReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	log := r.Log.WithValues("namespace", request.Namespace, "name", request.Name)

	ingress := &extensionsv1beta1.Ingress{}
	if err := r.Client.Get(ctx, request.NamespacedName, ingress); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	ctx = controllerutil.ContextWithControllerClient(ctx, r.Client)

	_ = r.patchIngressIfNeed(ctx, ingress)

	ok := r.isClusterWithIstio(ctx)
	if ok {
		if err := r.applyVirtualService(ctx, ingress); err != nil {
			log.Error(err, "apply virtual service failed")
			return reconcile.Result{}, nil
		}
	}

	return reconcile.Result{}, nil
}

func (r *IngressReconciler) patchIngressIfNeed(ctx context.Context, ingress *extensionsv1beta1.Ingress) error {
	needApplyPatch := false

	if ingress.Labels == nil {
		ingress.Labels = map[string]string{}
	}

	if len(ingress.Spec.Rules) > 0 {
		if _, ok := ingress.Labels[LabelBashHost]; !ok {
			ingress.Labels[LabelBashHost] = BaseHost(ingress.Spec.Rules[0].Host)
			needApplyPatch = true
		}

		backend := ingress.Spec.Rules[0].HTTP
		if len(backend.Paths) > 0 {
			expectServiceName := backend.Paths[0].Backend.ServiceName

			if serviceName := ingress.Labels[LabelServiceName]; serviceName != expectServiceName {
				ingress.Labels[LabelServiceName] = expectServiceName
				needApplyPatch = true
			}
		}
	}

	if needApplyPatch {
		if err := applyIngress(ctx, ingress); err != nil {
			return err
		}
	}

	return nil
}

func (r *IngressReconciler) applyVirtualService(ctx context.Context, ingress *extensionsv1beta1.Ingress) error {
	vss := converter.ToExportedVirtualServicesByIngress(ingress)
	for i := range vss {
		vs := vss[i]
		r.setControllerReference(vs, ingress)

		if err := applyVirtualService(ctx, vs); err != nil {
			return err
		}
	}
	return nil
}

func (r *IngressReconciler) setControllerReference(obj metav1.Object, owner metav1.Object) {
	if err := controllerutil.SetControllerReference(owner, obj, r.Scheme); err != nil {
		r.Log.Error(err, "")
	}
	obj.SetAnnotations(controllerutil.AnnotateControllerGeneration(obj.GetAnnotations(), owner.GetGeneration()))
}

func (r *IngressReconciler) isClusterWithIstio(ctx context.Context) bool {
	n := &corev1.Namespace{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: "istio-system", Namespace: ""}, n)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false
		}
		r.Log.Error(err, "")
		return false
	}
	return true
}
