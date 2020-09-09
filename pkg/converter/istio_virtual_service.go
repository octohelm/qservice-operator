package converter

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"

	istiotypes "istio.io/api/networking/v1alpha3"
	istioapis "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

func ToClusterVirtualServiceFromService(s *corev1.Service) *istioapis.VirtualService {
	vs := &istioapis.VirtualService{}
	vs.Namespace = s.Namespace
	vs.Name = s.Name

	vs.Labels = s.Labels
	vs.Annotations = s.Annotations

	vs.Spec.Hosts = []string{s.Name}

	vs.Spec.Http = []*istiotypes.HTTPRoute{
		{
			Route: []*istiotypes.HTTPRouteDestination{
				{
					Destination: &istiotypes.Destination{
						Host: s.Name,
					},
				},
			},
		},
	}

	return vs
}

func ToExportedVirtualServicesByIngress(ingress *extensionsv1beta1.Ingress) (list []*istioapis.VirtualService) {
	for i := range ingress.Spec.Rules {
		r := ingress.Spec.Rules[i]

		vs := &istioapis.VirtualService{}
		vs.Namespace = ingress.Namespace
		vs.Name = ingress.Name + "-" + hashID(r.Host)
		vs.Spec.Hosts = []string{r.Host}

		gatewayName := convertToGatewayName(r.Host)

		if gatewayName != "" {
			vs.Spec.Gateways = append(vs.Spec.Gateways, gatewayName)
		}

		for j := range r.HTTP.Paths {
			p := r.HTTP.Paths[j]

			route := &istiotypes.HTTPRoute{
				Route: []*istiotypes.HTTPRouteDestination{
					{
						Destination: &istiotypes.Destination{
							Host: p.Backend.ServiceName,
						},
					},
				},
			}

			if p.Path != "" {
				route.Match = []*istiotypes.HTTPMatchRequest{
					{
						Uri: &istiotypes.StringMatch{
							MatchType: &istiotypes.StringMatch_Prefix{Prefix: p.Path},
						},
					},
				}
			}

			v := vs.DeepCopy()
			v.Name = vs.Name + fmt.Sprintf("-%d", j)
			v.Spec.Http = append(vs.Spec.Http, route)

			list = append(list, v)
		}

	}

	return list
}

func convertToGatewayName(host string) string {
	return "istio-ingress/" + strings.Join(strings.Split(host, ".")[1:], "--")
}
