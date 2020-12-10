module github.com/octohelm/qservice-operator

go 1.15

require (
	cloud.google.com/go v0.73.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-courier/ptr v1.0.1
	github.com/go-courier/reflectx v1.3.4
	github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr v0.3.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/googleapis/gnostic v0.5.3 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0 // indirect
	github.com/prometheus/common v0.15.0 // indirect
	github.com/prometheus/procfs v0.2.0 // indirect
	github.com/ugorji/go v1.1.4 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9 // indirect
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20201207223542-d4d67f95c62d // indirect
	golang.org/x/term v0.0.0-20201207232118-ee85cb95a76b // indirect
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/yaml.v2 v2.4.0
	istio.io/api v0.0.0-20201125194658-3cee6a1d3ab4
	istio.io/client-go v1.8.1
	istio.io/gogo-genproto v0.0.0-20201125194658-285dd734f786 // indirect
	k8s.io/api v0.20.0
	k8s.io/apiextensions-apiserver v0.20.0
	k8s.io/apimachinery v0.20.0
	k8s.io/client-go v1.5.1
	k8s.io/klog/v2 v2.4.0 // indirect
	mvdan.cc/sh/v3 v3.2.1
	sigs.k8s.io/controller-runtime v0.6.4
	sigs.k8s.io/structured-merge-diff/v4 v4.0.2 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.20.0
