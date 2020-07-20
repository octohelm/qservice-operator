package qservice

var (
	LabelIstioInjection          = "istio-injection"
	AnnotationImageKeyPullSecret = "serving.octohelm.tech/imagePullSecret"
)

func FlagsFromNamespaceLabels(labels map[string]string) Flags {
	flags := Flags{}

	if v, ok := labels[LabelIstioInjection]; ok && v == "enabled" {
		flags.IstioEnabled = true
	}

	return flags
}

type Flags struct {
	IstioEnabled bool
}
