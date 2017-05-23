# Confluent Platform CLI
A CLI to start and manage Confluent Platform from command line.

## Installation

* Download and install [Confluent OSS](https://www.confluent.io/download/)

* Checkout *confluent-cli* by running:

    ```bash
    git clone git@github.com:confluentinc/confluent-cli.git
    ```

* Set *CONFLUENT_HOME* environment variable to point to the location of Confluent OSS. For instance:

    ```bash
    export CONFLUENT_HOME=/usr/local/confluent-3.2.0
    ```

* Install *confluent-cli*:

    ```bash
    cd confluent-cli; make install
    ```

## Usage
To get a list of available commands, run:

```bash
export PATH=${CONFLUENT_HOME}/bin:${PATH};
confluent help
```

Examples:

```bash
confluent start

confluent status

confluent stop

confluent current

confluent destroy
```
