# FME Server CLI

A command line interface for FME Server.

## Description

This is a command line interface that utilizes the FME Server REST API to interact with a running FME Server. It is meant to ease the pain of using the REST API by providing intuitive commands and flags for various operations on an FME Server.

## Getting Started

### Installing

* Simply download the binary for your system from the [releases](https://github.com/safesoftware/fmeflow-cli/releases) page and extract it.
* On Unix systems, you may need to give the file execute permissions (e.g. `chmod +x fmeflow`). You can move the executable to a desired location (e.g. `mv fmeflow /usr/local/bin/fmeflow`)

### Executing program

* Execute the program to get a high level overview of each command
```
fmeflow
```
* Log in to an existing FME Server. It is recommended to generate an API token using the FME Server Web UI initially and use that to log in.
```
fmeflow login https://my-fmeflow.com --token my-token-here
```
* Your token and URL will be saved to a config file located in $HOME/.fmeflow-cli.yaml. Config file location can be overridden with the `--config` flag
* Test your credentials work
```
fmeflow info
```

For full documentation of all commands, see the [Documentation](docs/fmeflow.md).


## Development

* Run while coding:
```
go run main.go
```
* Build binary
```
go build -o fmeflow
```

A great resource for adding new structs to represent JSON returned from FME Server is this [JSON to Go converter](https://mholt.github.io/json-to-go/) which will create a Go struct for you from a JSON sample.
