GOPATH:=GOPATH=$(shell pwd)/vendor
GOENV:=$(GOPATH)
GOFILES:= "./src/api.go ./src/dbops.go ./main.go subsonic_objects.go"

all: clean godeps build-server build-indexer install

clean:
	rm -rf bin/*

godeps:
	$(GOENV) go get github.com/mattn/go-sqlite3

build-server: $(GOFILES)
	$(GOENV) go build -o bin/ubersonic-server $(GOFILES)

build-indexer: src/indexer.go
	$(GOENV) go build -o bin/ubersonic-indexer ./src/indexer.go

install:
	mkdir -p /opt/ubersonic/bin
	cp bin/ubersonic-server bin/ubersonic-indexer /opt/ubersonic/bin/
