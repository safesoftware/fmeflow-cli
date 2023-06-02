## fmeflow license request status

Check the status of a license request.

### Synopsis

Check the status of a license request.

```
fmeflow license request status [flags]
```

### Examples

```

	# Output the license request status as a table
	fmeflow license request status
	
	# Output the license Request status in json
	fmeflow license request status --json
	
	# Output just the status message
	fmeflow license request status --output custom-columns=STATUS:.status --no-headers
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

* [fmeflow license request](fmeflow_license_request.md)	 - Request a license from the FME Server licensing server

