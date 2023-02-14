## fmeserver cancel

Cancel a running job on FME Server

### Synopsis

Cancels the job and marks it as aborted in the completed jobs section, but does not remove it from the database.

```
fmeserver cancel [flags]
```

### Examples

```

  # Cancel a job with id 42
  fmeserver cancel --id 42
	
```

### Options

```
  -h, --help        help for cancel
      --id string   	The ID of the job to cancel.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver](fmeserver.md)	 - A command line interface for interacting with FME Server.
