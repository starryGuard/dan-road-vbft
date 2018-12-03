GOFMT=gofmt
GC=go build
VERSION := $(shell git describe --abbrev=4 --always --tags)
BUILD_NODE_PAR = -ldflags "-X dan-road-vbft/common/config.Version=$(VERSION)" #-race

ARCH=$(shell uname -m)
DBUILD=docker build
DRUN=docker run
DOCKER_TAG=$(ARCH)-$(VERSION)
DOCKER_NS ?= dan

SRC_FILES = $(shell git ls-files | grep -e .go$ | grep -v _test.go)

dan-road-vbft: $(SRC_FILES)
	$(GC)  $(BUILD_NODE_PAR) -o dan-road-vbft main.go


format:
	$(GOFMT) -w main.go

docker/payload: docker/build/bin/cvbft docker/Dockerfile
	@echo "Building cvbft payload"
	@mkdir -p $@
	@cp docker/Dockerfile $@
	@cp docker/build/bin/cvbft $@
	@touch $@

docker/build/bin/%: Makefile
	@echo "Building cvbft in docker"
	@mkdir -p docker/build/bin docker/build/pkg
	@$(DRUN) --rm \
		-v $(abspath docker/build/bin):/go/bin \
		-v $(abspath docker/build/pkg):/go/pkg \
		-v $(GOPATH)/src:/go/src \
		-w /go/src/dan-road-vbft \
		golang:1.9.5-stretch \
		$(GC)  $(BUILD_NODE_PAR) -o docker/build/bin/cvbft main.go
	@touch $@

docker: Makefile docker/payload docker/Dockerfile
	@echo "Building cvbft docker"
	@$(DBUILD) -t $(DOCKER_NS)/cvbft docker/payload
	@docker tag $(DOCKER_NS)/cvbft $(DOCKER_NS)/cvbft:$(DOCKER_TAG)
	@touch $@
