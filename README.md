# Confluent CLI

## Developing

```
$ make compile-proto
$ go run main.go --help
```

The CLI automatically adds commands when their respctive plugins are installed. Enabling the connect
commands by installing the plugig:

```
$ make install-plugins
```

Now you can run:

```
$ go run main.go connect list
```
