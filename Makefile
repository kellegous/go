CPP = /usr/bin/cpp -P -undef -Wundef -std=c99 -nostdinc -Wtrigraphs -fdollars-in-identifiers -C -Wno-invalid-pp-token

SRC = $(shell find web/assets -maxdepth 1 -type f)
DST = $(patsubst %.scss,%.css,$(patsubst %.ts,%.js,$(subst web/assets,.build/assets,$(SRC))))

ALL: web/bindata.go

.build/bin/go-bindata:
	GOPATH=$(shell pwd)/.build go get github.com/a-urth/go-bindata/...

.build/assets:
	mkdir -p $@

.build/assets/%.css: web/assets/%.scss
	sass --no-source-map --style=compressed $< $@

.build/assets/%.js: web/assets/%.ts
	$(eval TMP := $(shell mktemp))
	tsc --out $(TMP) $< 
	google-closure-compiler --js $(TMP) --js_output_file $@
	rm -f $(TMP)

.build/assets/%: web/assets/%
	cp $< $@

web/bindata.go: .build/bin/go-bindata .build/assets $(DST)
	$< -o $@ -pkg web -prefix .build/assets -nomemcopy .build/assets/...

clean:
	rm -rf .build/assets web/bindata.go
