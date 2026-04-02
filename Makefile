SHA = $(shell git rev-parse HEAD)

ASSETS := \
	internal/ui/assets/edit/index.html \
	internal/ui/assets/index.html

.PHONY: all clean develop publish

all: bin/go

bin/go: cmd/go/main.go $(ASSETS) $(shell find internal -name '*.go')
	go build -o $@ ./cmd/go

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

bin/publish: cmd/publish/main.go
	go build -o $@ ./cmd/publish

publish: bin/publish
	bin/publish \
		--tag=latest \
		--tag=$(shell git rev-parse --short $(SHA)) \
		--platform=linux/arm64 \
		--platform=linux/amd64 \
		--image=kellegous/go
