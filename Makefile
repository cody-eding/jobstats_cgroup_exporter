# Needs to be defined before including Makefile.common to auto-generate targets
export GOPATH ?= $(firstword $(subst :, ,$(shell go env GOPATH)))
GOLANG_CROSS_VERSION ?= v1.24.0

release:
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
	--env-file .release-env -v `pwd`:/work -w /work \
	ghcr.io/goreleaser/goreleaser-cross:$(GOLANG_CROSS_VERSION) \
	release --clean
