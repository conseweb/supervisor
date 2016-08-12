PWD := $(shell pwd)
PKG := github.com/conseweb/supervisor
VERSION := $(shell cat VERSION.txt)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
LD_FLAGS := -X $(PKG)/version.version=$(VERSION) -X $(PKG)/version.gitCommit=$(GIT_COMMIT)

test:
	docker run --rm \
	 --name supervisor-testing \
	 -v $(PWD):/go/src/$(PKG) \
	 -w /go/src/$(PKG) \
	 ckeyer/obc:base make test-inner

test-inner:
	go test -ldflags="$(LD_FLAGS)" $$(go list ./... |grep -v "vendor")

integration-test:
	# TODO