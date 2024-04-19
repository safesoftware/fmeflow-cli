## fmeflow repositories

List, Create and Delete repositories

### Synopsis

Lists repositories on the given FME Server. Pass in a name to get information on a specific repository. Use the subcommands to create or delete repositories.

```
fmeflow repositories [flags]
```

### Examples

```

  # List all repositories
  fmeflow repositories
	
  # List all repositories owned by the admin user
  fmeflow repositories --owner admin
	
  # List a single repository with the name "Samples"
  fmeflow repositories --name Samples
	
  # Output just the name of all the repositories
  fmeflow repositories --output=custom-columns=NAME:.name --no-headers
	
  # Output all repositories in json format
  fmeflow repositories --json
```

### Options

```
      --filter-string string   Specify the output type. Should be one of table, json, or custom-columns. Only usable with V4 API.
  -h, --help                   help for repositories
      --name string            If specified, only the repository with that name will be returned
      --no-headers             Don't print column headers
  -o, --output string          Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --owner string           If specified, only repositories owned by the specified user uuid will be returned. With the V3 API, set this to the user name.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Flow.
* [fmeflow repositories create](fmeflow_repositories_create.md)	 - Create a new repository.
* [fmeflow repositories delete](fmeflow_repositories_delete.md)	 - Delete a repository.

