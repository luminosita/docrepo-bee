OS = $(shell uname | tr A-Z a-z)

export PATH := $(abspath bin/protoc/bin/):$(abspath bin/):${PATH}

PROJ=sample-bee
ORG_PATH=github.com/luminosita
REPO_PATH=$(ORG_PATH)/$(PROJ)

VERSION ?= $(shell ./scripts/git-version)

GOTOOLS = golang.org/x/lint/golint \
	github.com/golangci/golangci-lint/cmd/golangci-lint \
	golang.org/x/tools/cmd/goimports \
	mvdan.cc/gofumpt \
	github.com/bufbuild/buf/cmd/buf\
	github.com/gogo/protobuf/protoc-gen-gogo \
	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
	github.com/mwitkow/go-proto-validators/protoc-gen-govalidators \
	sigs.k8s.io/kind \

DOCKER_REPO=ghcr.io/luminosita/sample-bee
DOCKER_IMAGE=$(DOCKER_REPO):$(VERSION)

$( shell mkdir -p bin )

user=$(shell id -u -n)
group=$(shell id -g -n)

export GOBIN=$(PWD)/bin

LD_FLAGS="-w -X main.version=$(VERSION)"

# Dependency versions

KIND_NODE_IMAGE = "kindest/node:v1.26.0@sha256:3264cbae4b80c241743d12644b2506fff13dce07fcadf29079c1d06a47b399dd"
KIND_TMP_DIR = "$(PWD)/bin/test/bee-kind-kubeconfig"

milos:
	@echo $(LD_FLAGS)

.PHONY: generate
generate:
#	@go generate $(REPO_PATH)/storage/ent/

build: generate bin/bee

bin/bee:
	@mkdir -p bin/
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/sample-bee

.PHONY: release-binary
release-binary: LD_FLAGS = "-w -X main.version=$(VERSION) -extldflags \"-static\""
release-binary: generate
	@go build -o /go/bin/bee -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/bee

docker-compose.override.yaml:
	@cd deployments; \
		cp docker-compose.override.yaml.dist docker-compose.override.yaml

.PHONY: up
up: docker-compose.override.yaml ## Launch the development environment
	@cd deployments; \
		if [ docker-compose.override.yaml -ot docker-compose.override.yaml.dist ]; then \
		  diff -u docker-compose.override.yaml docker-compose.override.yaml.dist || \
		  (echo "!!! The distributed docker-compose.override.yaml example changed. Please update your file accordingly (or at least touch it). !!!" && false); \
		fi; \
		docker-compose up -d

.PHONY: down
down:  ## Destroy the development environment
	@cd deployments; \
		docker-compose down --volumes --remove-orphans --rmi local \

test:
	@go test -v ./...

testrace:
	@go test -v --race ./...

.PHONY: kind-up kind-down kind-tests
kind-up:
	@mkdir -p bin/test
	@kind create cluster --image ${KIND_NODE_IMAGE} --kubeconfig ${KIND_TMP_DIR}

kind-down:
	@kind delete cluster
	rm ${KIND_TMP_DIR}

kind-tests: export BEE_KUBERNETES_CONFIG_PATH=${KIND_TMP_DIR}
kind-tests: testall

.PHONY: lint lint-fix
lint: ## Run linter
	golangci-lint run

.PHONY: fix
fix: ## Fix lint violations
	golangci-lint run --fix

.PHONY: docker-image
docker-image:
	@docker build -f build/package/Dockerfile -t $(DOCKER_IMAGE) .

.PHONY: verify-proto
verify-proto: proto
	@./scripts/git-diff

clean:
	@rm -rf bin/

testall: testrace

FORCE:

.PHONY: test testrace testall

.PHONY: proto
proto:
	@protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. api/v1/*.proto
	#@cp api/v2/*.proto api/

.PHONY: proto-internal
proto-internal:
	#@protoc --go_out=paths=source_relative:. server/internal/*.proto

deps:
	go get $(GOTOOLS)

