package qservice

var (
	AnnotationImagePullSecret = "image-pull-secret"
	LabelIstioInjection       = "istio-injection"
	LabelIngressHost          = "ingress-host"
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
