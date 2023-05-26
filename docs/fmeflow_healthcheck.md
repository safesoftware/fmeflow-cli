## fmeflow healthcheck

Retrieves the health status of FME Server

### Synopsis

Retrieves the health status of FME Server. The health status is normal if the FME Server REST API is responsive. Note that this endpoint does not require authentication. This command can be used without calling the login command first. The FME Server url can be passed in using the --url flag without needing a config file. A config file without a token can also be used.

```
fmeflow healthcheck [flags]
```

### Examples

```

  # Check if the FME Server is healthy and accepting requests
  fmeflow healthcheck
		
  # Check if the FME Server is healthy and ready to process jobs
  fmeflow healthcheck --ready
		
  # Check if the FME Server is healthy and output in json
  fmeflow healthcheck --json
  
  # Check that the FME Server is healthy and output just the status
  fmeflow healthcheck --output=custom-columns=STATUS:.status
  
 # Check the FME Server is healthy without needing a config file
 fmeflow healthcheck --url https://my-fmeflow.internal
 
 # Check the FME Server is healthy with a manually created config file
 cat << EOF >fmeflow-cli.yaml
 build: 23235
 url: https://my-fmeflow.internal
 EOF
 fmeflow healthcheck --config fmeflow-cli.yaml
```

### Options

```
  -h, --help            help for healthcheck
      --no-headers      Don't print column headers
  -o, --output string   Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --ready           The health check will report the status of FME Server if it is ready to process jobs.
      --url string      The base URL of the FME Server to check the health of. Pass this in if checking the health of an FME Server that you haven't called the login command for.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Server.

