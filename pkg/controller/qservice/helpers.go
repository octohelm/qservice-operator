package qservice

import (
	"encoding/json"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var keyControllerGeneration = "controller-generation"

func isControllerResourceVersionEqual(cur metav1.Object, next metav1.Object) bool {
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

func JSONPatch(patchType types.PatchType) client.Patch {
	return &jsonPath{patchType: patchType}
}

type jsonPath struct {
	patchType types.PatchType
}

func (j *jsonPath) Type() types.PatchType {
	return j.patchType
}

func (j *jsonPath) Data(obj runtime.Object) ([]byte, error) {
	return json.Marshal(obj)
}
