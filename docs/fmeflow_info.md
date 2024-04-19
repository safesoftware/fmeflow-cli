## fmeflow info

Retrieves build, version and time information about FME Server

### Synopsis

Retrieves build, version and time information about FME Server

```
fmeflow info [flags]
```

### Examples

```

  # Output FME Server information in a table
  fmeflow info

  # Output FME Server information in json
  fmeflow info --json

  # Output just the build string with no column headers
  fmeflow info --output=custom-columns="BUILD:.build" --no-headers
	
```

### Options

```
  -h, --help            help for info
      --no-headers      Don't print column headers
  -o, --output string   Specify the output type. Should be one of table, json, or custom-columns (default "table")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Flow.

