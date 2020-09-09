PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = $(shell cat .version)
COMMIT_SHA ?= $(shell git describe --always)-devel

GOBUILD = CGO_ENABLED=0 go build -ldflags "-X ${PKG}/version.Version=${VERSION}+sha.${COMMIT_SHA}"

GOBIN ?= ./bin
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

run: apply-crd
	AUTO_INGRESS_HOSTS=hw-dev.rktl.xyz \
	WATCH_NAMESPACE=default \
	go run ./main.go

build:
	$(GOBUILD) -o $(GOBIN)/qservice-operator ./main.go

build.dockerx:
	docker buildx build \
		--push \
		--build-arg=GOPROXY=${GOPROXY} \
		--platform=linux/amd64,linux/arm64 \
		-t octohelm/qservice-operator:${VERSION} \
		-t hub-dev.demo.querycap.com/octohelm/qservice-operator:${VERSION} \
		-f Dockerfile .

lint:
	husky hook pre-commit
	husky hook commit-msg

release:
	git push
	git push origin v${VERSION}

apply-crd:
	kubectl apply -f deploy/crds/serving.octohelm.tech_qservices_crd.yaml

apply:
	kubectl apply -k ./deploy

delete:
	kubectl delete -k ./deploy

apply.example:
	kubectl apply --filename ./deploy/examples/service.yaml

delete.example:
	kubectl delete --filename ./deploy/examples/service.yaml

test.apply: apply
	$(MAKE) apply.example
	kubectl get deployments -n default | grep srv-test

test.delete:
	$(MAKE) delete.example
	kubectl get deployments -n default | grep srv-test

gen-deepcopy: install-deepcopy-gen
	deepcopy-gen \
		--output-file-base zz_generated.deepcopy \
		--go-header-file ./hack/boilerplate.go.txt \
		--input-dirs $(PKG)/apis/serving/v1alpha1,$(PKG)/pkg/strfmt

install-deepcopy-gen:
	go install k8s.io/code-generator/cmd/deepcopy-gen