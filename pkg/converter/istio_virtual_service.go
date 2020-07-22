package converter

import (
	"strings"

	"github.com/octohelm/qservice-operator/pkg/strfmt"
	istiotypes "istio.io/api/networking/v1alpha3"
	istioapis "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

func ToClusterVirtualService(s *QService) *istioapis.VirtualService {
	vs := &istioapis.VirtualService{}
	vs.Name = s.Name
	vs.Labels = s.Labels

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

func ToExportedVirtualService(s *QService, host string, ingresses []strfmt.Ingress) *istioapis.VirtualService {
	vs := &istioapis.VirtualService{}
	vs.Name = s.Name + "-" + hashID(host)
	vs.Labels = s.Labels
	vs.Labels["controlled-by"] = s.Name
	vs.Labels["host"] = vs.Name

	vs.Spec.Hosts = append(vs.Spec.Hosts, host)

	gatewayName := convertToGatewayName(host)

	if gatewayName != "" {
		vs.Spec.Gateways = append(vs.Spec.Gateways, gatewayName)
	}

	route := &istiotypes.HTTPRoute{
		Route: []*istiotypes.HTTPRouteDestination{
			{
				Destination: &istiotypes.Destination{
					Host: s.Name,
				},
			},
		},
	}

	for i := range ingresses {
		ingress := ingresses[i]

		if ingress.Path != "" {
			route.Match = []*istiotypes.HTTPMatchRequest{
				{
					Uri: &istiotypes.StringMatch{
						MatchType: &istiotypes.StringMatch_Prefix{Prefix: ingress.Path},
					},
				},
			}
		}
	}

	vs.Spec.Http = append(vs.Spec.Http, route)

	return vs
}

func convertToGatewayName(host string) string {
	return "istio-ingress/" + strings.Join(strings.Split(host, ".")[1:], "--")
}
