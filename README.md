# FME Server CLI

A command line interface for FME Server.

## Description

This is a command line interface that utilizes the FME Server REST API to interact with a running FME Server. It is meant to ease the pain of using the REST API by providing intuitive commands and flags for various operations on an FME Server.

## Getting Started

### Installing

* Simply download the binary for your system from the [releases](https://github.com/safesoftware/fmeserver-cli/releases) page.

### Executing program

* Execute the program to get a high level overview of each command
```
fmeserver
```
* Log in to an existing FME Server. It is recommended to generate an API token using the FME Server Web UI initially and use that to log in.
```
fmeserver login https://my-fmeserver.com --token my-token-here
```
* Your token and URL will be saved to a config file located in $HOME/.fmeserver-cli.yaml. Config file location can be overridden with the `--config` flag
* Test your credentials work
```
fmeserver info
```

## Development

* `cobra-cli` will be needed to add new commands
```
go install github.com/spf13/cobra-cli@latest
```
* Run while coding:
```
go run main.go
```
* Build binary
```
go build -o fmeserver
```
* Add a new command
```
cobra-cli add new-command
```
More details [here](https://github.com/spf13/cobra-cli/blob/main/README.md)

## Acknowledgments

* Created using [cobra](https://github.com/spf13/cobra)