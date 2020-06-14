package converter

import (
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

	return d
}
