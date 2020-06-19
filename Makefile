VERSION=$(shell cat .version)

up: apply-crd
	operator-sdk run local

dockerx:
	docker buildx build \
		--push \
		--build-arg=GOPROXY=${GOPROXY} \
		--platform=linux/amd64,linux/arm64 \
		-f build/Dockerfile \
		-t octohelm/qservice-operator:${VERSION} .

lint:
	husky hook pre-commit
	husky hook commit-msg

release:
	git push
	git push origin v${VERSION}

include Makefile.apply
include Makefile.gen