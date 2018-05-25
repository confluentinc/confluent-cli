IMAGE_NAME ?=
SERVICE_NAME ?=
MAIN_GO ?=

MODULE_NAME ?= $(SERVICE_NAME)
BASE_IMAGE ?= 368821881613.dkr.ecr.us-west-2.amazonaws.com/confluentinc/cc-service-base
BASE_VERSION ?= 1.9

_include_prefix := mk-include/
include ./$(_include_prefix)cc-docker-targets.mk
include ./$(_include_prefix)cc-base.mk

GO_LDFLAGS := "-X main.version=$(VERSION)"

clean: clean-images
	rm -rf $(SERVICE_NAME) .netrc

vet:
	@go list ./... | grep -v vendor | xargs go vet

deps:
	which dep 2>/dev/null || go get -u github.com/golang/dep/cmd/dep
	dep ensure $(ARGS)

build:
	go build -o $(SERVICE_NAME) -ldflags $(GO_LDFLAGS) $(MAIN_GO)

install:
	go build -o $(GOBIN)/$(SERVICE_NAME) $(MAIN_GO)

run: deps
	go run $(MAIN_GO)

test:
	go test -v -cover $(TEST_ARGS) ./...

generate:
	go generate

seed-local-mothership:
	psql -d postgres -c 'DROP DATABASE IF EXISTS mothership;'
	psql -d postgres -c 'CREATE DATABASE mothership;'
	psql -d mothership -f mk-include/seed-db/mothership-seed.sql
