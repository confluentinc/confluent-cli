include ./mk-include/cc-semver.mk

CCSTRUCTS := $(GOPATH)/src/github.com/confluentinc/cc-structs

.PHONY: deps
deps:
	@which dep >/dev/null 2>&1 || go get github.com/golang/dep/cmd/dep
	@which gometalinter >/dev/null 2>&1 || ( go get github.com/alecthomas/gometalinter && gometalinter --install &> /dev/null )
	@which gox >/dev/null 2>&1 || go get github.com/mitchellh/gox
	@which goreleaser >/dev/null 2>&1 || ( curl -s https://api.github.com/repos/goreleaser/goreleaser/releases/latest | grep "browser_download_url" | grep -i $(shell go env GOOS) | grep -i $(shell go env GOARCH | sed -e 's/amd64/x86_64/') | cut -d : -f 2,3 | tr -d \" | wget -O goreleaser.tar.gz -qi - && tar -xzvf goreleaser.tar.gz goreleaser && chmod +x goreleaser && mv goreleaser $(GOPATH)/bin/goreleaser)
	dep ensure $(ARGS)

.PHONY: compile-proto
compile-proto:
	protoc -I shared/connect -I $(CCSTRUCTS) -I $(CCSTRUCTS)/vendor shared/connect/*.proto --gogo_out=plugins=grpc:shared/connect
	protoc -I shared/kafka -I $(CCSTRUCTS) -I $(CCSTRUCTS)/vendor shared/kafka/*.proto --gogo_out=plugins=grpc:shared/kafka

.PHONY: install-plugins
install-plugins:
	go install ./plugin/...

.PHONY: dev
dev:
	@gox -os="$(shell go env GOOS)" -arch="$(shell go env GOARCH)" \
	  -output="{{if eq .Dir \"cli\"}}confluent{{else}}{{.Dir}}{{end}}" ./...

.PHONY: release
release: get-release-image commit-release tag-release
	goreleaser

.PHONY: release-ci
release-ci:
ifeq ($(BRANCH_NAME),master)
	make release
else
	true
endif

.PHONY: lint
lint:
	gometalinter ./... --vendor

.PHONY: coverage
coverage:
      ifdef CI
	@echo "" > coverage.txt
	@for d in $$(go list ./... | grep -v vendor); do \
	  go test -v -race -coverprofile=profile.out -covermode=atomic $$d || exit 2; \
	  if [ -f profile.out ]; then \
	    cat profile.out >> coverage.txt; \
	    rm profile.out; \
	  fi; \
	done
      else
	@go test -race -cover $(TEST_ARGS) ./...
      endif

.PHONY: test
test: lint coverage

.PHONY: clean
clean:
	rm $(PROTO)/*.pb.go
