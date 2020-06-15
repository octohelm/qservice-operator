module github.com/octohelm/qservice-operator

go 1.13

require (
	github.com/onsi/gomega v1.9.0
	github.com/operator-framework/operator-sdk v0.18.1
	github.com/spf13/pflag v1.0.5
	gopkg.in/yaml.v2 v2.2.8
	istio.io/api v0.0.0-20200518203817-6d29a38039bd
	istio.io/client-go v0.0.0-20200518205227-b36721b13e0b
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.18.2 // Required by prometheus-operator
)
