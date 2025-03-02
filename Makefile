ifndef SHA
	SHA := $(shell git rev-parse HEAD)
endif

ASSETS := \
	internal/ui/assets/edit/index.html \
	internal/ui/assets/index.html

.PHONY: all clean develop

all: bin/go

bin/go: cmd/go/main.go $(ASSETS) $(shell find internal -name '*.go')
	go build -o $@ ./cmd/go

bin/devserver: cmd/devserver/main.go $(ASSETS)
	go build -o $@ ./cmd/devserver

node_modules/.build: package.json
	npm install
	touch $@

internal/ui/assets/edit/index.html: node_modules/.build $(shell find ui -type f)
	npm run build

internal/ui/assets/index.html: node_modules/.build $(shell find ui -type f)
	npm run build

develop: bin/devserver bin/go
	bin/devserver

clean:
	rm -rf bin internal/ui/assets

bin/buildimg:
	GOBIN="$(CURDIR)/bin" go install github.com/kellegous/buildimg@latest

bin/publish: cmd/publish/main.go
	go build -o $@ ./cmd/publish

publish: bin/publish
	bin/publish \
		--tag=latest \
		--tag=$(shell git rev-parse --short $(SHA))
		--platform=linux/arm64 \
		--platform=linux/amd64 \
		--build-arg=SHA=${SHA} \
		--image=kellegous/go
