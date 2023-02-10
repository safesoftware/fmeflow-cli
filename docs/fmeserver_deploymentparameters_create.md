## fmeserver deploymentparameters create

Create a deployment parameter

### Synopsis

Create a deployment parameter.

```
fmeserver deploymentparameters create [flags]
```

### Examples

```

	Examples:
	# Create a deployment parameter with the name "myParam" and the value "myValue"
	fmeserver deploymentparameter create --name myParam --value myValue

```

### Options

```
  -h, --help           help for create
      --name string    Name of the deployment parameter to create.
      --value string   The value to set the deployment parameter to.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver deploymentparameters](fmeserver_deploymentparameters.md)	 - List Deployment Parameters

