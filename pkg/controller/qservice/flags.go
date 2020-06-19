package qservice

var (
	LabelIngressHost             = "ingress-host"
	LabelIstioInjection          = "istio-injection"
	AnnotationImageKeyPullSecret = "serving.octohelm.tech/imagePullSecret"
)

func FlagsFromNamespaceLabels(labels map[string]string) Flags {
	flags := Flags{}

	if v, ok := labels[LabelIstioInjection]; ok && v == "enabled" {
		flags.IstioEnabled = true
	}

	if v, ok := labels[LabelIngressHost]; ok && v != "" && v != "none" {
		flags.HostBase = v
	}

	return flags
}

type Flags struct {
	IstioEnabled bool
	HostBase     string
}
