package controllers

import (
	"context"

	"github.com/octohelm/qservice-operator/pkg/apis/serving/v1alpha1"
	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	istioneteworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func applyDeployment(ctx context.Context, deployment *appsv1.Deployment) error {
	deployment.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("Deployment"))
	return applyResource(ctx, deployment)
}

func applyIngress(ctx context.Context, ingress *networkingv1.Ingress) error {
	ingress.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Ingress"))
	return applyResource(ctx, ingress)
}

func applyService(ctx context.Context, service *corev1.Service) error {
	service.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Service"))
	return applyResource(ctx, service)
}

func applySecret(ctx context.Context, secret *corev1.Secret) error {
	secret.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))
	return applyResource(ctx, secret)
}

func applyVirtualService(ctx context.Context, vs *istioneteworkingv1alpha3.VirtualService) error {
	vs.SetGroupVersionKind(istioneteworkingv1alpha3.SchemeGroupVersion.WithKind("VirtualService"))
	return applyResource(ctx, vs)
}

func applyServiceEntry(ctx context.Context, se *istioneteworkingv1alpha3.ServiceEntry) error {
	se.SetGroupVersionKind(istioneteworkingv1alpha3.SchemeGroupVersion.WithKind("ServiceEntry"))
	return applyResource(ctx, se)
}

func applyQIngress(ctx context.Context, qingress *v1alpha1.QIngress) error {
	qingress.SetGroupVersionKind(v1alpha1.SchemeGroupVersion.WithKind("QIngress"))
	return applyResource(ctx, qingress)
}

func applyResource(ctx context.Context, ro runtime.Object) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	obj, err := ClientObject(ro)
	if err != nil {
		return err
	}

	live, _ := DeepCopyClientObject(obj)

	if err := c.Get(ctx, client.ObjectKeyFromObject(obj), live); err != nil {
		if apierrors.IsNotFound(err) {
			return c.Create(ctx, obj)
		}
		return err
	}

	if !controllerutil.IsControllerGenerationEqual(live, obj) {
		return c.Patch(ctx, obj, PatchFor(live))
	}

	return nil
}

func ClientObject(ro runtime.Object) (client.Object, error) {
	o, err := meta.Accessor(ro)
	if err != nil {
		return nil, err
	}
	return o.(client.Object), nil
}

func DeepCopyClientObject(object client.Object) (client.Object, error) {
	return ClientObject(object.DeepCopyObject())
}

func PatchFor(live client.Object) client.Patch {
	gvk := live.GetObjectKind().GroupVersionKind()

	if gvk.Group == corev1.GroupName && gvk.Kind == "Service" {
		return client.Merge
	}

	// TODO handle more
	if gvk.Group == corev1.GroupName || gvk.Group == appsv1.GroupName {
		return client.StrategicMergeFrom(live)
	}

	return client.MergeFrom(live)
}
