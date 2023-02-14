## fmeserver license refresh status

Check the status of a license refresh request.

### Synopsis

Check the status of a license refresh request.

```
fmeserver license refresh status [flags]
```

### Examples

```

	# Output the license refresh status as a table
	fmeserver license refresh status
	
	# Output the license refresh status in json
	fmeserver license refresh status --json
	
	# Output just the status message
	fmeserver license refresh status --output custom-columns=STATUS:.status --no-headers
```

### Options

```
  -h, --help            help for status
      --no-headers      Don't print column headers
  -o, --output string   Specify the output type. Should be one of table, json, or custom-columns (default "table")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver license refresh](fmeserver_license_refresh.md)	 - Refreshes the installed license file with a current license from Safe Software.
