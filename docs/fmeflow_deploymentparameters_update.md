## fmeflow deploymentparameters update

Update a deployment parameter

### Synopsis

Update a deployment parameter.

```
fmeflow deploymentparameters update [flags]
```

### Examples

```

	Examples:
	# Update a deployment parameter with the name "myParam" and the value "myValue"
	fmeflow deploymentparameters update --name myParam --value myValue

```

### Options

```
      --database-type string           The type of the database to use for the database deployment parameter. (Optional)
      --excluded-service stringArray   Service to exclude in the deployment parameter. Can be passed in multiple times if there are multiple Web services to exclude.
  -h, --help                           help for update
      --included-service stringArray   Service to include in the deployment parameter. Can be passed in multiple times if there are multiple Web services to include.
      --name string                    Name of the deployment parameter to update.
      --type string                    Update the type of the parameter. Must be one of text, database, or web. Default is text.
      --value string                   The value to set the deployment parameter to.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow deploymentparameters](fmeflow_deploymentparameters.md)	 - List Deployment Parameters

