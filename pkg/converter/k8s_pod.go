package converter

import (
	"sort"

	"github.com/octohelm/qservice-operator/pkg/strfmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func toPodSpec(s *QService, pod Pod) v1.PodSpec {
	podSpec := v1.PodSpec{}

	if s.Spec.ImagePullSecret != "" {
		podSpec.ImagePullSecrets = []v1.LocalObjectReference{{Name: s.Spec.ImagePullSecret}}
	}

	podSpec.Volumes = toVolumes(s)
	podSpec.Containers = toContainers(s, pod)

	podSpec.Tolerations = toTolerations(s)
	podSpec.HostAliases = toHostAlias(s)

	podSpec.RestartPolicy = v1.RestartPolicy(pod.RestartPolicy)
	podSpec.TerminationGracePeriodSeconds = pod.TerminationGracePeriodSeconds
	podSpec.ActiveDeadlineSeconds = pod.ActiveDeadlineSeconds
	podSpec.DNSPolicy = v1.DNSPolicy(pod.DNSPolicy)
	podSpec.NodeSelector = pod.NodeSelector
	podSpec.ServiceAccountName = pod.ServiceAccountName

	return podSpec
}

func toContainers(s *QService, pod Pod) []v1.Container {
	c := toContainer(s, pod.Container)
	c.Ports = toContainerPorts(s, s.Spec.Ports)
	return []v1.Container{c}
}

func toContainer(s *QService, c Container) v1.Container {
	container := v1.Container{}

	container.Image = c.Image
	container.Name = s.Name
	container.ImagePullPolicy = v1.PullPolicy(c.ImagePullPolicy)
	container.WorkingDir = c.WorkingDir
	container.Command = c.Command
	container.Args = c.Args
	container.TTY = c.TTY

	if c.LivenessProbe != nil {
		container.LivenessProbe = toProbe(c.LivenessProbe)
	}

	if c.ReadinessProbe != nil {
		container.ReadinessProbe = toProbe(c.ReadinessProbe)
	}

	if c.Lifecycle != nil {
		container.Lifecycle = &v1.Lifecycle{}
		if c.Lifecycle.PostStart != nil {
			container.Lifecycle.PostStart = &c.Lifecycle.PostStart.Handler
		}
		if c.Lifecycle.PreStop != nil {
			container.Lifecycle.PreStop = &c.Lifecycle.PreStop.Handler
		}
	}

	if s.Spec.Resources != nil {
		requests := v1.ResourceList{}
		limits := v1.ResourceList{}

		for resourceName, res := range s.Spec.Resources {
			if res.Request != 0 {
				requests[v1.ResourceName(resourceName)] = resource.MustParse(res.RequestString())
			}
			if res.Limit != 0 {
				limits[v1.ResourceName(resourceName)] = resource.MustParse(res.LimitString())
			}
		}

		container.Resources.Requests = requests
		container.Resources.Limits = limits
	}

	if s.Spec.Envs != nil {
		if c.Envs == nil {
			c.Envs = Envs{}
		}
		container.Env = toEnvVars(c.Envs.Merge(s.Spec.Envs))
	}

	container.VolumeMounts = toVolumeMounts(c)

	return container
}

func toEnvVars(envs Envs) []v1.EnvVar {
	keys := make([]string, 0)

	for k := range envs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	envVars := make([]v1.EnvVar, 0)

	for _, k := range keys {
		envVars = append(envVars, v1.EnvVar{
			Name:  k,
			Value: envs[k],
		})
	}

	return envVars
}

func toVolumeMounts(container Container) []v1.VolumeMount {
	volumeMounts := make([]v1.VolumeMount, 0)
	for _, volumeMount := range container.Mounts {
		volumeMounts = append(volumeMounts, toVolumeMount(volumeMount))
	}
	return volumeMounts
}

func toVolumes(s *QService) []v1.Volume {
	volumes := make([]v1.Volume, 0)
	for name, v := range s.Spec.Volumes {
		volumes = append(volumes, toVolume(name, v))
	}
	return volumes
}

func toVolumeMount(volumeMount strfmt.VolumeMount) v1.VolumeMount {
	return v1.VolumeMount{
		MountPath: volumeMount.MountPath,
		Name:      volumeMount.Name,
		SubPath:   volumeMount.SubPath,
		ReadOnly:  volumeMount.ReadOnly,
	}
}

func toVolume(name string, vs v1.VolumeSource) v1.Volume {
	volume := v1.Volume{
		Name:         name,
		VolumeSource: vs,
	}
	return volume
}

func toContainerPorts(s *QService, ports []strfmt.PortForward) []v1.ContainerPort {
	containerPorts := make([]v1.ContainerPort, 0)

	for _, port := range ports {
		p := v1.ContainerPort{
			ContainerPort: int32(port.ContainerPort),
		}

		if p.ContainerPort == 0 {
			port.ContainerPort = port.Port
		}

		p.Protocol = toProtocol(port.Protocol)

		containerPorts = append(containerPorts, p)
	}

	return containerPorts
}

func toProbe(specProb *Probe) *v1.Probe {
	p := &v1.Probe{}

	p.Handler = specProb.Action.Handler
	p.InitialDelaySeconds = specProb.InitialDelaySeconds
	p.TimeoutSeconds = specProb.TimeoutSeconds
	p.PeriodSeconds = specProb.PeriodSeconds
	p.SuccessThreshold = specProb.SuccessThreshold
	p.FailureThreshold = specProb.FailureThreshold

	return p
}

func toTolerations(s *QService) []v1.Toleration {
	tolerations := make([]v1.Toleration, 0)

	for _, toleration := range s.Spec.Tolerations {
		t := v1.Toleration{
			Key:    toleration.Key,
			Value:  toleration.Value,
			Effect: v1.TaintEffect(toleration.Effect),
		}

		if t.Value == "" {
			t.Operator = "Exists"
		} else {
			t.Operator = "Equal"
		}

		if toleration.TolerationSeconds != nil {
			t.TolerationSeconds = toleration.TolerationSeconds
		}

		tolerations = append(tolerations, t)
	}
	return tolerations
}

func toHostAlias(s *QService) []v1.HostAlias {
	hostAlias := make([]v1.HostAlias, 0)

	for _, h := range s.Spec.Hosts {
		hostAlias = append(hostAlias, v1.HostAlias{
			IP:        h.IP,
			Hostnames: h.HostNames,
		})
	}

	return hostAlias
}
