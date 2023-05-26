## fmeflow license refresh status

Check the status of a license refresh request.

### Synopsis

Check the status of a license refresh request.

```
fmeflow license refresh status [flags]
```

### Examples

```

	# Output the license refresh status as a table
	fmeflow license refresh status
	
	# Output the license refresh status in json
	fmeflow license refresh status --json
	
	# Output just the status message
	fmeflow license refresh status --output custom-columns=STATUS:.status --no-headers
```

### Options

```
  -h, --help            help for status
      --no-headers      Don't print column headers
  -o, --output string   Specify the output type. Should be one of table, json, or custom-columns (default "table")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow license refresh](fmeflow_license_refresh.md)	 - Refreshes the installed license file with a current license from Safe Software.

