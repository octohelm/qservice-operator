package serving

import (
	"github.com/octohelm/qservice-operator/apis/serving/v1alpha1"
	"github.com/octohelm/qservice-operator/pkg/apiutil"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var CRDs = []*apiextensionsv1.CustomResourceDefinition{
	qserviceCustomResourceDefinition(),
	qingressCustomResourceDefinition(),
	qegressCustomResourceDefinition(),
}

func qserviceCustomResourceDefinition() *apiextensionsv1.CustomResourceDefinition {
	return apiutil.ToCRD(&apiutil.CustomResourceDefinition{
		GroupVersion: v1alpha1.SchemeGroupVersion,
		KindType:     &v1alpha1.QService{},
		ListKindType: &v1alpha1.QServiceList{},
		ShortNames:   []string{"qsvc"},
	})
}

func qingressCustomResourceDefinition() *apiextensionsv1.CustomResourceDefinition {
	return apiutil.ToCRD(&apiutil.CustomResourceDefinition{
		GroupVersion: v1alpha1.SchemeGroupVersion,
		KindType:     &v1alpha1.QIngress{},
		ListKindType: &v1alpha1.QIngressList{},
		Plural:       "qingresses",
		ShortNames:   []string{"qing"},
	})
}

func qegressCustomResourceDefinition() *apiextensionsv1.CustomResourceDefinition {
	return apiutil.ToCRD(&apiutil.CustomResourceDefinition{
		GroupVersion: v1alpha1.SchemeGroupVersion,
		KindType:     &v1alpha1.QEgress{},
		ListKindType: &v1alpha1.QEgressList{},
		Plural:       "qegresses",
		ShortNames:   []string{"qeg"},
	})
}
