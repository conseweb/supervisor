PWD := $(shell pwd)
PKG := github.com/conseweb/supervisor
VERSION := $(shell cat VERSION.txt)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
LD_FLAGS := -X $(PKG)/version.version=$(VERSION) -X $(PKG)/version.gitCommit=$(GIT_COMMIT)
APP := supervisor
IMAGE := conseweb/supervisor:$(GIT_BRANCH)
INNER_GOPATH := /opt/gopath

UNIT_TEST_CONTAINER := supervisor-unittest-$(GIT_COMMIT)
SV_CONTAINER := supervisor-$(GIT_COMMIT)
INTE_TEST_CONTAINER := supervisor-testing-$(GIT_COMMIT)
BUILD_CONTAINER := supervisor-building-$(GIT_COMMIT)

default: unit-test

test: unit-test integration-test clear

unit-test: 
	docker run --rm \
	 --name $(UNIT_TEST_CONTAINER) \
	 -v $(PWD):$(INNER_GOPATH)/src/$(PKG) \
	 -w $(INNER_GOPATH)/src/$(PKG) \
	 ckeyer/obc:dev make testInner

testInner: 
	go test -ldflags="$(LD_FLAGS)" $$(go list ./... |grep -v "vendor"|grep -v "integration-tests")

integration-test: clear build-image
	docker run -d --name $(SV_CONTAINER) $(IMAGE)
	docker run --rm \
	 --name $(INTE_TEST_CONTAINER) \
	 --link supervisor \
	 -e SUPERVISOR_ADDR="supervisor:9376" \
	 -v $(PWD):$(INNER_GOPATH)/src/$(PKG) \
	 -w $(INNER_GOPATH)/src/$(PKG) \
	 ckeyer/obc:dev tools/integration-test.sh
	-docker rm -f supervisor

build: 
	docker run --rm \
	 --name $(BUILD_CONTAINER) \
	 -v $(PWD):$(INNER_GOPATH)/src/$(PKG) \
	 -w $(INNER_GOPATH)/src/$(PKG) \
	 ckeyer/obc:dev go build -o bundles/$(APP) -ldflags="$(LD_FLAGS)" .

build-local:
	go build -o bundles/$(APP) -ldflags="$(LD_FLAGS)" .

build-image:
	docker build -t $(IMAGE) .

clear:
	-rm -rf bundles
	-docker rm -f UNIT_TEST_CONTAINER
	-docker rm -f SV_CONTAINER
	-docker rm -f INTE_TEST_CONTAINER
	-docker rm -f BUILD_CONTAINER
	-docker rmi $(IMAGE)
