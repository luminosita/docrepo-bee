OS = $(shell uname | tr A-Z a-z)

export PATH := $(abspath bin/protoc/bin/):$(abspath bin/):${PATH}

PROJ=docrepo-bee
ORG_PATH=github.com/luminosita
REPO_PATH=$(ORG_PATH)/$(PROJ)

VERSION ?= $(shell ./scripts/git-version)

DOCKER_REPO=ghcr.io/luminosita/docrepo-bee
DOCKER_IMAGE=$(DOCKER_REPO):$(VERSION)

$( shell mkdir -p bin )

user=$(shell id -u -n)
group=$(shell id -g -n)

export GOBIN=$(PWD)/bin

LD_FLAGS="-w -X main.version=$(VERSION)"

# Dependency versions

KIND_NODE_IMAGE = "kindest/node:v1.26.0@sha256:3264cbae4b80c241743d12644b2506fff13dce07fcadf29079c1d06a47b399dd"
KIND_TMP_DIR = "$(PWD)/bin/test/bee-kind-kubeconfig"

.PHONY: generate
# generates all
generate: wire api config

# generates wire file
wire:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	GOFLAGS=-mod=mod go generate ./...
	swag init -d internal/bee,internal/infra/http/handlers/documents -g beeServer.go

# runs generate and go install
build: generate bin/bee

bin/bee:
	@mkdir -p bin/
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/bee

.PHONY: release-binary
# create release build
release-binary: LD_FLAGS = "-w -X main.version=$(VERSION) -extldflags \"-static\""
release-binary: generate
	@go build -o /go/bin/bee -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/bee

docker-compose.override.yaml:
	@cd deployments; \
		cp docker-compose.override.yaml.dist docker-compose.override.yaml

.PHONY: up
# Launch docker environment
up: docker-compose.override.yaml
	@cd deployments; \
		if [ docker-compose.override.yaml -ot docker-compose.override.yaml.dist ]; then \
		  diff -u docker-compose.override.yaml docker-compose.override.yaml.dist || \
		  (echo "!!! The distributed docker-compose.override.yaml example changed. Please update your file accordingly (or at least touch it). !!!" && false); \
		fi; \
		docker-compose up -d

.PHONY: down
# Destroy docker environment
down:
	@cd deployments; \
		docker-compose down --volumes --remove-orphans --rmi local \

# Run all tests
test:
	@go test -v ./...

# Run all tests with race check
testrace:
	@go test -v --race ./...

.PHONY: kind-up kind-down kind-tests
# create kind cluster
kind-up:
	@mkdir -p bin/test
	@kind create cluster --image ${KIND_NODE_IMAGE} --kubeconfig ${KIND_TMP_DIR}

# shutdown kind cluster
kind-down:
	@kind delete cluster
	rm ${KIND_TMP_DIR}

kind-tests: export BEE_KUBERNETES_CONFIG_PATH=${KIND_TMP_DIR}
kind-tests: testall

.PHONY: lint lint-fix
# run linter
lint:
	golangci-lint run

.PHONY: fix
# fix lint violations
fix:
	golangci-lint run --fix

.PHONY: docker-image
# start docker image build
docker-image:
	@docker build -f build/package/Dockerfile -t $(DOCKER_IMAGE) .

# delete bin folder
clean:
	@rm -rf bin/

# run all tests with testrace
testall: testrace

FORCE:

.PHONY: test testrace testall

.PHONY: proto
# generate api files
proto:
	@cd api; \
		buf generate --debug

.PHONY: api
api:
	@protoc --proto_path=./api \
		   --proto_path=/Users/milos/go/proto \
		   --go-http_out=paths=source_relative:./api \
		   api/documents/v1/documents.proto


# generate config files
config:
	@cd internal/conf; \
		buf generate --debug

# install dependencies
deps:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
      go install github.com/bufbuild/buf/cmd/buf@latest; \
      go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest; \
	  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest; \
	  go install github.com/envoyproxy/protoc-gen-validate@latest; \
	  go install sigs.k8s.io/kind@latest; \
	  go install github.com/google/wire/cmd/wire@latest; \
	  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
	  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
