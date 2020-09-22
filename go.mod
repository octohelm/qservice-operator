module github.com/octohelm/qservice-operator

go 1.15

require (
	github.com/ghodss/yaml v1.0.0
	github.com/go-courier/ptr v1.0.1
	github.com/go-courier/reflectx v1.3.4
	github.com/go-logr/logr v0.2.1
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.13.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/net v0.0.0-20200904194848-62affa334b73 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43 // indirect
	golang.org/x/sys v0.0.0-20200909081042-eff7692f9009 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	istio.io/api v0.0.0-20200910154833-da5469b620b9
	istio.io/client-go v0.0.0-20200908160912-f99162621a1a
	istio.io/gogo-genproto v0.0.0-20200908160912-66171252e3db // indirect
	k8s.io/api v0.19.1
	k8s.io/apiextensions-apiserver v0.19.1
	k8s.io/apimachinery v0.19.1
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog/v2 v2.3.0 // indirect
	k8s.io/kube-openapi v0.0.0-20200831175022-64514a1d5d59 // indirect
	k8s.io/utils v0.0.0-20200821003339-5e75c0163111 // indirect
	mvdan.cc/sh/v3 v3.1.2
	sigs.k8s.io/controller-runtime v0.6.3

)

replace k8s.io/client-go => k8s.io/client-go v0.19.1
