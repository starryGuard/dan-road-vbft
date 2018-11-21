GOFMT=gofmt
GC=go build
VERSION := $(shell git describe --abbrev=4 --always --tags)
BUILD_NODE_PAR = -ldflags "-X github.com/ontio/ontology/common/config.Version=$(VERSION)" #-race

ARCH=$(shell uname -m)
DBUILD=docker build
DRUN=docker run
DOCKER_TAG=$(ARCH)-$(VERSION)

SRC_FILES = $(shell git ls-files | grep -e .go$ | grep -v _test.go)
TOOLS=./tools
ABI=$(TOOLS)/abi
NATIVE_ABI_SCRIPT=./cmd/abi/native_abi_script

dan-road-vbft: $(SRC_FILES)
	$(GC)  $(BUILD_NODE_PAR) -o dan-road-vbft main.go

docker/payload: docker/build/bin/dan-road-vbft docker/Dockerfile
	@mkdir -p $@
	@cp docker/Dockerfile $@
	@cp docker/build/bin/dan-road-vbft $@
	@touch $@

docker/build/bin/%: Makefile
	@mkdir -p docker/build/bin docker/build/pkg
	@$(DRUN) --rm \
		-v $(abspath docker/build/bin):/go/bin \
		-v $(abspath docker/build/pkg):/go/pkg \
		-v $(GOPATH)/src:/go/src \
		-w /go/src/dan-road/vbft \
		golang:1.9.5-stretch \
		$(GC)  $(BUILD_NODE_PAR) -o docker/build/bin/dan-road-vbft main.go
	@touch $@

docker: Makefile docker/payload docker/Dockerfile
	@$(DBUILD) -t /dan-road-vbft docker/payload
	@docker tag /dan-road-vbft /dan-road-vbft:$(DOCKER_TAG)
	@touch $@
