CPP = /usr/bin/cpp -P -undef -Wundef -std=c99 -nostdinc -Wtrigraphs -fdollars-in-identifiers -C -Wno-invalid-pp-token

SRC = $(shell find web/assets -maxdepth 1 -type f)
DST = $(patsubst %.scss,%.css,$(patsubst %.ts,%.js,$(subst web/assets,.build/assets,$(SRC))))

ALL: web/bindata.go

.build/bin/go-bindata:
	GOPATH=$(shell pwd)/.build go get github.com/go-bindata/go-bindata/...

.build/assets:
	mkdir -p $@

.build/assets/%.css: web/assets/%.scss
	sass --no-source-map --style=compressed $< $@

.build/assets/%.js: web/assets/%.ts
	$(eval TMP := $(shell mktemp))
	tsc --out $(TMP) $<
	npx google-closure-compiler --js $(TMP) --js_output_file $@
	rm -f $(TMP)

.build/assets/%: web/assets/%
	cp $< $@

web/bindata.go: .build/bin/go-bindata .build/assets $(DST)
	$< -o $@ -pkg web -prefix .build/assets -nomemcopy .build/assets/...

clean:
	rm -rf .build/assets web/bindata.go

build-local: clean web/bindata.go
	@echo $(shell echo 'Stopping running container'; docker {stop,rm} go-links 2>/dev/null)
	docker build --rm -t stgarf/go-links:dev .
	$(eval dangling = $(shell docker images -f dangling=true -q --no-trunc | tr '\n' ' ' | sed -e 's/sha256://g'))
	docker rmi $(dangling) 2>/dev/null || true

run-local: build-local
	docker run --rm -dv /tmp/data:/data -p 8067:8067 --name go-links stgarf/go-links:dev

docker-push-remote:
	docker build -t stgarf/go-links:latest .
	docker push stgarf/go-links:latest
