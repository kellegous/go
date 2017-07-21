ALL: web/bindata.go

.build/bin/go-bindata:
	GOPATH=$(shell pwd)/.build go get github.com/jteeuwen/go-bindata/...

web/bindata.go: .build/bin/go-bindata $(wildcard web/assets/*)
	$< -o $@ -pkg web -prefix web/assets -nomemcopy web/assets/...
