## fmeserver info

Retrieves build, version and time information about FME Server

### Synopsis

Retrieves build, version and time information about FME Server

```
fmeserver info [flags]
```

### Examples

```

  # Output FME Server information in a table
  fmeserver info

  # Output FME Server information in json
  fmeserver info --json

  # Output just the build string with no column headers
  fmeserver info --output=custom-columns="BUILD:.build" --no-headers
	
```

### Options

```
  -h, --help            help for info
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
