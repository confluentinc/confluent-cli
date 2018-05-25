# Downstream TF Consumer Format: (replace : with ^)
#	GIT_URI//TF_FILE_PATH//MODULE_NAME
# Exampe: git@github.com^confluentinc/cc-terraform-module-mothership.git//internal-resources/main.tf//gateway-service
#
# Note that this is a list so list all downstream consumers here
#
_noop :=
_space := $(_noop) $(_noop)
_include_prefix ?= mk-include/
_tf_include_semver ?= false
_tf_include_base ?= false

ifeq ($(_tf_include_semver),true)
include ./$(_include_prefix)cc-semver.mk
endif

ifeq ($(_tf_include_base),true)
include ./$(_include_prefix)cc-base.mk
endif

DOWNSTREAM_TF_CONSUMERS ?=

show-downstream-tf-consumers:
	@echo 'Downstream TF Consumers:'
	@$(foreach consumer,$(DOWNSTREAM_TF_CONSUMERS),echo "  $(subst ^,:,$(consumer))";)

bump-downstream-tf-consumers: $(DOWNSTREAM_TF_CONSUMERS)

$(DOWNSTREAM_TF_CONSUMERS):
	$(eval consumer := $(subst ^,:,$@))
	$(eval split_consumer := $(subst //,$(_space),$(consumer)))
	$(eval git_uri := $(word 1, $(split_consumer)))
	$(eval tf_path := $(word 2, $(split_consumer)))
	$(eval tf_module := $(word 3, $(split_consumer)))
	$(eval repo_name := $(basename $(notdir $(git_uri))))
	git clone $(git_uri)
	vim -Ec '/module "$(tf_module)" {/|/source\s\+=/|s/?ref=v[0-9.]\+/?ref=v$(CLEAN_VERSION)/|p|x' "./$(repo_name)/$(tf_path)"
	git -C "./$(repo_name)" add "$(tf_path)"
	git -C "./$(repo_name)" diff --exit-code --cached --name-status || \
		(git -C "./$(repo_name)" commit -m "TF Module Bump: $(tf_path)//$(tf_module) to $(CLEAN_VERSION)" && \
		 git -C "./$(repo_name)" push origin master)
	rm -rf $(repo_name)

# Release target for tf only modules
release-tf: get-release-image commit-release tag-release
	make bump-downstream-tf-consumers
