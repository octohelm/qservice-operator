PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = $(shell cat internal/version/version)
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
TAG ?= $(VERSION)

GOBIN ?= ./bin

PUSH ?= true
NAMESPACES ?= docker.io/octohelm
TARGETS ?= qservice-operator

DOCKER_BUILDX_BUILD = docker buildx build \
	--label=org.opencontainers.image.source=https://github.com/$(REPO) \
	--label=org.opencontainers.image.revision=$(COMMIT_SHA) \
	--platform=linux/arm64,linux/amd64

run:
	INGRESS_GATEWAYS=auto-internal:hw-infra.rktl.xyz \
	WATCH_NAMESPACE=default \
	go run ./cmd/qservice-operator/main.go

build:
	goreleaser build --snapshot --rm-dist

lint:
	husky hook pre-commit
	husky hook commit-msg

eval:
	cuem eval -w --output=./bin/qservice-operator.cue ./deploy/component

apply.example:
	cuem k apply ./deploy/qservice-opreator.cue


gen-deepcopy:
	deepcopy-gen \
		--output-file-base zz_generated.deepcopy \
		--go-header-file ./boilerplate.go.txt \
		--input-dirs $(PKG)/pkg/apis/serving/v1alpha1,$(PKG)/pkg/strfmt

dockerx: build
	$(foreach target,$(TARGETS),\
		$(DOCKER_BUILDX_BUILD) \
		--build-arg=VERSION=$(VERSION) \
		$(foreach namespace,$(NAMESPACES),--tag=$(namespace)/$(target):$(TAG)) \
		--file=cmd/$(target)/Dockerfile . ;\
	)
