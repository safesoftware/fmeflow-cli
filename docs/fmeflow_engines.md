## fmeflow engines

Get information about the FME Engines

### Synopsis

Gets information and status about FME Engines currently connected to FME Server

```
fmeflow engines [flags]
```

### Examples

```

  # List all engines
  fmeflow engines
	
  # Output number of engines
  fmeflow engines --count
	
  # Output engines in json form
  fmeflow engines --json
	
  # Output just the names of the engines with no column headers (V4)
  fmeflow engines --output=custom-columns=NAME:.name --no-headers
```

### Options

```
      --count           Prints the total count of engines.
  -h, --help            help for engines
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

