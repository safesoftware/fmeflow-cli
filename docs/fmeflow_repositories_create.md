## fmeflow repositories create

Create a new repository.

### Synopsis

Create a new repository.

```
fmeflow repositories create [flags]
```

### Examples

```

  # Create a repository with the name "myRepository" and no description
  fmeflow repositories create --name myRepository
	
  # Output just the name of all the repositories
  fmeflow repositories create --name myRepository --description "This is my new repository"

```

### Options

```
      --description string   Description of the new repository.
  -h, --help                 help for create
      --name string          Name of the repository to create.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow repositories](fmeflow_repositories.md)	 - List repositories

