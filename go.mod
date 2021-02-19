module github.com/octohelm/qservice-operator

go 1.16

require (
	cloud.google.com/go v0.77.0 // indirect
	github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6 // indirect
	github.com/coreos/go-etcd v2.0.0+incompatible // indirect
	github.com/cpuguy83/go-md2man v1.0.10 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-courier/ptr v1.0.1
	github.com/go-courier/reflectx v1.3.4
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0 // indirect
	github.com/go-openapi/validate v0.19.5 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.4 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/onsi/gomega v1.10.4
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/prometheus/common v0.17.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/ugorji/go v1.1.4 // indirect
	github.com/ugorji/go/codec v0.0.0-20181204163529-d75b2dcb6bc8 // indirect
	github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20210218145215-b8e89b74b9df // indirect
	golang.org/x/oauth2 v0.0.0-20210218202405-ba52d332ba99 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20210218155724-8ebf48af031b // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	istio.io/api v0.0.0-20210219010445-724943e9da20
	istio.io/client-go v1.9.0
	istio.io/gogo-genproto v0.0.0-20210218000043-b2cfc07e0f6b // indirect
	k8s.io/api v0.20.4
	k8s.io/apiextensions-apiserver v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
	k8s.io/klog/v2 v2.5.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210216185858-15cd8face8d6 // indirect
	mvdan.cc/sh/v3 v3.2.2
	sigs.k8s.io/controller-runtime v0.8.2
	sigs.k8s.io/structured-merge-diff/v4 v4.0.3 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.20.4
