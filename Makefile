CPP = /usr/bin/cpp -P -undef -Wundef -std=c99 -nostdinc -Wtrigraphs -fdollars-in-identifiers -C -Wno-invalid-pp-token

SRC = $(shell find web/assets -maxdepth 1 -type f)
DST = $(subst web/assets,.build/assets,$(SRC))

ALL: web/bindata.go

.build/bin/go-bindata:
	GOPATH=$(shell pwd)/.build go get github.com/jteeuwen/go-bindata/...

.build/assets:
	mkdir -p $@

.build/assets/%.js: web/assets/%.js
	$(CPP) $< | closure-compiler --js_output_file $@

.build/assets/%: web/assets/%
	cp $< $@

web/bindata.go: .build/bin/go-bindata .build/assets $(DST)
	$< -o $@ -pkg web -prefix .build/assets -nomemcopy .build/assets/...

clean:
	rm -rf .build/assets web/bindata.go