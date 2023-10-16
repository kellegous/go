SHA := $(shell git rev-parse HEAD)

all: bin/golinks

bin/golinks: $(shell find pkg -type f) pkg/web/ui/build
	go build -o $@ ./cmd/golinks

pkg/web/ui/build: node_modules/build $(shell find ui -type f) vite.config.js
	npm run build
	echo $(SHA) > $@

node_modules/build:
	npm install --verbose
	touch $@

# CPP = /usr/bin/cpp -P -undef -Wundef -std=c99 -nostdinc -Wtrigraphs -fdollars-in-identifiers -C -Wno-invalid-pp-token

# SRC = $(shell find web/assets -maxdepth 1 -type f)
# DST = $(patsubst %.scss,%.css,$(patsubst %.ts,%.js,$(subst web/assets,.build/assets,$(SRC))))

# ALL: web/bindata.go

# .build/bin/go-bindata:
# 	GOPATH=$(shell pwd)/.build go get github.com/a-urth/go-bindata/...

# .build/assets:
# 	mkdir -p $@

# .build/assets/%.css: web/assets/%.scss node_modules/build
# 	npx sass --no-source-map --style=compressed $< $@

# .build/assets/%.js: web/assets/%.ts node_modules/build closure-compiler.jar
# 	$(eval TMP := $(shell mktemp))
# 	npx tsc --out $(TMP) --ignoreDeprecations 5.0 $< 
# 	java -jar closure-compiler.jar --js $(TMP) --js_output_file $@
# 	rm -f $(TMP)

# .build/assets/%: web/assets/%
# 	cp $< $@

# node_modules/build:
# 	npm install --verbose
# 	touch $@

# closure-compiler.jar:
# 	curl -L -o $@ https://repo1.maven.org/maven2/com/google/javascript/closure-compiler/v20230802/closure-compiler-v20230802.jar

# web/bindata.go: .build/bin/go-bindata .build/assets $(DST)
# 	$< -o $@ -pkg web -prefix .build/assets -nomemcopy .build/assets/...

# clean:
# 	rm -rf .build/assets web/bindata.go