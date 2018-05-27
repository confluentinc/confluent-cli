# Confluent CLI

[![Build Status](https://semaphoreci.com/api/v1/projects/accef4bb-d1db-491f-b22e-0d438211c888/1992525/shields_badge.svg)](https://semaphoreci.com/confluent/cli)
![Release](release.svg)
[![codecov](https://codecov.io/gh/confluentinc/cli/branch/master/graph/badge.svg?token=67t1cdciLU)](https://codecov.io/gh/confluentinc/cli)

## Install

Right now, we're publishing to a private S3 bucket. Make sure you have your AWS
engineering creds setup in your `[default]` AWS profile locally, like normal.

To install the CLI itself, run:

    VERSION=12ffed3
    aws s3 cp s3://cloud-confluent-bin/cli/cli-${VERSION}-darwin-amd64 ./cli

To install all components, run:

    VERSION=12ffed3
    for COMPONENT in confluent-kafka-plugin confluent-connect-plugin ; do
        aws s3 cp s3://cloud-confluent-bin/cli/components/${COMPONENT}/${COMPONENT}-${VERSION}-darwin-amd64 ./${COMPONENT}
    done

Note: components must be installed in your `$PATH` for the CLI to pick them up.

## Developing

```
$ make compile-proto
$ go run main.go --help
```

The CLI automatically adds commands when their respective plugins are installed. Enabling the connect
commands by installing the plugins:

```
$ make install-plugins
```

Now you can run:

```
$ go run main.go connect list
```
