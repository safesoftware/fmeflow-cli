## fmeflow deploymentparameters delete

Delete a deployment parameter

### Synopsis

Delete a deployment parameter.

```
fmeflow deploymentparameters delete [flags]
```

### Examples

```

	Examples:
	# Delete adeployment parameter with the name "myParam"
	fmeflow deploymentparameters delete --name myParam
	
	# Delete a repository with the name "myRepository" and no confirmation
	fmeflow deploymentparameters delete --name myParam --no-prompt

```

### Options

```
  -h, --help          help for delete
      --name string   Name of the repository to create.
  -y, --no-prompt     Description of the new repository.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow deploymentparameters](fmeflow_deploymentparameters.md)	 - List Deployment Parameters

