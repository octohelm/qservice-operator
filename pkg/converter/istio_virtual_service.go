package converter

import (
	"fmt"
	"strings"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"

	istiotypes "istio.io/api/networking/v1alpha3"
	istioapis "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

//func ToClusterVirtualServiceFromService(s *corev1.Service) *istioapis.VirtualService {
//	if s.Spec.Type != corev1.ServiceTypeClusterIP {
//		return nil
//	}
//
//	vs := &istioapis.VirtualService{}
//	vs.Namespace = s.Namespace
//	vs.Name = s.Name
//	vs.Labels = s.Labels
//	vs.Annotations = s.Annotations
//
//	vs.Spec.Hosts = []string{s.Name}
//
//	for i := range s.Spec.Ports {
//		port := s.Spec.Ports[i]
//
//		if port.Name != "" {
//			istio.P
//
//
//		}
//
//		vs.Spec.Tcp = []*istiotypes.TCPRoute{
//			{
//				Route: []*istiotypes.RouteDestination{
//					{
//						Destination: &istiotypes.Destination{
//							Host: s.Name,
//						},
//					},
//				},
//			},
//		}
//	}
//
//	return vs
//}

func ToExportedVirtualServicesByIngress(ingress *extensionsv1beta1.Ingress) (list []*istioapis.VirtualService) {
	for i := range ingress.Spec.Rules {
		r := ingress.Spec.Rules[i]

		vs := &istioapis.VirtualService{}
		vs.Namespace = ingress.Namespace
		vs.Name = ingress.Name + "-" + HashID(r.Host)
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
				stringMatch := &istiotypes.StringMatch{
					MatchType: &istiotypes.StringMatch_Prefix{Prefix: p.Path},
				}

				if p.PathType != nil && *p.PathType == extensionsv1beta1.PathTypeExact {
					stringMatch = &istiotypes.StringMatch{
						MatchType: &istiotypes.StringMatch_Exact{Exact: p.Path},
					}
				}

				route.Match = []*istiotypes.HTTPMatchRequest{
					{
						Uri: stringMatch,
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

func GatewaySafeName(gateway string) string {
	return strings.Join(strings.Split(gateway, "."), "--")
}
