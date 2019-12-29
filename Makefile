DIR=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

OLDGOPATH := $(GOPATH)
GOPATH := $(GOPATH)
DATE=$(shell date -u +%Y%m%d.%H%M%S.%Z)
GOGENERATE   = $(shell if [ -f .gogenerate ]; then cat .gogenerate; fi)
TESTPACKETS  = $(shell if [ -f .testpackages ]; then cat .testpackages; fi)
BENCHPACKETS = $(shell if [ -f .benchpackages ]; then cat .benchpackages; fi)
GO111MODULE ?= $(GO111MODULE:on)

default: dep gen

## Generate code by go generate or other utilities
gen:
.PHONY: gen

## Dependence managers
dep:
	@go get -u -v ./...
	@go mod download
	@go mod tidy
	@go mod vendor
.PHONY: dep

test:
	@echo "mode: set" > $(DIR)/coverage.log
	@for PACKET in $(TESTPACKETS); do \
		touch $(DIR)/coverage-tmp.log; \
		GOPATH=${GOPATH} go test -v -covermode=count -coverprofile=$(DIR)/coverage-tmp.log $$PACKET; \
		if [ "$$?" -ne "0" ]; then exit $$?; fi; \
		tail -n +2 $(DIR)/coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> $(DIR)/coverage.log; \
		rm -f $(DIR)/coverage-tmp.log; true; \
	done
.PHONY: test

cover: test
	@GOPATH=${GOPATH} go tool cover -html=$(DIR)/coverage.log
.PHONY: cover

bench:
	@for PACKET in $(BENCHPACKETS); do GOPATH=${GOPATH} go test -race -bench=. -benchmem $$PACKET; done
.PHONY: bench

## Code quality testing
lint:
	@golangci-lint \
		run \
		--enable-all \
		--disable nakedret \
		--disable gochecknoinits \
		--disable goimports \
		--disable maligned \
		./...; \
.PHONY: lint

clean:
	@rm -rf ${DIR}/src; true
	@rm -rf ${DIR}/bin; true
	@rm -rf ${DIR}/pkg; true
	@rm -rf ${DIR}/*.log; true
	@rm -rf ${DIR}/*.lock; true
.PHONY: clean
