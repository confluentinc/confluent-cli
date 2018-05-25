# Confluent Cloud Makefile Includes
This is a set of Makefile include targets that are used in cloud applications.

## Install
Add this repo to your repo with the command:
```shell
git subtree add --prefix mk-include git@github.com:confluentinc/cc-mk-include.git master --squash
```

Then update your makefile like so:

### Go + Docker Service
```make
SERVICE_NAME := example-service
IMAGE_NAME := confluentinc/cc-$(SERVICE_NAME)
MAIN_GO := cmd/server/main.go

include ./mk-include/cc-go-targets.mk
```

### Docker Only Service
```make
IMAGE_NAME := confluentinc/cc-example
MODULE_NAME := example
BASE_IMAGE := 368821881613.dkr.ecr.us-west-2.amazonaws.com/confluentinc/caas-base-alpine
BASE_VERSION := v0.6.1

include ./mk-include/cc-docker-targets.mk
```

## Updating
Once you have the make targets installed, you can update at any time by running

```shell
make update-mk-include
```

## Targets
### Go
TODO: Document go targets

### Docker

TODO: Document docker targets
