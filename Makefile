LIBS := gswat
REV_VAR := gswat.Rev
VERSION_VAR := gswat.Version
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
GOBUILD_VERSION_ARGS := -ldflags "-X $(REV_VAR) $(REPO_REV) -X $(VERSION_VAR) $(REPO_VERSION)"

all: test

test: build
	go test $(GO_TAG_ARGS) -x $(LIBS)

build: deps
	go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) -x $(LIBS)
	go build -o $${GOPATH%%:*}/bin/gswat-server $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) ./gswat-server

deps:
	if [ ! -L $${GOPATH%%:*}/src/gswat ] ; then gvm linkthis ; fi

clean:
	go clean -x $(LIBS) || true
	if [ -d $${GOPATH%%:*}/pkg ] ; then \
		find $${GOPATH%%:*}/pkg -name '*gswat*' -exec rm -v {} \; ; \
	fi

.PHONY: all test build deps clean
