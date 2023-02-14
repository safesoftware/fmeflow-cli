## fmeserver license requestfile

Generates a JSON file for requesting a FME Server license file.

### Synopsis

Generates a JSON file for requesting a FME Server license file.
		
	Example:
	
	# Generate a license request file and output to the console
	fmeserver license requestfile --first-name "Billy" --last-name "Bob" --email "billy.bob@example.com" --company "Example Company Inc."
	
	# Generate a license request file and output to a local file
	fmeserver license requestfile --first-name "Billy" --last-name "Bob" --email "billy.bob@example.com" --company "Example Company Inc." --file my-request-file.json

```
fmeserver license requestfile [flags]
```

### Options

```
      --category string        License Category
      --company string         Company for the licensing request
      --email string           Email address for license request.
      --file string            Path to file to output to.
      --first-name string      First name to use for license request.
  -h, --help                   help for requestfile
      --industry string        Industry for the licensing request
      --last-name string       Last name to use for license request.
      --sales-source string    Sales source
      --serial-number string   Serial Number for the license request.
      --subscribe-to-updates   Subscribe to Updates
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver license](fmeserver_license.md)	 - Interact with licensing an FME Server

