.PHONY: default build docker-image test stop clean-images clean

BINARY = kubewatch

VERSION=
BUILD=

PKG            = github.com/bitnami-labs/kubewatch
TRAVIS_COMMIT ?= `git describe --tags`
GOCMD          = go
BUILD_DATE     = `date +%FT%T%z`
GOFLAGS       ?= $(GOFLAGS:)
LDFLAGS       := "-X '$(PKG)/cmd.gitCommit=$(TRAVIS_COMMIT)' \
		          -X '$(PKG)/cmd.buildDate=$(BUILD_DATE)'"

default: build test

build:
	"$(GOCMD)" build ${GOFLAGS} -ldflags ${LDFLAGS} -o "${BINARY}"

docker-image:
	@docker build -t "${BINARY}" .

test:
	"$(GOCMD)" test -race -v ./...

stop:
	@docker stop "${BINARY}"

clean-images: stop
	@docker rmi "${BUILDER}" "${BINARY}"

clean:
	"$(GOCMD)" clean -i

update_minikube: docker-image
	docker tag kubewatch lwabish/kubewatch
	docker push lwabish/kubewatch
	kubectl rollout restart deployment kubewatch -n kubewatch