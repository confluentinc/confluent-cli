_empty := 
_space := $(_empty) $(empty)

# Default to patch release
BUMP ?= auto
DEFAULT_BUMP ?= minor

VERSION := $(shell [ -d .git ] && git describe --tags --always --dirty)
CLEAN_VERSION := $(shell [ -d .git ] && git describe --tags --always --dirty | grep -Eo '([0-9]+\.){2}[0-9]+')

ifeq ($(CLEAN_VERSION),$(_empty))
CLEAN_VERSION := 0.0.0
else
GIT_MESSAGES := $(shell git log --pretty='%s' v$(CLEAN_VERSION)...HEAD | tr '\n' ' ')
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
BUMPED_VERSION := v$(bump).0.0
else ifeq ($(BUMP),minor)
bump := $(shell expr $(word 2,$(split_version)) + 1)
BUMPED_VERSION := v$(word 1,$(split_version)).$(bump).0
else ifeq ($(BUMP),patch)
bump := $(shell expr $(word 3,$(split_version)) + 1)
BUMPED_VERSION := v$(word 1,$(split_version)).$(word 2,$(split_version)).$(bump)
endif

RELEASE_SVG := <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="94" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="\#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><clipPath id="a"><rect width="94" height="20" rx="3" fill="\#fff"/></clipPath><g clip-path="url(\#a)"><path fill="\#555" d="M0 0h49v20H0z"/><path fill="\#007ec6" d="M49 0h45v20H49z"/><path fill="url(\#b)" d="M0 0h94v20H0z"/></g><g fill="\#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110"><text x="255" y="150" fill="\#010101" fill-opacity=".3" transform="scale(.1)" textLength="390">release</text><text x="255" y="140" transform="scale(.1)" textLength="390">release</text><text x="705" y="150" fill="\#010101" fill-opacity=".3" transform="scale(.1)" textLength="350">$(BUMPED_VERSION)</text><text x="705" y="140" transform="scale(.1)" textLength="350">$(BUMPED_VERSION)</text></g> </svg>

show-version:
	@echo version: $(VERSION)
	@echo clean version: $(CLEAN_VERSION)
	@echo version bump: $(BUMP) $(_auto_bump_msg)
	@echo bumped version: $(BUMPED_VERSION)
	@echo 'release svg: $(RELEASE_SVG)'

tag-release:
	git tag $(BUMPED_VERSION)
	git push origin master --tags

get-release-image:
	echo '$(RELEASE_SVG)' > release.svg
	git add release.svg

commit-release:
	git diff --exit-code --cached --name-status || \
	git commit -m "$(BUMPED_VERSION): $(BUMP) version bump [ci skip]"
