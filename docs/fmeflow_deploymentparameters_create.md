## fmeflow deploymentparameters create

Create a deployment parameter

### Synopsis

Create a deployment parameter.

```
fmeflow deploymentparameters create [flags]
```

### Examples

```

	Examples:
	# Create a deployment parameter with the name "myParam" and the value "myValue"
	fmeflow deploymentparameters create --name myParam --value myValue

```

### Options

```
      --database-type string           The type of the database to use for the database deployment parameter. (Optional)
      --excluded-service stringArray   Service to exclude in the deployment parameter. Can be passed in multiple times if there are multiple Web services to exclude.
  -h, --help                           help for create
      --included-service stringArray   Service to include in the deployment parameter. Can be passed in multiple times if there are multiple Web services to include.
      --name string                    Name of the deployment parameter to create.
      --type string                    Type of parameter to create. Must be one of text, database, or web. Default is text.
      --value string                   The value to set the deployment parameter to. (Optional)
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow deploymentparameters](fmeflow_deploymentparameters.md)	 - List, Create, Update and Delete Deployment Parameters

