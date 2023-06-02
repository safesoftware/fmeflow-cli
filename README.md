# FME Flow CLI

A command line interface for FME Flow.

## Description

This is a command line interface that utilizes the FME Flow REST API to interact with a running FME Flow instance. It is meant to ease the pain of using the REST API by providing intuitive commands and flags for various operations on an FME Flow instance.

## Getting Started

### Installing

* Simply download the binary for your system from the [releases](https://github.com/safesoftware/fmeflow-cli/releases) page and extract it.
* On Unix systems, you may need to give the file execute permissions (e.g. `chmod +x fmeflow`). You can move the executable to a desired location (e.g. `mv fmeflow /usr/local/bin/fmeflow`)

### Executing program

* Execute the program to get a high level overview of each command
```
fmeflow
```
* Log in to an existing FME Flow instance. It is recommended to generate an API token using the FME Flow Web UI initially and use that to log in.
```
fmeflow login https://my-fmeflow.com --token my-token-here
```
* Your token and URL will be saved to a config file located in $HOME/.fmeflow-cli.yaml. Config file location can be overridden with the `--config` flag
* Test your credentials work
```
fmeflow info
```

For full documentation of all commands, see the [Documentation](docs/fmeflow.md).

## Supported Versions of FME Flow

This CLI has been written with backwards compatibilty in mind. Officially this will support FME Flow 2022.2 and later. However, we have tested back to FME Flow 2019 and are able to log in and run commands. Not all commands are guaranteed to work on builds before FME Flow 2022.2.

## Development

* Run while coding:
```
go run main.go
```
* Build binary
```
go build -o fmeflow
```

A great resource for adding new structs to represent JSON returned from FME Flow is this [JSON to Go converter](https://mholt.github.io/json-to-go/) which will create a Go struct for you from a JSON sample.
