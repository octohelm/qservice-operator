package qservice

import (
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var keyControllerGeneration = "controller-generation"

func isControllerGenerationEqual(cur metav1.Object, next metav1.Object) bool {
	return controllerGeneration(cur) == controllerGeneration(next)
}

func controllerGeneration(obj metav1.Object) string {
	return obj.GetAnnotations()[keyControllerGeneration]
}

func annotateControllerGeneration(annotations map[string]string, generation int64) map[string]string {
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[keyControllerGeneration] = strconv.FormatInt(generation, 10)
	return annotations
}
