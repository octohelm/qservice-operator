package converter

import (
	"encoding/json"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ToDeployment(s *QService) *appsv1.Deployment {
	d := &appsv1.Deployment{}
	d.Name = s.Name
	d.Namespace = s.Namespace
	d.Labels = cloneKV(s.Labels)

	d.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": d.Name,
		},
	}

	d.Spec.Replicas = s.Spec.Replicas
	d.Spec.Template.Labels = cloneKV(d.Labels)
	d.Spec.Template.Labels["app"] = d.Name
	d.Spec.Template.Spec = toPodSpec(s, s.Spec.Pod)

	if s.Spec.Strategy != nil {
		switch appsv1.DeploymentStrategyType(s.Spec.Strategy.Type) {
		case appsv1.RecreateDeploymentStrategyType:
			d.Spec.Strategy.Type = appsv1.RecreateDeploymentStrategyType
		case appsv1.RollingUpdateDeploymentStrategyType:
			d.Spec.Strategy.Type = appsv1.RollingUpdateDeploymentStrategyType
			if len(s.Spec.Strategy.Flags) > 0 {
				data, _ := json.Marshal(s.Spec.Strategy.Flags)
				_ = json.Unmarshal(data, d.Spec.Strategy.RollingUpdate)
			}
		}

	}

	return d
}
