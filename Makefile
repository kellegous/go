ALL: bindata.go

.build/bin/go-bindata:
	GOPATH=$(shell pwd)/.build go get github.com/jteeuwen/go-bindata/...

bindata.go: .build/bin/go-bindata $(wildcard pub/**)
	$< -o $@ -pkg main -prefix pub -nomemcopy pub/...
