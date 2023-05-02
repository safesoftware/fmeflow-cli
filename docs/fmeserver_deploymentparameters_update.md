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
  -h, --help           help for update
      --name string    Name of the deployment parameter to update.
      --value string   The value to set the deployment parameter to.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow deploymentparameters](fmeflow_deploymentparameters.md)	 - List Deployment Parameters

