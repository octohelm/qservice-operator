package controllerutil

import (
	"context"
	"encoding/json"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensionstypesv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

func ApplyCRDs(c *rest.Config, crds ...*apiextensionsv1.CustomResourceDefinition) error {
	cs, err := apiextensionsclientset.NewForConfig(c)
	if err != nil {
		return err
	}

	apis := cs.ApiextensionsV1().CustomResourceDefinitions()

	ctx := context.Background()

	for i := range crds {
		if err := applyCRD(ctx, apis, crds[i]); err != nil {
			return err
		}
	}

	return nil
}

func applyCRD(ctx context.Context, apis apiextensionstypesv1.CustomResourceDefinitionInterface, crd *apiextensionsv1.CustomResourceDefinition) error {
	_, err := apis.Get(ctx, crd.Name, v1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		_, err := apis.Create(ctx, crd, v1.CreateOptions{})
		return err
	}
	data, err := json.Marshal(crd)
	if err != nil {
		return err
	}
	_, err = apis.Patch(ctx, crd.Name, types.MergePatchType, data, v1.PatchOptions{})
	return err
}
