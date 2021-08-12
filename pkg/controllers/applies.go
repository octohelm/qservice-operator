package controllers

import (
	"context"

	"github.com/octohelm/qservice-operator/pkg/apis/serving/v1alpha1"
	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	pkgerrors "github.com/pkg/errors"
	istioneteworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func applyDeployment(ctx context.Context, deployment *appsv1.Deployment) error {
	deployment.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("Deployment"))
	if err := applyResource(ctx, deployment); err != nil {
		return pkgerrors.Wrap(err, "Apply Deployment")
	}
	return nil
}

func applyIngress(ctx context.Context, ingress *networkingv1.Ingress) error {
	ingress.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Ingress"))
	if err := applyResource(ctx, ingress); err != nil {
		return pkgerrors.Wrap(err, "Apply Ingress")
	}
	return nil
}

func applyService(ctx context.Context, service *corev1.Service) error {
	service.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Service"))
	if err := applyResource(ctx, service); err != nil {
		return pkgerrors.Wrap(err, "Apply Service")
	}
	return nil
}

func applySecret(ctx context.Context, secret *corev1.Secret) error {
	secret.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))
	if err := applyResource(ctx, secret); err != nil {
		return pkgerrors.Wrap(err, "Apply Secret")
	}
	return nil
}

func applyVirtualService(ctx context.Context, vs *istioneteworkingv1alpha3.VirtualService) error {
	vs.SetGroupVersionKind(istioneteworkingv1alpha3.SchemeGroupVersion.WithKind("VirtualService"))
	if err := applyResource(ctx, vs); err != nil {
		return pkgerrors.Wrap(err, "Apply VirtualService")
	}
	return nil
}

func applyServiceEntry(ctx context.Context, se *istioneteworkingv1alpha3.ServiceEntry) error {
	se.SetGroupVersionKind(istioneteworkingv1alpha3.SchemeGroupVersion.WithKind("ServiceEntry"))
	if err := applyResource(ctx, se); err != nil {
		return pkgerrors.Wrap(err, "Apply ServiceEntry")
	}
	return nil
}

func applyQIngress(ctx context.Context, qingress *v1alpha1.QIngress) error {
	qingress.SetGroupVersionKind(v1alpha1.SchemeGroupVersion.WithKind("QIngress"))
	if err := applyResource(ctx, qingress); err != nil {
		return pkgerrors.Wrap(err, "Apply QIngress")
	}
	return nil
}

func applyResource(ctx context.Context, ro runtime.Object) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	obj, err := ClientObject(ro)
	if err != nil {
		return err
	}

	gvk := obj.GetObjectKind().GroupVersionKind()

	live, _ := DeepCopyClientObject(obj)

	if err := c.Get(ctx, client.ObjectKeyFromObject(obj), live); err != nil {
		if apierrors.IsNotFound(err) {
			return c.Create(ctx, obj)
		}
		return err
	}

	if !controllerutil.IsControllerGenerationEqual(live, obj) {
		return c.Patch(ctx, obj, PatchFor(gvk, live))
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
	c, err := ClientObject(object.DeepCopyObject())
	if err != nil {
		return nil, err
	}
	return c, nil
}

func PatchFor(gvk schema.GroupVersionKind, live client.Object) client.Patch {
	if gvk.Group == corev1.GroupName && gvk.Kind == "Service" {
		return client.Merge
	}

	// TODO handle more
	if gvk.Group == corev1.GroupName || gvk.Group == appsv1.GroupName {
		return client.StrategicMergeFrom(live)
	}

	return client.MergeFromWithOptions(live, client.MergeFromWithOptimisticLock{})
}

func ClientWithoutCache(c client.Client, r client.Reader) client.Client {
	return &clientWithoutCache{Client: c, r: r}
}

type clientWithoutCache struct {
	r client.Reader
	client.Client
}

func (c *clientWithoutCache) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
	return c.r.Get(ctx, key, obj)
}

func (c *clientWithoutCache) List(ctx context.Context, key client.ObjectList, opts ...client.ListOption) error {
	return c.r.List(ctx, key, opts...)
}
