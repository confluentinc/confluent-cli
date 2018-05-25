_empty := 
_space := $(_empty) $(empty)
_include_prefix := mk-include/

include ./$(_include_prefix)cc-semver.mk

# Base Stuff
BASE_IMAGE ?=
BASE_VERSION ?=

# Image Variables
IMAGE_NAME ?= unknown
IMAGE_VERSION ?= $(VERSION)
BUILD_TAG ?= $(IMAGE_NAME):$(IMAGE_VERSION)

# ECR Stuff
DOCKER_REPO ?= 368821881613.dkr.ecr.us-west-2.amazonaws.com
AWS_PROFILE ?= default
AWS_REGION ?= us-west-2
LIFECYCLE_POLICY ?= '{"rules":[{"rulePriority":10,"description":"keeps 50 latest tagged images","selection":{"tagStatus":"tagged","countType":"imageCountMoreThan","countNumber":50,"tagPrefixList":["v"]},"action":{"type":"expire"}},{"rulePriority":20,"description":"keeps 5 latest untagged images","selection":{"tagStatus":"untagged","countType":"imageCountMoreThan","countNumber":5},"action":{"type":"expire"}},{"rulePriority":30,"description":"keeps latest 20 numeric-tagged images","selection":{"tagStatus":"tagged","countType":"imageCountMoreThan","tagPrefixList":["0","1","2","3","4","5","6","7","8","9"],"countNumber":20},"action":{"type":"expire"}},{"rulePriority":40,"description":"keeps latest 20 a-f tagged images","selection":{"tagStatus":"tagged","countType":"imageCountMoreThan","tagPrefixList":["a","b","c","d","e","f"],"countNumber":20},"action":{"type":"expire"}}]}'

# Terraform Variables
MODULE_NAME ?= $(IMAGE_NAME)
define BUMPED_IMAGE_VERSION_OVERRIDE
{ "variable": { "image_version": { "default": "$(BUMPED_VERSION)" } } }
endef

ifeq ($(IMAGE_NAME),unknown)
$(error IMAGE_NAME must be set)
endif

_docker_release_target ?= release-docker

include ./$(_include_prefix)cc-terraform-targets.mk

show-args: show-version
	@echo 'IMAGE_NAME: $(IMAGE_NAME)'
	@echo 'IMAGE_VERSION: $(IMAGE_VERSION)'
	@echo 'MODULE_NAME: $(MODULE_NAME)'
	@echo 'DOCKER_REPO: $(DOCKER_REPO)'
	@echo 'BUILD_TAG: $(BUILD_TAG)'
	@echo 'DOWNSTREAM_TF_CONSUMERS: $(DOWNSTREAM_TF_CONSUMERS)'

$(HOME)/.netrc:
	$(eval user := $(shell bash -c 'read -p "GitHub Username: " user; echo $$user'))
	$(eval pass := $(shell bash -c 'read -p "GitHub Password: " pass; echo $$pass'))
	@printf "machine github.com\n\tlogin $(user)\n\tpassword $(pass)" > $(HOME)/.netrc

.netrc: $(HOME)/.netrc
	cp $(HOME)/.netrc .netrc

create-repo:
	aws ecr describe-repositories --region us-west-2 --repository-name $(IMAGE_NAME) || aws ecr create-repository --region us-west-2 --repository-name $(IMAGE_NAME)
	aws ecr put-lifecycle-policy --region us-west-2 --repository-name $(IMAGE_NAME) --lifecycle-policy-text $(LIFECYCLE_POLICY) || echo "Failed to put lifecycle policy on $(IMAGE_NAME) repo"

pull-base: repo-login
ifneq ($(BASE_IMAGE),$(_emmpty))
	docker image ls -f reference="$(BASE_IMAGE):$(BASE_VERSION)" | grep -Eq "$(BASE_IMAGE)[ ]*$(BASE_VERSION)" || \
		docker pull $(BASE_IMAGE):$(BASE_VERSION)
endif

build-docker: .netrc pull-base
	docker build --no-cache --build-arg version=$(IMAGE_VERSION) -t $(IMAGE_NAME):$(IMAGE_VERSION) .
	rm -f .netrc

push-docker: create-repo push-docker-latest push-docker-version

push-docker-latest: tag-docker-latest
	@echo 'push latest to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(IMAGE_NAME):latest

push-docker-version: tag-docker-version
	@echo 'push $(IMAGE_VERSION) to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(IMAGE_NAME):$(IMAGE_VERSION)

tag-docker: tag-docker-latest tag-docker-version

tag-docker-latest:
	@echo 'create docker tag latest'
	docker tag $(IMAGE_NAME):$(IMAGE_VERSION) $(DOCKER_REPO)/$(IMAGE_NAME):latest

tag-docker-version:
	@echo 'create docker tag $(IMAGE_VERSION)'
	docker tag $(IMAGE_NAME):$(IMAGE_VERSION) $(DOCKER_REPO)/$(IMAGE_NAME):$(IMAGE_VERSION)

clean-images:
	docker images -q -f label=io.confluent.caas=true -f reference='*$(IMAGE_NAME)' | uniq | xargs docker rmi -f

clean-all:
	docker images -q -f label=io.confluent.caas=true | uniq | xargs docker rmi -f

clean-terraform:
	rm -rf terraform/deployments/minikube/*.tfstate* terraform/deployments/minikube/.terraform*

set-tf-bumped-version:
	test -d terraform \
		&& (echo '$(BUMPED_IMAGE_VERSION_OVERRIDE)' > terraform/modules/$(MODULE_NAME)/image_version_override.tf.json && \
			git add terraform/modules/$(MODULE_NAME)/image_version_override.tf.json) \
		|| true

set-node-bumped-version:
	test -f package.json \
		&& (npm version $(BUMPED_VERSION) --git-tag-version=false &&\
			git add package.json) \
		|| true

deploy-terraform:
	cd terraform/deployments/minikube && \
		terraform init && \
		terraform apply -var "image_version=$(VERSION)"

release-docker: build-docker push-docker bump-downstream-tf-consumers

release: set-tf-bumped-version set-node-bumped-version get-release-image commit-release tag-release
	make $(_docker_release_target)

release-ci:
ifeq ($(BRANCH_NAME),master)
	make release
else
	true
endif

repo-login:
	@eval "$(shell aws ecr get-login --no-include-email --region us-west-2 --profile default)"
