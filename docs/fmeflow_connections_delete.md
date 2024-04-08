## fmeflow connections delete

Delete a connection

### Synopsis

Delete a connection.

```
fmeflow connections delete [flags]
```

### Examples

```

	Examples:
	# Delete a connection with the name "myConnection"
	fmeflow connections delete --name myConnection

```

### Options

```
  -h, --help          help for delete
      --name string   Name of the connection to delete.
  -y, --no-prompt     Description of the new repository.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow connections](fmeflow_connections.md)	 - Lists connections on FME Flow

