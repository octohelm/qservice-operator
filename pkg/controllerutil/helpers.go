package controllerutil

import (
	"context"
	"strconv"

	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var SetControllerReference = controllerutil.SetControllerReference

var keyControllerGeneration = "controller-generation"

func IsControllerGenerationEqual(cur metav1.Object, next metav1.Object) bool {
	return ControllerGeneration(cur) == ControllerGeneration(next)
}

func ControllerGeneration(obj metav1.Object) string {
	return obj.GetAnnotations()[keyControllerGeneration]
}

func AnnotateControllerGeneration(annotations map[string]string, generation int64) map[string]string {
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[keyControllerGeneration] = strconv.FormatInt(generation, 10)
	return annotations
}

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
