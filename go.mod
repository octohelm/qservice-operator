module github.com/octohelm/qservice-operator

go 1.15

require (
	cloud.google.com/go v0.66.0 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/ghodss/yaml v1.0.0
	github.com/go-courier/ptr v1.0.1
	github.com/go-courier/reflectx v1.3.4
	github.com/go-logr/logr v0.2.1
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/onsi/gomega v1.10.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.14.0 // indirect
	github.com/prometheus/procfs v0.2.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/net v0.0.0-20200927032502-5d4f70055728 // indirect
	golang.org/x/sys v0.0.0-20200926100807-9d91bd62050c // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	istio.io/api v0.0.0-20200926011135-d7cf1f5167bf
	istio.io/client-go v0.0.0-20200916161914-94f0e83444ca
	istio.io/gogo-genproto v0.0.0-20200916161914-c65bfcb51be9 // indirect
	k8s.io/api v0.19.2
	k8s.io/apiextensions-apiserver v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog/v2 v2.3.0 // indirect
	k8s.io/kube-openapi v0.0.0-20200923155610-8b5066479488 // indirect
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800 // indirect
	mvdan.cc/sh/v3 v3.1.2
	sigs.k8s.io/controller-runtime v0.6.3
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2
