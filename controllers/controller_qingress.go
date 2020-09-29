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
	"github.com/octohelm/qservice-operator/pkg/converter"
	istiotypes "istio.io/api/networking/v1alpha3"
	istiov1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
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

// QIngressReconciler reconciles a Ingress object
type QIngressReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *QIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.QIngress{}).
		Owns(&extensionsv1beta1.Ingress{}).
		Owns(&istiov1alpha3.VirtualService{}).
		Complete(r)
}

func (r *QIngressReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

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

	if ok := r.isClusterWithIstio(ctx); ok {
		if err := r.applyVirtualService(ctx, qingress, hostname); err != nil {
			log.Error(err, "apply virtual service failed")
			return reconcile.Result{}, err
		}
	} else {
		if err := r.applyIngress(ctx, qingress, hostname); err != nil {
			log.Error(err, "apply ingress failed")
			return reconcile.Result{}, err
		}
	}

	// todo added status

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

func (r *QIngressReconciler) applyVirtualService(ctx context.Context, qingress *v1alpha1.QIngress, hostname string) error {
	vs := toExportedVirtualServicesByQIngress(qingress, hostname)
	r.setControllerReference(vs, qingress)
	if err := applyVirtualService(ctx, vs); err != nil {
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

func (r *QIngressReconciler) isClusterWithIstio(ctx context.Context) bool {
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

func toIngressByQIngress(qingress *v1alpha1.QIngress, hostname string) *extensionsv1beta1.Ingress {
	ing := &extensionsv1beta1.Ingress{}

	ing.Namespace = qingress.Namespace
	ing.Name = qingress.Name
	ing.Annotations = map[string]string{
		"kubernetes.io/ingress.class": "nginx",
	}

	paths := make([]extensionsv1beta1.HTTPIngressPath, 0)

	if len(qingress.Spec.Ingress.Paths) > 0 {
		for i := range qingress.Spec.Ingress.Paths {
			p := qingress.Spec.Ingress.Paths[i]

			htp := extensionsv1beta1.HTTPIngressPath{}
			htp.Path = p.Path

			if p.Exactly {
				pt := extensionsv1beta1.PathTypeExact
				htp.PathType = &pt
			}

			htp.Backend = qingress.Spec.Backend
			paths = append(paths, htp)
		}
	} else {
		htp := extensionsv1beta1.HTTPIngressPath{}

		htp.Backend = qingress.Spec.Backend
		paths = append(paths, htp)
	}

	ing.Spec.Rules = []extensionsv1beta1.IngressRule{
		{
			Host: hostname,
			IngressRuleValue: extensionsv1beta1.IngressRuleValue{
				HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
					Paths: paths,
				},
			},
		},
	}

	return ing
}

func toExportedVirtualServicesByQIngress(qingress *v1alpha1.QIngress, hostname string) *istiov1alpha3.VirtualService {
	vs := &istiov1alpha3.VirtualService{}
	vs.Namespace = qingress.Namespace
	vs.Name = qingress.Name

	vs.Spec.Hosts = []string{hostname}

	gatewayName := converter.ConvertToGatewayName(hostname)

	if gatewayName != "" {
		vs.Spec.Gateways = append(vs.Spec.Gateways, gatewayName)
	}

	route := &istiotypes.HTTPRoute{
		Route: []*istiotypes.HTTPRouteDestination{
			{
				Destination: &istiotypes.Destination{
					Host: qingress.Spec.Backend.ServiceName,
				},
			},
		},
	}

	for j := range qingress.Spec.Ingress.Paths {
		p := qingress.Spec.Ingress.Paths[j]

		if p.Path != "" {
			if p.Exactly {
				route.Match = append(route.Match, &istiotypes.HTTPMatchRequest{
					Uri: &istiotypes.StringMatch{
						MatchType: &istiotypes.StringMatch_Exact{Exact: p.Path},
					},
				})
			} else {
				route.Match = append(route.Match, &istiotypes.HTTPMatchRequest{
					Uri: &istiotypes.StringMatch{
						MatchType: &istiotypes.StringMatch_Prefix{Prefix: p.Path},
					},
				})
			}
		}
	}

	vs.Spec.Http = append(vs.Spec.Http, route)

	return vs
}
