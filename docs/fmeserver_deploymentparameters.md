## fmeserver deploymentparameters

List Deployment Parameters

### Synopsis

Lists Deployment Parameters on the given FME Server.

```
fmeserver deploymentparameters [flags]
```

### Examples

```

	Examples:
	# List all deployment parameters
	fmeserver deploymentparameters
	
	# List a single deployment parameter
	fmeserver deploymentparameters --name testParameter
	
	# Output all deploymentparameters in json format
	fmeserver deploymentparameters --json
```

### Options

```
  -h, --help            help for deploymentparameters
      --name string     If specified, only the repository with that name will be returned
      --no-headers      Don't print column headers
  -o, --output string   Specify the output type. Should be one of table, json, or custom-columns (default "table")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver](fmeserver.md)	 - A command line interface for interacting with FME Server.
* [fmeserver deploymentparameters create](fmeserver_deploymentparameters_create.md)	 - Create a deployment parameter
* [fmeserver deploymentparameters delete](fmeserver_deploymentparameters_delete.md)	 - Delete a deployment parameter
* [fmeserver deploymentparameters update](fmeserver_deploymentparameters_update.md)	 - Update a deployment parameter

