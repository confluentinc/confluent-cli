# Confluent Platform CLI
A CLI to start and manage Confluent Platform from command line.

## Installation

* Download and install [Confluent OSS](https://www.confluent.io/download/)

* Checkout *confluent-cli* by running:

    ```bash
    $ git clone git@github.com:confluentinc/confluent-cli.git
    ```

* Set *CONFLUENT_HOME* environment variable to point to the location of Confluent OSS. For instance:

    ```bash
    $ export CONFLUENT_HOME=/usr/local/confluent-3.3.0
    ```

* Install *confluent-cli*:

    ```bash
    $ cd confluent-cli; make install
    ```

## Usage
To get a list of available commands, run:

```bash
$ export PATH=${CONFLUENT_HOME}/bin:${PATH};
$ confluent help
```

Examples:

* Start all the services!
```bash
$ confluent start
```

* Retrieve their status:
```bash
$ confluent status
```

* Open the log file of a service:
```bash
$ confluent log connect
```

* Access runtime stats of a service:
```bash
$ confluent top kafka
```

* Discover the availabe Connect plugins:
```bash
$ confluent list plugins
```

* or list the predefined connector names:
```bash
$ confluent list connectors
```

* Load a couple connectors:
```bash
$ confluent load file-source
$ confluent load file-sink
```

* Get a list with the currently loaded connectors:
```bash
$ confluent status connectors
```

* Check the status of a loaded connector:
```bash
$ confluent status file-source
```

* Read the configuration of a connector:
```bash
$ confluent config file-source
```

* Reconfigure a connector:
```bash
$ confluent config file-source -d ./updated-file-source-config.json
```

* or reconfigure using a properties file:
```bash
$ confluent config file-source -d ./updated-file-source-config.properties
```

* Figure out where the data and the logs of the current confluent run are stored:
```bash
$ confluent current
```

* Unload a specific connector:
```bash
$ confluent unload file-sink
```

* Stop the services:
```bash
$ confluent stop
```

* Start on a clean slate next time (deletes data and logs of a confluent run):
```bash
$ confluent destroy
```

Set CONFLUENT_CURRENT if you want to use a top directory for confluent runs other than your platform's tmp directory.

```bash
$ cd $CONFLUENT_HOME
$ mkdir -p var
$ export CONFLUENT_CURRENT=${CONFLUENT_HOME}/var
$ confluent current
```
