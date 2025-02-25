SHA := $(shell git rev-parse HEAD)

ASSETS := \
	internal/ui/assets/edit.html \
	internal/ui/assets/links.html

.PHONY: all clean

all: bin/go

bin/go: $(ASSETS) $(shell find internal -name '*.go')
	go build -o $@ ./cmd/go

node_modules/.build: package.json
	npm install
	touch $@

internal/ui/assets/edit.html: node_modules/.build $(shell find ui -type f)
	npm run build

internal/ui/assets/links.html: node_modules/.build $(shell find ui -type f)
	npm run build

clean:
	rm -rf bin internal/ui/assets
