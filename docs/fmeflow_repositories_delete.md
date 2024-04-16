## fmeflow repositories delete

Delete a repository.

### Synopsis

Delete a repository.

```
fmeflow repositories delete [flags]
```

### Examples

```

  # Delete a repository with the name "myRepository"
  fmeflow repositories delete --name myRepository
	
  # Delete a repository with the name "myRepository" and no confirmation
  fmeflow repositories delete --name myRepository --no-prompt

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

* [fmeflow repositories](fmeflow_repositories.md)	 - List repositories

