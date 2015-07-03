ALL: web/bindata.go

.build/bin/go-bindata:
	GOPATH=$(shell pwd)/.build go get github.com/jteeuwen/go-bindata/...

web/bindata.go: .build/bin/go-bindata $(wildcard pub/**/*)
	$< -o $@ -pkg web -prefix pub -nomemcopy pub/...
