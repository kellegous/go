PROTOC_GEN_GO_VERSION := v1.36.5
PROTOC_GEN_CONNECT_GO_VERSION := v1.19.1
PROTOC_VERSION := 34.1
GOLANGCI_LINT_VERSION := v2.11.4

SHA := $(shell go run github.com/kellegous/glue/build/info --format="{{.SHA}}")
BUILD_NAME := $(shell go run github.com/kellegous/glue/build/info --format="{{.Name}}")

GO_MOD := $(shell go list -m)

BE_PROTOS := \
	go.pb.go \
	go_connect/go.connect.go

ASSETS := \
	internal/ui/assets/edit/index.html \
	internal/ui/assets/index.html

.PHONY: all clean develop nuke publish

.PRECIOUS: $(BE_PROTOS)

all: bin/go bin/golinks

bin/%: cmd/%/main.go $(ASSETS) $(shell find internal -name '*.go')
	go build -o $@ ./cmd/$*

node_modules/.build: package.json
	npm install
	touch $@

internal/ui/assets/edit/index.html: node_modules/.build $(shell find ui -type f)
	npm run build

internal/ui/assets/index.html: node_modules/.build $(shell find ui -type f)
	npm run build

develop: bin/go
	bin/go --addr=:4025 --metrics=true --dev-mode=.:4026

clean:
	rm -rf bin internal/ui/assets

nuke: clean
	rm -rf node_modules

bin/publish: cmd/publish/main.go
	go build -o $@ ./cmd/publish

publish: bin/publish
	bin/publish \
		--tag=latest \
		--tag=$(shell git rev-parse --short $(SHA)) \
		--platform=linux/arm64 \
		--platform=linux/amd64 \
		--image=kellegous/go

%.pb.go: %.proto bin/protoc-gen-go bin/protoc
	bin/protoc --proto_path=. \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go_out=. \
		--go_opt=module=$(GO_MOD) \
		$<

go_connect/go.connect.go: go.proto bin/protoc-gen-connect-go bin/protoc
	bin/protoc --proto_path=. \
		--plugin=protoc-gen-connect-go=bin/protoc-gen-connect-go \
		--connect-go_out=. \
		--connect-go_opt=module=$(GO_MOD) \
		--connect-go_opt=package_suffix=_connect \
		$<

bin/protoc:
	etc/download-protoc $(PROTOC_VERSION)

bin/protoc-gen-go:
	GOBIN="$(CURDIR)/bin" go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

bin/protoc-gen-connect-go:
	GOBIN="$(CURDIR)/bin" go install connectrpc.com/connect/cmd/protoc-gen-connect-go@$(PROTOC_GEN_CONNECT_GO_VERSION)

bin/golangci-lint: Makefile
	GOBIN=$(abspath $(dir $@)) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)