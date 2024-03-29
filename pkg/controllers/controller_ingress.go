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
	istiotypes "istio.io/api/networking/v1alpha3"
	istiov1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	if !controllerutil.IsResourceRegistered(r.Client, istiov1alpha3.SchemeGroupVersion.WithKind("VirtualService")) {
		return nil
	}

	if controllerutil.IsResourceRegistered(r.Client, istiov1alpha3.SchemeGroupVersion.WithKind("Gateway")) {
		if err := controllerutil.ApplyGateways(mgr.GetConfig(), IngressGateways.ToGateways()...); err != nil {
			return err
		}
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.Ingress{}).
		Owns(&istiov1alpha3.VirtualService{}).
		Complete(r)
}

func (r *IngressReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.Log.WithValues("namespace", request.Namespace, "name", request.Name)

	ingress := &networkingv1.Ingress{}
	if err := r.Client.Get(ctx, request.NamespacedName, ingress); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	ctx = controllerutil.ContextWithControllerClient(ctx, r.Client)

	if err := r.applyVirtualService(ctx, ingress); err != nil {
		log.Error(err, "apply virtual service failed")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *IngressReconciler) applyVirtualService(ctx context.Context, ingress *networkingv1.Ingress) error {
	vss := toExportedVirtualServicesByIngress(ingress)

	for i := range vss {
		r.setControllerReference(vss[i], ingress)

		if err := applyVirtualService(ctx, vss[i]); err != nil {
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

func toExportedVirtualServicesByIngress(ingress *networkingv1.Ingress) (vss []*istiov1alpha3.VirtualService) {
	isForceSSLRedirect := false

	if annotations := ingress.GetAnnotations(); annotations != nil && annotations["nginx.ingress.kubernetes.io/force-ssl-redirect"] == "true" {
		isForceSSLRedirect = true
	}

	for i := range ingress.Spec.Rules {
		rule := ingress.Spec.Rules[i]

		if rule.Host == "" || rule.Host == "*" {
			// skip wild ingress
			continue
		}

		vs := &istiov1alpha3.VirtualService{}
		vs.Namespace = ingress.Namespace
		vs.Name = ingress.Name

		vs.Spec.Hosts = []string{rule.Host}

		gatewayName := converter.ConvertToGatewayName(rule.Host)

		if gatewayName != "" {
			vs.Spec.Gateways = append(vs.Spec.Gateways, gatewayName)
		}

		var fallbackRoute *istiotypes.HTTPRoute

		for j := range rule.HTTP.Paths {
			p := rule.HTTP.Paths[j]

			route, isFallback := ingressPathToHttpRoute(&p, isForceSSLRedirect)

			if !isFallback {
				vs.Spec.Http = append(vs.Spec.Http, route)
			} else {
				fallbackRoute = route
			}
		}

		if fallbackRoute != nil {
			vs.Spec.Http = append(vs.Spec.Http, fallbackRoute)
		}

		vss = append(vss, vs)
	}

	return
}

func ingressPathToHttpRoute(p *networkingv1.HTTPIngressPath, isForceSSLRedirect bool) (*istiotypes.HTTPRoute, bool) {
	isFallback := false

	dest := &istiotypes.Destination{
		Host: p.Backend.Service.Name,
	}

	if p.Backend.Service.Port.Number != 0 {
		dest.Port = &istiotypes.PortSelector{
			Number: uint32(p.Backend.Service.Port.Number),
		}
	}

	route := &istiotypes.HTTPRoute{
		Route: []*istiotypes.HTTPRouteDestination{
			{Destination: dest},
		},
	}

	if isForceSSLRedirect {
		route.Headers = &istiotypes.Headers{Request: &istiotypes.Headers_HeaderOperations{
			Set: map[string]string{
				"X-Forwarded-Proto": "https",
			},
		}}
	}

	if p.Path != "" {
		pathType := networkingv1.PathTypePrefix

		if p.PathType != nil {
			pathType = *p.PathType
		}

		switch pathType {
		case networkingv1.PathTypePrefix:
			route.Match = append(route.Match, &istiotypes.HTTPMatchRequest{
				Uri: &istiotypes.StringMatch{
					MatchType: &istiotypes.StringMatch_Prefix{Prefix: p.Path},
				},
			})

			if p.Path == "/" {
				isFallback = true
			}
		case networkingv1.PathTypeExact:
			route.Match = append(route.Match, &istiotypes.HTTPMatchRequest{
				Uri: &istiotypes.StringMatch{
					MatchType: &istiotypes.StringMatch_Exact{Exact: p.Path},
				},
			})
		}
	}

	return route, isFallback
}
