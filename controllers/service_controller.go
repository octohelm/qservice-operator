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
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	"github.com/octohelm/qservice-operator/pkg/converter"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
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
		Owns(&extensionsv1beta1.Ingress{}).
		Complete(r)
}

func (r *ServiceReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("namespace", request.Namespace, "name", request.Name)

	s := &corev1.Service{}
	if err := r.Client.Get(ctx, request.NamespacedName, s); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	ctx = controllerutil.ContextWithControllerClient(ctx, r.Client)

	if err := r.applyAutoIngress(ctx, s); err != nil {
		log.Error(err, "apply ingress failed")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ServiceReconciler) setControllerReference(obj metav1.Object, owner metav1.Object) {
	if err := controllerutil.SetControllerReference(owner, obj, r.Scheme); err != nil {
		r.Log.Error(err, "")
	}
	obj.SetAnnotations(controllerutil.AnnotateControllerGeneration(obj.GetAnnotations(), owner.GetGeneration()))
}

var autoIngressHosts = func(v string) map[string]bool {
	if v == "" {
		return map[string]bool{}
	}
	m := map[string]bool{}
	for _, h := range strings.Split(v, ",") {
		m[h] = true
	}
	return m
}(os.Getenv("AUTO_INGRESS_HOSTS"))

func (r *ServiceReconciler) applyAutoIngress(ctx context.Context, svc *v1.Service) error {
	if len(autoIngressHosts) == 0 {
		return nil
	}

	if len(svc.Spec.Ports) == 0 {
		return nil
	}

	portName := svc.Spec.Ports[0].Name

	if strings.HasPrefix(portName, "http") || strings.HasPrefix(portName, "grpc") {
		for autoIngressHost := range autoIngressHosts {
			ingress := serviceToIngress(svc, fmt.Sprintf("%s---%s.%s", svc.Name, svc.Namespace, autoIngressHost))
			r.setControllerReference(ingress, svc)

			if err := applyIngress(ctx, ingress); err != nil {
				return err
			}
		}
	}

	return nil
}

func serviceToIngress(svc *v1.Service, hostname string) *extensionsv1beta1.Ingress {
	ingress := &extensionsv1beta1.Ingress{}
	ingress.Namespace = svc.Namespace
	ingress.Name = svc.Name + "-" + converter.HashID(hostname)
	ingress.Labels = svc.Labels

	ingress.Annotations = map[string]string{
		"kubernetes.io/ingress.class": "nginx",
	}

	port := uint16(80)

	if len(svc.Spec.Ports) > 0 {
		port = uint16(svc.Spec.Ports[0].Port)
	}

	ingress.Spec.Rules = append(ingress.Spec.Rules, extensionsv1beta1.IngressRule{
		Host: hostname,
		IngressRuleValue: extensionsv1beta1.IngressRuleValue{
			HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
				Paths: []extensionsv1beta1.HTTPIngressPath{
					{
						Backend: extensionsv1beta1.IngressBackend{
							ServiceName: svc.Name,
							ServicePort: intstr.FromInt(int(port)),
						},
					},
				},
			},
		},
	})

	return ingress
}
