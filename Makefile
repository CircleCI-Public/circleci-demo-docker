CIRCLE_BUILD_NUM ?= 0
TAG = 0.0.$(CIRCLE_BUILD_NUM)-$(shell git rev-parse --short HEAD)

GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

default: build

workdir:
	mkdir -p workdir

build: workdir/service

workdir/service: $(GOFILES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o workdir/service .

test: test-all

test-all:
	@go test -v $(GOPACKAGES)
