## fmeserver repositories

List repositories

### Synopsis

Lists repositories on the given FME Server.

```
fmeserver repositories [flags]
```

### Examples

```

	Examples:
	# List all repositories
	fmeserver repositories
	
	# List all repositories owned by the admin user
	fmeserver repositories --owner admin
	
	# List a single repository with the name "Samples"
	fmeserver repositories --name Samples
	
	# Output just the name of all the repositories
	fmeserver repositories --output=custom-columns=NAME:.name --no-headers
	
	# Output all repositories in json format
	fmeserver repositories --json
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
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver](fmeserver.md)	 - A command line interface for interacting with FME Server.
* [fmeserver repositories create](fmeserver_repositories_create.md)	 - Create a repository
* [fmeserver repositories delete](fmeserver_repositories_delete.md)	 - Delete a repository

