package converter

import (
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func ToIngress(s *QService) *extensionsv1beta1.Ingress {
	ingress := &extensionsv1beta1.Ingress{}
	ingress.Namespace = s.Namespace
	ingress.Name = s.Name
	ingress.Labels = s.Labels

	ingress.Annotations = map[string]string{
		"kubernetes.io/ingress.class": "nginx",
	}

	for _, r := range s.Spec.Ingresses {
		rule := extensionsv1beta1.IngressRule{
			Host: r.Host,
			IngressRuleValue: extensionsv1beta1.IngressRuleValue{
				HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
					Paths: []extensionsv1beta1.HTTPIngressPath{
						{
							Path: r.Path,
							Backend: extensionsv1beta1.IngressBackend{
								ServiceName: s.Name,
								ServicePort: intstr.FromInt(int(r.Port)),
							},
						},
					},
				},
			},
		}

		ingress.Spec.Rules = append(ingress.Spec.Rules, rule)
	}

	return ingress
}
