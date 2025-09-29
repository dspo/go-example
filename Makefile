BUILD_DATE ?= "$(shell date +"%Y-%m-%dT%H:%M")"
GIT_SHA=$(shell git rev-parse --short=7 HEAD)
REGISTRY ?= registry.cn-hangzhou.aliyuncs.com/dspo
IMAGE_TAG ?= dev

GOOS ?= linux
GOARCH ?= amd64

ifeq ($(shell uname -s),Darwin)
	GOOS = darwin
endif

ifeq ($(shell uname -m),arm64)
	GOARCH = arm64
endif
ifeq ($(shell uname -m), aarch64)
	GOARCH = arm64
endif

build-app0-image:
	docker build -f Dockerfile/app0.Dockerfile -t ${REGISTRY}/go-exmaple-app0:${IMAGE_TAG} .

push-app0-image:
	docker push ${REGISTRY}/go-exmaple-app0:${IMAGE_TAG}

openapi:
	# not implement yet

e2e-test:
	# not implement yet

gofmt: ## Apply go fmt
	@gofmt -w -r 'interface{} -> any' .
	@gofmt -w -r 'ginkgo.FIt -> ginkgo.It' test
	@gofmt -w -r 'ginkgo.FContext -> ginkgo.Context' test
	@gofmt -w -r 'ginkgo.FDescribe -> ginkgo.Describe' test
	@gofmt -w -r 'ginkgo.FDescribeTable -> ginkgo.DescribeTable' test
	@go fmt ./...
.PHONY: gofmt
