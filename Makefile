LIBS := gsm
REV_VAR := gsm.Rev
VERSION_VAR := gsm.Version
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
GOBUILD_VERSION_ARGS := -ldflags "-X $(REV_VAR) $(REPO_REV) -X $(VERSION_VAR) $(REPO_VERSION)"

all: test

test: build
	go test $(GO_TAG_ARGS) -x $(LIBS)

build: deps
	go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) -x $(LIBS)
	go build -o $${GOPATH%%:*}/bin/gsm-server $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) ./gsm-server

deps:
	if [ ! -L $${GOPATH%%:*}/src/gsm ] ; then gvm linkthis ; fi

clean:
	go clean -x $(LIBS) || true
	if [ -d $${GOPATH%%:*}/pkg ] ; then \
		find $${GOPATH%%:*}/pkg -name '*gsm*' -exec rm -v {} \; ; \
	fi

.PHONY: all test build deps clean
