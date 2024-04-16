## fmeflow workspaces

List workspaces.

### Synopsis

Lists workspaces that exist on the FME Server. Filter by repository, specify a name to retrieve a specific workspace, or specify a filter string to narrow down by name or title.

```
fmeflow workspaces [flags]
```

### Examples

```

  # List all workspaces on the FME Server
  fmeflow workspaces
	
  # List all workspaces in Samples repository
  fmeflow workspaces --repository Samples
	
  # List all workspaces in the Samples repository and output it in json
  fmeflow workspaces --repository Samples --json
	
  # List all workspaces in the Samples repository with custom columns showing the last publish date and number of times run
  fmeflow workspaces --repository Samples --output="custom-columns=NAME:.name,PUBLISH DATE:.lastPublishDate,TOTAL RUNS:.totalRuns"
	
  # Get information on a single workspace 
  fmeflow workspaces --repository Samples --name austinApartments.fmw
	
  # Get the name, source format, and destination format for this workspace
  fmeflow workspaces --repository Samples --name austinApartments.fmw --output=custom-columns=NAME:.name,SOURCE:.datasets.source[*].format,DEST:.datasets.destination[*].format
```

### Options

```
      --filter-string string   If specified, only workspaces with a matching name or title will be returned. Only usable with V4 API.
  -h, --help                   help for workspaces
      --name string            If specified, get details about a specific workspace
      --no-headers             Don't print column headers
  -o, --output string          Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --repository string      Name of repository to list workspaces in.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Flow.

