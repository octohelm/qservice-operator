package qservice

import (
	"context"

	istiov1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type contextKeyControllerClient int

func ContextWithControllerClient(ctx context.Context, client client.Client) context.Context {
	return context.WithValue(ctx, contextKeyControllerClient(1), client)
}

func ControllerClientFromContext(ctx context.Context) client.Client {
	if i, ok := ctx.Value(contextKeyControllerClient(1)).(client.Client); ok {
		return i
	}
	return nil
}

func applyDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) error {
	c := ControllerClientFromContext(ctx)

	setNamespace(deployment, namespace)

	current := &appsv1.Deployment{}

	err := c.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		return c.Create(ctx, deployment)
	}

	if !isControllerResourceVersionEqual(current, deployment) {
		return c.Patch(ctx, deployment, JSONPatch(types.StrategicMergePatchType))
	}

	return nil
}

func applyIngress(ctx context.Context, namespace string, ingress *extensionsv1beta1.Ingress) error {
	c := ControllerClientFromContext(ctx)

	setNamespace(ingress, namespace)

	current := &extensionsv1beta1.Ingress{}

	err := c.Get(ctx, types.NamespacedName{Name: ingress.Name, Namespace: namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		return c.Create(ctx, ingress)
	}

	if !isControllerResourceVersionEqual(current, ingress) {
		return c.Patch(ctx, ingress, JSONPatch(types.StrategicMergePatchType))
	}

	return nil
}

func applyService(ctx context.Context, namespace string, service *corev1.Service) error {
	c := ControllerClientFromContext(ctx)

	setNamespace(service, namespace)

	current := &corev1.Service{}

	err := c.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		return c.Create(ctx, service)
	}

	if !isControllerResourceVersionEqual(current, service) {
		return c.Patch(ctx, service, JSONPatch(types.MergePatchType))
	}

	return nil
}

func applyVirtualService(ctx context.Context, namespace string, vs *istiov1alpha3.VirtualService) error {
	c := ControllerClientFromContext(ctx)

	setNamespace(vs, namespace)

	current := &istiov1alpha3.VirtualService{}

	err := c.Get(ctx, types.NamespacedName{Name: vs.Name, Namespace: namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		return c.Create(ctx, vs)
	}

	if !isControllerResourceVersionEqual(current, vs) {
		return c.Patch(ctx, vs, JSONPatch(types.MergePatchType))
	}

	return nil
}

func applySecret(ctx context.Context, namespace string, secret *corev1.Secret) error {
	c := ControllerClientFromContext(ctx)

	setNamespace(secret, namespace)

	current := &corev1.Secret{}

	err := c.Get(ctx, types.NamespacedName{Name: secret.Name, Namespace: namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		return c.Create(ctx, secret)
	}

	if !isControllerResourceVersionEqual(current, secret) {
		return c.Patch(ctx, secret, JSONPatch(types.MergePatchType))
	}

	return nil
}

func setNamespace(o metav1.Object, namespace string) {
	o.SetNamespace(namespace)
}
