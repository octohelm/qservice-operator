package controllerutil

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/octohelm/qservice-operator/pkg/converter"
	"istio.io/api/networking/v1alpha3"
	istioneteworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	istioclientset "istio.io/client-go/pkg/clientset/versioned"
	typesv1alpha3 "istio.io/client-go/pkg/clientset/versioned/typed/networking/v1alpha3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"strings"
)

func ApplyGateways(c *rest.Config, crds ...*istioneteworkingv1alpha3.Gateway) error {
	cs, err := istioclientset.NewForConfig(c)
	if err != nil {
		return err
	}

	apis := cs.NetworkingV1alpha3().Gateways("istio-ingress")

	ctx := context.Background()

	for i := range crds {
		if err := applyGateway(ctx, apis, crds[i]); err != nil {
			return err
		}
	}

	return nil
}

func applyGateway(ctx context.Context, apis typesv1alpha3.GatewayInterface, crd *istioneteworkingv1alpha3.Gateway) error {
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

func ParseIngressGateways(gatewaystr string) IngressGateways {
	pairs := strings.Split(gatewaystr, ",")

	ingressGateways := IngressGateways{}

	for i := range pairs {
		nameHost := strings.Split(pairs[i], ":")
		if len(nameHost) != 2 {
			continue
		}
		ingressGateways[nameHost[0]] = nameHost[1]
	}

	return ingressGateways
}

type IngressGateways map[string]string

func (s IngressGateways) ToGateways() []*istioneteworkingv1alpha3.Gateway {
	gateways := make([]*istioneteworkingv1alpha3.Gateway, 0)

	for _, gateway := range s {
		g := &istioneteworkingv1alpha3.Gateway{}
		g.Name = converter.GatewaySafeName(gateway)
		g.Namespace = "istio-ingress"
		g.Labels = map[string]string{
			"kiali_wizard": "Gateway",
		}
		g.Spec.Selector = map[string]string{
			"istio": "ingressgateway",
		}
		g.Spec.Servers = []*v1alpha3.Server{
			{
				Hosts: []string{
					"*." + gateway,
				},
				Port: &v1alpha3.Port{
					Number:   80,
					Name:     "http-80",
					Protocol: "HTTP",
				},
			},
		}

		gateways = append(gateways, g)
	}

	return gateways
}

func (s IngressGateways) IngressGateway(gateway string) (string, bool) {
	for n := range s {
		g := s[n]

		if g == gateway || n == gateway {
			return g, true
		}
	}

	return "", false
}

func (s IngressGateways) IngressGatewayHost(hostname string) (string, bool) {
	for n := range s {
		gateway := s[n]

		if strings.HasSuffix(hostname, "."+gateway) {
			return hostname, true
		}

		if strings.HasSuffix(hostname, "."+n) {
			return strings.TrimSuffix(hostname, n) + gateway, true
		}
	}

	return "", false
}

func (s IngressGateways) AutoInternalIngress(serviceName string, namespace string) (string, bool) {
	h, ok := s["auto-internal"]
	if ok {
		return host(serviceName, namespace, h), true
	}
	return "", false
}

func host(serviceName string, namespace string, gateway string) string {
	return fmt.Sprintf("%s---%s.%s", serviceName, namespace, gateway)
}
