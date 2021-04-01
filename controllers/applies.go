package controllers

import (
	"context"

	"github.com/octohelm/qservice-operator/apis/serving/v1alpha1"
	"github.com/octohelm/qservice-operator/pkg/apiutil"
	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	istioneteworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func applyDeployment(ctx context.Context, deployment *appsv1.Deployment) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	current := &appsv1.Deployment{}

	err := c.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return c.Create(ctx, deployment)
	}

	if !controllerutil.IsControllerGenerationEqual(current, deployment) {
		return c.Patch(ctx, deployment, apiutil.JSONPatch(types.StrategicMergePatchType))
	}

	return nil
}

func applyIngress(ctx context.Context, ingress *networkingv1.Ingress) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	current := &networkingv1.Ingress{}

	err := c.Get(ctx, types.NamespacedName{Name: ingress.Name, Namespace: ingress.Namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
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
			return err
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
			return err
		}
		return c.Create(ctx, secret)
	}

	if !controllerutil.IsControllerGenerationEqual(current, secret) {
		return c.Patch(ctx, secret, apiutil.JSONPatch(types.MergePatchType))
	}

	return nil
}

func applyVirtualService(ctx context.Context, vs *istioneteworkingv1alpha3.VirtualService) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	current := &istioneteworkingv1alpha3.VirtualService{}

	err := c.Get(ctx, types.NamespacedName{Name: vs.Name, Namespace: vs.Namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return c.Create(ctx, vs)
	}

	if !controllerutil.IsControllerGenerationEqual(current, vs) {
		return c.Patch(ctx, vs, apiutil.JSONPatch(types.MergePatchType))
	}

	return nil
}

func applyServiceEntry(ctx context.Context, se *istioneteworkingv1alpha3.ServiceEntry) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	current := &istioneteworkingv1alpha3.ServiceEntry{}

	err := c.Get(ctx, types.NamespacedName{Name: se.Name, Namespace: se.Namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return c.Create(ctx, se)
	}

	if !controllerutil.IsControllerGenerationEqual(current, se) {
		return c.Patch(ctx, se, apiutil.JSONPatch(types.MergePatchType))
	}

	return nil
}

func applyQIngress(ctx context.Context, qingress *v1alpha1.QIngress) error {
	c := controllerutil.ControllerClientFromContext(ctx)

	current := &v1alpha1.QIngress{}

	err := c.Get(ctx, types.NamespacedName{Name: qingress.Name, Namespace: qingress.Namespace}, current)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return c.Create(ctx, qingress)
	}

	if !controllerutil.IsControllerGenerationEqual(current, qingress) {
		return c.Patch(ctx, qingress, apiutil.JSONPatch(types.StrategicMergePatchType))
	}

	return nil
}
