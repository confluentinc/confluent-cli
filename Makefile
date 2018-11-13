CCSTRUCTS 						:= $(GOPATH)/src/github.com/confluentinc/cc-structs
ALL_SRC               := $(shell find . -name "*.go" | grep -v -e vendor)
GOLANGCI_LINT_VERSION := 1.12.2

.PHONY: deps
deps:
	@which gox >/dev/null 2>&1 || go get github.com/mitchellh/gox >/dev/null 2>&1
	@which goreleaser >/dev/null 2>&1 || go get github.com/goreleaser/goreleaser >/dev/null 2>&1
	@(golangci-lint --version | grep $(GOLANGCI_LINT_VERSION)) >/dev/null 2>&1 || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v$(GOLANGCI_LINT_VERSION) >/dev/null 2>&1
	@GO111MODULE=on go mod download >/dev/null 2>&1

.PHONY: compile
compile:
	protoc -I shared/connect -I $(CCSTRUCTS) -I $(CCSTRUCTS)/vendor shared/connect/*.proto --gogo_out=plugins=grpc:shared/connect
	protoc -I shared/kafka -I $(CCSTRUCTS) -I $(CCSTRUCTS)/vendor shared/kafka/*.proto --gogo_out=plugins=grpc:shared/kafka
	protoc -I shared/ksql -I $(CCSTRUCTS) -I $(CCSTRUCTS)/vendor shared/ksql/*.proto --gogo_out=plugins=grpc:shared/ksql

.PHONY: install-plugins
install-plugins:
	@GO111MODULE=on go install ./plugin/...

.PHONY: binary
binary: install-plugins
	@GO111MODULE=on go build

.PHONY: dev
dev:
	@gox -os="$(shell go env GOOS)" -arch="$(shell go env GOARCH)" \
	  -output="{{if eq .Dir \"cli\"}}confluent{{else}}{{.Dir}}{{end}}" ./...

.PHONY: release
release: get-release-image commit-release tag-release
	goreleaser

.PHONY: fmt
fmt:
	@gofmt -e -s -l -w $(ALL_SRC)

.PHONY: release-ci
release-ci:
ifeq ($(BRANCH_NAME),master)
	make release
else
	true
endif

.PHONY: lint
lint:
	@GO111MODULE=on golangci-lint run

.PHONY: coverage
coverage:
      ifdef CI
	@echo "" > coverage.txt
	@for d in $$(go list ./... | grep -v vendor); do \
	  GO111MODULE=on go test -v -race -coverprofile=profile.out -covermode=atomic $$d || exit 2; \
	  if [ -f profile.out ]; then \
	    cat profile.out >> coverage.txt; \
	    rm profile.out; \
	  fi; \
	done
      else
	@GO111MODULE=on go test -race -cover $(TEST_ARGS) ./...
      endif

.PHONY: test
test: lint coverage
