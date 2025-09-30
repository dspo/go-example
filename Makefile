BUILD_DATE ?= "$(shell date +"%Y-%m-%dT%H:%M")"
GIT_SHA=$(shell git rev-parse --short=7 HEAD)
REGISTRY ?= registry.cn-hangzhou.aliyuncs.com/dspo
IMAGE_TAG ?= dev

KIND_NAME ?= go-example-e2e
CLUSTER_NAME ?= go-example-e2e
E2E_NAMESPACE ?= go-example-e2e

export KUBECONFIG = /tmp/$(CLUSTER_NAME).kubeconfig

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

before-build: gofmt openapi
	go mod tidy
	go vet ./...

build-app0-image: before-build
	docker build -f Dockerfile/app0.Dockerfile -t go-example-app0:${IMAGE_TAG} .

push-app0-image:
	docker tag go-example-app0:${IMAGE_TAG} ${REGISTRY}/go-example-app0:${IMAGE_TAG}
	docker push ${REGISTRY}/go-example-app0:${IMAGE_TAG}

build-ginkgo-image: gofmt openapi
	docker build -f Dockerfile/ginkgo.Dockerfile -t ginkgo:dev .

gofmt: ## Apply go fmt
	@gofmt -w -r 'interface{} -> any' .
	@gofmt -w -r 'ginkgo.FIt -> ginkgo.It' test
	@gofmt -w -r 'ginkgo.FContext -> ginkgo.Context' test
	@gofmt -w -r 'ginkgo.FDescribe -> ginkgo.Describe' test
	@gofmt -w -r 'ginkgo.FDescribeTable -> ginkgo.DescribeTable' test
	@go fmt ./...
.PHONY: gofmt

openapi:
	# not implement yet

e2e: kind-up build-app0-image build-ginkgo-image kind-load-images e2e-run

e2e-ginkgo: build-ginkgo-image kind-load-images e2e-run

e2e-run:
	@kubectl delete deployment -l testGroup=application --all-namespaces
	@kubectl apply -f test/framework/manifests/configmap.yaml --create-namespace ${E2E_NAMESPACE}
	@kubectl apply -f test/framework/manifests/ginkgo.yaml
	@kubectl run -n go-example-e2e --rm -i ginkgo --env="DB=mysql" --image ginkgo:dev --overrides='{"spec":{"serviceAccount":"ginkgo" }}' --restart=Never

.PHONY: kind-up
kind-up:
	@kind get clusters 2>&1 | grep -v $(KIND_NAME) \
		&& kind create cluster --name $(KIND_NAME) \
		|| echo "kind cluster already exists"
	@kind get kubeconfig --name $(KIND_NAME) > $$KUBECONFIG
	kubectl wait --for=condition=Ready nodes --all

.PHONY: kind-down
kind-down:
	@kind get clusters 2>&1 | grep $(KIND_NAME) \
		&& kind delete cluster --name $(KIND_NAME) \
		|| echo "kind cluster does not exist"

.PHONY: kind-load-images
kind-load-images:
	@kind load docker-image go-example-app0:dev --name $(KIND_NAME)
	@kind load docker-image ginkgo:dev --name $(KIND_NAME)
