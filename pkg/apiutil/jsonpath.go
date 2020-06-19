package apiutil

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
