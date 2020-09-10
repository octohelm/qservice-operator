package controllers

import (
	"context"

	"github.com/octohelm/qservice-operator/pkg/apiutil"
	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	istiov1beta1 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func applyDeployment(ctx context.Context, deployment *appsv1.Deployment) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	current := &appsv1.Deployment{}

	err := c.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		return c.Create(ctx, deployment)
	}

	if !controllerutil.IsControllerGenerationEqual(current, deployment) {
		return c.Patch(ctx, deployment, apiutil.JSONPatch(types.StrategicMergePatchType))
	}

	return nil
}

func applyIngress(ctx context.Context, ingress *extensionsv1beta1.Ingress) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	current := &extensionsv1beta1.Ingress{}

	err := c.Get(ctx, types.NamespacedName{Name: ingress.Name, Namespace: ingress.Namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		return c.Create(ctx, ingress)
	}

	if !controllerutil.IsControllerGenerationEqual(current, ingress) {
		return c.Patch(ctx, ingress, apiutil.JSONPatch(types.StrategicMergePatchType))
	}

	return nil
}

func applyService(ctx context.Context, service *corev1.Service) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	current := &corev1.Service{}

	err := c.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		return c.Create(ctx, service)
	}

	if !controllerutil.IsControllerGenerationEqual(current, service) {
		return c.Patch(ctx, service, apiutil.JSONPatch(types.MergePatchType))
	}

	return nil
}

func applySecret(ctx context.Context, secret *corev1.Secret) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	current := &corev1.Secret{}

	err := c.Get(ctx, types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		return c.Create(ctx, secret)
	}

	if !controllerutil.IsControllerGenerationEqual(current, secret) {
		return c.Patch(ctx, secret, apiutil.JSONPatch(types.MergePatchType))
	}

	return nil
}

func applyVirtualService(ctx context.Context, vs *istiov1beta1.VirtualService) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	current := &istiov1beta1.VirtualService{}

	err := c.Get(ctx, types.NamespacedName{Name: vs.Name, Namespace: vs.Namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil
		}
		return c.Create(ctx, vs)
	}

	if !controllerutil.IsControllerGenerationEqual(current, vs) {
		return c.Patch(ctx, vs, apiutil.JSONPatch(types.MergePatchType))
	}

	return nil
}
