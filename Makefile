PWD := $(shell pwd)
PKG := github.com/conseweb/supervisor
VERSION := $(shell cat VERSION.txt)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
LD_FLAGS := -X $(PKG)/version.version=$(VERSION) -X $(PKG)/version.gitCommit=$(GIT_COMMIT)
APP := supervisor
IMAGE := conseweb/supervisor:$(GIT_BRANCH)

default: unit-test

test: unit-test integration-test clear

unit-test: 
	docker run --rm \
	 --name supervisor-testing \
	 -v $(PWD):/go/src/$(PKG) \
	 -w /go/src/$(PKG) \
	 ckeyer/obc:base make testInner

testInner: 
	go test -ldflags="$(LD_FLAGS)" $$(go list ./... |grep -v "vendor"|grep -v "integration-tests")

integration-test: clear build-image
	docker run -d --name supervisor $(IMAGE)
	docker run --rm \
	 --name supervisor-testing \
	 --link supervisor \
	 -e SUPERVISOR_ADDR="supervisor:9376" \
	 -v $(PWD):/go/src/$(PKG) \
	 -w /go/src/$(PKG) \
	 ckeyer/obc:dev tools/integration-test.sh
	-docker rm -f supervisor

build: 
	docker run --rm \
	 --name supervisor-building \
	 -v $(PWD):/go/src/$(PKG) \
	 -w /go/src/$(PKG) \
	 ckeyer/obc:base go build -o bundles/$(APP) -ldflags="$(LD_FLAGS)" .

build-local:
	go build -o bundles/$(APP) -ldflags="$(LD_FLAGS)" .

build-image:
	docker build -t $(IMAGE) .

clear:
	-rm -rf bundles
	-docker rm -f supervisor-testing
	-docker rm -f supervisor-building
	-docker rm -f supervisor
	-docker rmi $(IMAGE)
