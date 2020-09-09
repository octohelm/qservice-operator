package converter

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func ToService(s *QService) *v1.Service {
	service := &v1.Service{}
	service.Namespace = s.Namespace
	service.Name = s.Name
	service.Labels = s.Labels
	service.Spec = toServiceSpec(s)
	service.Spec.Selector = map[string]string{
		"app": service.Name,
	}
	return service
}

func toServiceSpec(s *QService) v1.ServiceSpec {
	serviceSpec := v1.ServiceSpec{
		Type: v1.ServiceTypeClusterIP,
	}

	for _, port := range s.Spec.Ports {
		servicePort := v1.ServicePort{}

		appProtocol := port.AppProtocol

		if appProtocol == "" {
			appProtocol = "http"
		}

		servicePort.Name = fmt.Sprintf("%s-%d", appProtocol, port.Port)
		servicePort.Port = int32(port.Port)
		servicePort.TargetPort = intstr.FromInt(int(port.ContainerPort))

		if port.IsNodePort {
			serviceSpec.Type = v1.ServiceTypeNodePort
			servicePort.Name = "np-" + servicePort.Name
			servicePort.NodePort = int32(port.Port)
		}

		servicePort.Protocol = toProtocol(port.Protocol)
		serviceSpec.Ports = append(serviceSpec.Ports, servicePort)
	}

	return serviceSpec
}
