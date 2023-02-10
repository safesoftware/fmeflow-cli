## fmeserver workspaces

List workspaces by repository

### Synopsis

Lists workspaces on the given FME Server in the repository.

```
fmeserver workspaces [flags]
```

### Examples

```

	Examples:
	# List all workspaces in Samples repository
	fmeserver workspaces --repository Samples
	
	# List all workspaces in the Samples repository and output it in json
	fmeserver workspaces --repository Samples --json
	
	# List all workspaces in the Samples repository with custom columns showing the last publish date and number of times run
	fmeserver workspaces --repository Samples --output="custom-columns=NAME:.name,PUBLISH DATE:.lastPublishDate,TOTAL RUNS:.totalRuns"
	
	# Get information on a single workspace 
	fmeserver workspaces --repository Samples --name austinApartments.fmw
	
	# Get the name, source format, and destination format for this workspace
	fmeserver workspaces --repository Samples --name austinApartments.fmw --output=custom-columns=NAME:.name,SOURCE:.datasets.source[*].format,DEST:.datasets.destination[*].format
```

### Options

```
      --filter-string string   Specify the output type. Should be one of table, json, or custom-columns. Only usable with V4 API.
  -h, --help                   help for workspaces
      --name string            If specified, get details about a specific workspace
      --no-headers             Don't print column headers
  -o, --output string          Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --repository string      Name of repository to list workspaces in.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver](fmeserver.md)	 - A command line interface for interacting with FME Server.

