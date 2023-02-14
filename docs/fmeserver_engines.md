## fmeserver engines

Get information about the FME Engines

### Synopsis

Gets information and status about FME Engines currently connected to FME Server

```
fmeserver engines [flags]
```

### Examples

```

  # List all engines
  fmeserver engines
	
  # Output number of engines
  fmeserver engines --count
	
  # Output engines in json form
  fmeserver engines --json
	
  # Output just the names of the engines with no column headers
  fmeserver engines --output=custom-columns=NAME:.instanceName --no-headers
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
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver](fmeserver.md)	 - A command line interface for interacting with FME Server.
