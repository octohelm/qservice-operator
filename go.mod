module github.com/octohelm/qservice-operator

go 1.15

require (
	cloud.google.com/go v0.74.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-courier/ptr v1.0.1
	github.com/go-courier/reflectx v1.3.4
	github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr v0.3.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.3 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/onsi/gomega v1.10.4
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20201217014255-9d1352758620 // indirect
	golang.org/x/net v0.0.0-20201216054612-986b41b23924 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20201218084310-7d0127a74742 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	gopkg.in/yaml.v2 v2.4.0
	istio.io/api v0.0.0-20201217173512-1f62aaeb5ee3
	istio.io/client-go v1.8.1
	istio.io/gogo-genproto v0.0.0-20201214221444-a4a00ea4011d // indirect
	k8s.io/api v0.20.1
	k8s.io/apiextensions-apiserver v0.20.1
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v1.5.1
	mvdan.cc/sh/v3 v3.2.1
	sigs.k8s.io/controller-runtime v0.7.0
)

replace k8s.io/client-go => k8s.io/client-go v0.20.1
