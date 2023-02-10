## fmeserver license request

Request a license from the FME Server licensing server

### Synopsis

Request a license file from the FME Server licensing server. First name, Last name and email are required for requesting a license file.
  If no serial number is passed in, a trial license will be requested.

```
fmeserver license request [flags]
```

### Examples

```

  # Request a trial license and wait for it to be downloaded and installed
  fmeserver license request --first-name "Billy" --last-name "Bob" --email "billy.bob@example.com" --company "Example Company Inc." --wait
	
  # Request a license with a serial number
  fmeserver license request --first-name "Billy" --last-name "Bob" --email "billy.bob@example.com" --company "Example Company Inc." --serial-number "AAAA-BBBB-CCCC"
	
```

### Options

```
      --category string        License Category
      --company string         Company for the licensing request
      --email string           Email address for license request.
      --first-name string      First name to use for license request.
  -h, --help                   help for request
      --industry string        Industry for the licensing request
      --last-name string       Last name to use for license request.
      --sales-source string    Sales source
      --serial-number string   Serial Number for the license request.
      --subscribe-to-updates   Subscribe to Updates
      --wait                   Wait for licensing request to finish
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver license](fmeserver_license.md)	 - Interact with licensing an FME Server
* [fmeserver license request status](fmeserver_license_request_status.md)	 - Check the status of a license request.

