package serving

import (
	"github.com/octohelm/qservice-operator/apis/serving/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func QServiceCustomResourceDefinition() *apiextensionsv1.CustomResourceDefinition {
	crd := &apiextensionsv1.CustomResourceDefinition{}

	crdNames := apiextensionsv1.CustomResourceDefinitionNames{
		Kind:       "QService",
		ListKind:   "QServiceList",
		Plural:     "qservices",
		Singular:   "qservice",
		ShortNames: []string{"qsvc"},
	}

	crd.Name = crdNames.Plural + "." + v1alpha1.SchemeGroupVersion.Group
	crd.Spec.Group = v1alpha1.SchemeGroupVersion.Group
	crd.Spec.Scope = apiextensionsv1.NamespaceScoped

	crd.Spec.Names = crdNames
	crd.Spec.Versions = []apiextensionsv1.CustomResourceDefinitionVersion{
		{
			Name:    "v1alpha1",
			Served:  true,
			Storage: true,
			Subresources: &apiextensionsv1.CustomResourceSubresources{
				Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
			},
		},
	}

	return crd
}
