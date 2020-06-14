VERSION=$(shell cat .version)-dev

up: apply-crd
	operator-sdk run local

gen-deepcopy:
	go run ./hack/tools.go deepcopy ./pkg/strfmt

apply-crd:
	kubectl apply -f deploy/crds/serving.octohelm.tech_qservices_crd.yaml

apply-example:
	kubectl apply --filename ./deploy/examples/service.yaml

dockerx:
	docker buildx build --build-arg=GOPROXY=${GOPROXY} --platform=linux/amd64,linux/arm64 -f build/Dockerfile -t octohelm/qservice-operator:${VERSION} . --push