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
	"fmt"

	"github.com/go-logr/logr"
	"github.com/octohelm/qservice-operator/apis/serving/v1alpha1"
	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// QIngressReconciler reconciles a Ingress object
type QIngressReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *QIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.QIngress{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}

func (r *QIngressReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.Log.WithValues("namespace", request.Namespace, "name", request.Name)

	qingress := &v1alpha1.QIngress{}
	if err := r.Client.Get(ctx, request.NamespacedName, qingress); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	ctx = controllerutil.ContextWithControllerClient(ctx, r.Client)

	hostname, ok := IngressGateways.IngressGatewayHost(qingress.Spec.Ingress.Host)
	if !ok {
		return reconcile.Result{}, fmt.Errorf("invalid gateway of %s", qingress.Spec.Ingress.Host)
	}

	if err := r.applyIngress(ctx, qingress, hostname); err != nil {
		log.Error(err, "apply ingress failed")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *QIngressReconciler) applyIngress(ctx context.Context, qingress *v1alpha1.QIngress, hostname string) error {
	ing := toIngressByQIngress(qingress, hostname)
	r.setControllerReference(ing, qingress)
	if err := applyIngress(ctx, ing); err != nil {
		return err
	}
	return nil
}

func (r *QIngressReconciler) setControllerReference(obj metav1.Object, owner metav1.Object) {
	if err := controllerutil.SetControllerReference(owner, obj, r.Scheme); err != nil {
		r.Log.Error(err, "")
	}
	obj.SetAnnotations(controllerutil.AnnotateControllerGeneration(obj.GetAnnotations(), owner.GetGeneration()))
}

func toIngressByQIngress(qingress *v1alpha1.QIngress, hostname string) *networkingv1.Ingress {
	ing := &networkingv1.Ingress{}

	ing.Namespace = qingress.Namespace
	ing.Name = qingress.Name
	ing.Annotations = qingress.Annotations

	paths := make([]networkingv1.HTTPIngressPath, 0)

	if len(qingress.Spec.Ingress.Paths) > 0 {
		for i := range qingress.Spec.Ingress.Paths {
			p := qingress.Spec.Ingress.Paths[i]

			htp := networkingv1.HTTPIngressPath{}
			htp.Path = p.Path

			if p.Exactly {
				pt := networkingv1.PathTypeExact
				htp.PathType = &pt
			} else {
				pt := networkingv1.PathTypeImplementationSpecific
				htp.PathType = &pt
			}

			htp.Backend = qingress.Spec.Backend
			paths = append(paths, htp)
		}
	} else {
		htp := networkingv1.HTTPIngressPath{}

		pt := networkingv1.PathTypeImplementationSpecific
		htp.PathType = &pt

		htp.Backend = qingress.Spec.Backend
		paths = append(paths, htp)
	}

	ing.Spec.Rules = []networkingv1.IngressRule{
		{
			Host: hostname,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: paths,
				},
			},
		},
	}

	return ing
}
