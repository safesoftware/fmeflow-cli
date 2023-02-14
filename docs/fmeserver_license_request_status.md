## fmeserver license request status

Check the status of a license request.

### Synopsis

Check the status of a license request.

```
fmeserver license request status [flags]
```

### Examples

```

	# Output the license request status as a table
	fmeserver license request status
	
	# Output the license Request status in json
	fmeserver license request status --json
	
	# Output just the status message
	fmeserver license request status --output custom-columns=STATUS:.status --no-headers
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

* [fmeserver license request](fmeserver_license_request.md)	 - Request a license from the FME Server licensing server
