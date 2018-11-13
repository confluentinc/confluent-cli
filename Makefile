_empty :=
_space := $(_empty) $(empty)

# Auto bump by default
BUMP ?= auto
# If on master branch bump the minor by default
ifeq ($(RELEASE_BRANCH),master)
DEFAULT_BUMP ?= minor
# Else bump the patch by default
else
DEFAULT_BUMP ?= patch
endif

VERSION := $(shell [ -d .git ] && git describe --tags --always --dirty)
ifneq (,$(findstring dirty,$(VERSION)))
VERSION := $(VERSION)-$(USER)
endif
CLEAN_VERSION := $(shell echo $(VERSION) | grep -Eo '([0-9]+\.){2}[0-9]+')
VERSION_NO_V := $(shell echo $(VERSION) | sed 's,^v,,' )

CI_SKIP ?= [ci skip]

ifeq ($(CLEAN_VERSION),$(_empty))
CLEAN_VERSION := 0.0.0
else
GIT_MESSAGES := $(shell [ -d .git ] && git log --pretty='%s' v$(CLEAN_VERSION)...HEAD | tr '\n' ' ')
endif

# If auto bump enabled, search git messages for bump hash
ifeq ($(BUMP),auto)
_auto_bump_msg := \(auto\)
ifneq (,$(findstring \#major,$(GIT_MESSAGES)))
BUMP := major
else ifneq (,$(findstring \#minor,$(GIT_MESSAGES)))
BUMP := minor
else ifneq (,$(findstring \#patch,$(GIT_MESSAGES)))
BUMP := patch
else
BUMP := $(DEFAULT_BUMP)
endif
endif

# Figure out what the next version should be
split_version := $(subst .,$(_space),$(CLEAN_VERSION))
ifeq ($(BUMP),major)
bump := $(shell expr $(word 1,$(split_version)) + 1)
BUMPED_CLEAN_VERSION := $(bump).0.0
else ifeq ($(BUMP),minor)
bump := $(shell expr $(word 2,$(split_version)) + 1)
BUMPED_CLEAN_VERSION := $(word 1,$(split_version)).$(bump).0
else ifeq ($(BUMP),patch)
bump := $(shell expr $(word 3,$(split_version)) + 1)
BUMPED_CLEAN_VERSION := $(word 1,$(split_version)).$(word 2,$(split_version)).$(bump)
endif

BUMPED_VERSION := v$(BUMPED_CLEAN_VERSION)

RELEASE_SVG := <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="94" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="\#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><clipPath id="a"><rect width="94" height="20" rx="3" fill="\#fff"/></clipPath><g clip-path="url(\#a)"><path fill="\#555" d="M0 0h49v20H0z"/><path fill="\#007ec6" d="M49 0h45v20H49z"/><path fill="url(\#b)" d="M0 0h94v20H0z"/></g><g fill="\#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110"><text x="255" y="150" fill="\#010101" fill-opacity=".3" transform="scale(.1)" textLength="390">release</text><text x="255" y="140" transform="scale(.1)" textLength="390">release</text><text x="705" y="150" fill="\#010101" fill-opacity=".3" transform="scale(.1)" textLength="350">$(BUMPED_VERSION)</text><text x="705" y="140" transform="scale(.1)" textLength="350">$(BUMPED_VERSION)</text></g> </svg>

CCSTRUCTS 						:= $(GOPATH)/src/github.com/confluentinc/cc-structs
ALL_SRC               := $(shell find . -name "*.go" | grep -v -e vendor)
GOLANGCI_LINT_VERSION := 1.12.2

.PHONY: show-version
## Show version variables
show-version:
	@echo version: $(VERSION)
	@echo version no v: $(VERSION_NO_V)
	@echo clean version: $(CLEAN_VERSION)
	@echo version bump: $(BUMP) $(_auto_bump_msg)
	@echo bumped version: $(BUMPED_VERSION)
	@echo bumped clean version: $(BUMPED_CLEAN_VERSION)
	@echo 'release svg: $(RELEASE_SVG)'

.PHONY: tag-release
tag-release:
	# Delete tag if it already exists
	git tag -d $(BUMPED_VERSION) || true
	git push $(GIT_REMOTE_NAME) :$(BUMPED_VERSION) || true
	git tag $(BUMPED_VERSION)
	git push $(GIT_REMOTE_NAME) $(RELEASE_BRANCH) --tags

.PHONY: get-release-image
get-release-image:
	echo '$(RELEASE_SVG)' > release.svg
	git add release.svg

.PHONY: commit-release
commit-release:
	git diff --exit-code --cached --name-status || \
	git commit -m "$(BUMPED_VERSION): $(BUMP) version bump $(CI_SKIP)"

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
binary:
	@GO111MODULE=on gox -os="$(shell go env GOOS)" -arch="$(shell go env GOARCH)" \
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
