package converter

import (
	"encoding/json"
	"strings"

	"github.com/octohelm/qservice-operator/pkg/constants"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ToDeployment(s *QService) *appsv1.Deployment {
	d := &appsv1.Deployment{}
	d.Namespace = s.Namespace
	d.Name = s.Name

	d.Labels = cloneKV(s.Labels)
	d.Labels["app"] = d.Name

	if _, ok := d.Labels["role"]; !ok {
		if strings.HasPrefix(d.Name, "srv-") || strings.HasPrefix(d.Name, "svc-") {
			d.Labels["role"] = "svc"
		}
		if strings.HasPrefix(s.Name, "web-") {
			d.Labels["role"] = "web"
		}
	}

	d.Annotations = cloneKV(s.Annotations)

	d.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": d.Name,
		},
	}

	d.Spec.Replicas = s.Spec.Replicas
	d.Spec.Template.Labels = cloneKV(d.Labels)
	d.Spec.Template.Annotations = map[string]string{}

	// sync restartedAt to trigger pod restarted
	if restartedAt, ok := d.Annotations[constants.AnnotationRestartedAt]; ok {
		d.Spec.Template.Annotations[constants.AnnotationRestartedAt] = restartedAt
	}

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
