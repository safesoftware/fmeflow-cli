## fmeflow deploymentparameters

List Deployment Parameters

### Synopsis

Lists Deployment Parameters on the given FME Server.

```
fmeflow deploymentparameters [flags]
```

### Examples

```

  # List all deployment parameters
  fmeflow deploymentparameters
	
  # List a single deployment parameter
  fmeflow deploymentparameters --name testParameter
	
  # Output all deploymentparameters in json format
  fmeflow deploymentparameters --json
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
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Server.
* [fmeflow deploymentparameters create](fmeflow_deploymentparameters_create.md)	 - Create a deployment parameter
* [fmeflow deploymentparameters delete](fmeflow_deploymentparameters_delete.md)	 - Delete a deployment parameter
* [fmeflow deploymentparameters update](fmeflow_deploymentparameters_update.md)	 - Update a deployment parameter

