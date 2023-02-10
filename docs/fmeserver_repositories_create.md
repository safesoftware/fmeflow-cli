## fmeserver repositories create

Create a repository

### Synopsis

Create a new repository.

```
fmeserver repositories create [flags]
```

### Examples

```

	Examples:
	# Create a repository with the name "myRepository" and no description
	fmeserver repositories create --name myRepository
	
	# Output just the name of all the repositories
	fmeserver repositories create --name myRepository --description "This is my new repository"

```

### Options

```
      --description string   Description of the new repository.
  -h, --help                 help for create
      --name string          Name of the repository to create.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver repositories](fmeserver_repositories.md)	 - List repositories

