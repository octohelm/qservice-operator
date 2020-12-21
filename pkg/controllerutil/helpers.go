package controllerutil

import (
	"context"
	"strconv"

	"github.com/octohelm/qservice-operator/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var SetControllerReference = controllerutil.SetControllerReference

func IsControllerGenerationEqual(cur metav1.Object, next metav1.Object) bool {
	if nextOwner := metav1.GetControllerOf(next); nextOwner != nil {
		if curOwner := metav1.GetControllerOf(cur); curOwner != nil {
			if curOwner.UID != nextOwner.UID {
				return false
			}
		}
	}

	annotations := cur.GetAnnotations()
	nextAnnotations := next.GetAnnotations()

	return isEqualProp(annotations, nextAnnotations, constants.AnnotationControllerGeneration) && isEqualProp(annotations, nextAnnotations, constants.AnnotationRestartedAt)
}

func isEqualProp(cur map[string]string, next map[string]string, prop string) bool {
	if cur == nil {
		return false
	}
	if next == nil {
		return false
	}
	return cur[prop] == next[prop]
}

func AnnotateControllerGeneration(annotations map[string]string, generation int64) map[string]string {
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[constants.AnnotationControllerGeneration] = strconv.FormatInt(generation, 10)
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
