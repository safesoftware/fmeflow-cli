## fmeserver healthcheck

Retrieves the health status of FME Server

### Synopsis

Retrieves the health status of FME Server. The health status is normal if the FME Server REST API is responsive. Note that this endpoint does not require authentication. Load balancer or other systems can monitor FME Server using this endpoint without supplying token or password credentials.

```
fmeserver healthcheck [flags]
```

### Examples

```

  # Check if the FME Server is healthy and accepting requests
  fmeserver healthcheck
		
  # Check if the FME Server is healthy and ready to process jobs
  fmeserver healthcheck --ready
		
  # Check if the FME Server is healthy and output in json
  fmeserver healthcheck --json
  
  # Check that the FME Server is healthy and output just the status
  fmeserver healthcheck --output=custom-columns=STATUS:.status
```

### Options

```
  -h, --help            help for healthcheck
      --no-headers      Don't print column headers
  -o, --output string   Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --ready           The health check will report the status of FME Server if it is ready to process jobs.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver](fmeserver.md)	 - A command line interface for interacting with FME Server.
