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
	"strings"

	"github.com/go-logr/logr"
	"github.com/octohelm/qservice-operator/apis/serving/v1alpha1"
	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	"github.com/octohelm/qservice-operator/pkg/converter"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
		Owns(&v1alpha1.QIngress{}).
		Complete(r)
}

func (r *ServiceReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.Log.WithValues("namespace", request.Namespace, "name", request.Name)

	s := &corev1.Service{}
	if err := r.Client.Get(ctx, request.NamespacedName, s); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	ctx = controllerutil.ContextWithControllerClient(ctx, r.Client)

	if err := r.applyAutoQIngress(ctx, s); err != nil {
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

func (r *ServiceReconciler) applyAutoQIngress(ctx context.Context, svc *v1.Service) error {
	if len(svc.Spec.Ports) == 0 {
		return nil
	}

	host := controllerutil.ServiceIngressHost(svc.Name, svc.Namespace, "auto-internal")

	portName := svc.Spec.Ports[0].Name

	if strings.HasPrefix(portName, "http") || strings.HasPrefix(portName, "grpc") {
		ingress := serviceToQIngress(svc, host)
		r.setControllerReference(ingress, svc)

		if err := applyQIngress(ctx, ingress); err != nil {
			return err
		}
	}

	return nil
}

func serviceToQIngress(svc *v1.Service, hostname string) *v1alpha1.QIngress {
	qingress := &v1alpha1.QIngress{}
	qingress.Namespace = svc.Namespace
	qingress.Name = svc.Name + "-" + converter.HashID(hostname)
	qingress.Labels = svc.GetLabels()
	if qingress.Labels == nil {
		qingress.Labels = map[string]string{}
	}
	qingress.Labels[LabelServiceName] = svc.Name
	qingress.Labels[LabelGateway] = getGateway(hostname)

	port := uint16(80)

	if len(svc.Spec.Ports) > 0 {
		port = uint16(svc.Spec.Ports[0].Port)
	}

	qingress.Spec.Ingress.Host = safeDNS1121Host(hostname)
	qingress.Spec.Ingress.Port = port
	qingress.Spec.Backend = networkingv1.IngressBackend{
		Service: &networkingv1.IngressServiceBackend{
			Name: svc.Name,
			Port: networkingv1.ServiceBackendPort{
				Number: int32(port),
			},
		},
	}

	return qingress
}
