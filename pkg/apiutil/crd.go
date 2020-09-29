package apiutil

import (
	"reflect"
	"strings"

	"github.com/go-courier/ptr"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type CustomResourceDefinition struct {
	GroupVersion schema.GroupVersion
	KindType     runtime.Object
	ListKindType runtime.Object
	Plural       string
	ShortNames   []string
}

func ToCRD(d *CustomResourceDefinition) *apiextensionsv1.CustomResourceDefinition {
	crd := &apiextensionsv1.CustomResourceDefinition{}

	crdNames := apiextensionsv1.CustomResourceDefinitionNames{
		Kind:       reflect.Indirect(reflect.ValueOf(d.KindType)).Type().Name(),
		ListKind:   reflect.Indirect(reflect.ValueOf(d.ListKindType)).Type().Name(),
		ShortNames: d.ShortNames,
	}

	crdNames.Singular = strings.ToLower(crdNames.Kind)

	if d.Plural != "" {
		crdNames.Plural = d.Plural
	} else {
		crdNames.Plural = crdNames.Singular + "s"
	}

	crd.Name = crdNames.Plural + "." + d.GroupVersion.Group
	crd.Spec.Group = d.GroupVersion.Group
	crd.Spec.Scope = apiextensionsv1.NamespaceScoped

	crd.Spec.Names = crdNames
	crd.Spec.Versions = []apiextensionsv1.CustomResourceDefinitionVersion{
		{
			Name:    d.GroupVersion.Version,
			Served:  true,
			Storage: true,
			Schema: &apiextensionsv1.CustomResourceValidation{
				OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
					XPreserveUnknownFields: ptr.Bool(true),
				},
			},
			Subresources: &apiextensionsv1.CustomResourceSubresources{
				Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
			},
		},
	}

	return crd
}
