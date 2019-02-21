CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-readwrite-github
	cp -r reader src/github.com/whosonfirst/go-whosonfirst-readwrite-github/
	cp -r writer src/github.com/whosonfirst/go-whosonfirst-readwrite-github/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-readwrite/..."
	@GOPATH=$(shell pwd) go get -u "github.com/google/go-github/..."

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt reader/*.go
	go fmt writer/*.go

bin: 	self
	GOPATH=$(GOPATH) go build -o bin/wof-github-copy cmd/wof-github-copy.go
	GOPATH=$(GOPATH) go build -o bin/api-reader cmd/api-reader.go