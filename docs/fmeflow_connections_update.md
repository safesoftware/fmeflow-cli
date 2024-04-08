## fmeflow connections update

Update a connection

### Synopsis

Update a connection.

```
fmeflow connections update [flags]
```

### Examples

```

	Examples:
	# Update a connection with the name "myConnection" and the category "PostgreSQL" and the type "database" with username "myUser" and password "myPassword"
	fmeflow connections update --name myConnection --category database --username myUser --password myPassword

```

### Options

```
      --authenticationMethod string   Authentication method of the connection to update.
  -h, --help                          help for update
      --name string                   Name of the connection to update.
      --parameter stringArray         Parameters of the connection to update. Must be of the form name=value. Can be specified multiple times.
      --password string               Password of the connection to update.
      --username string               Username of the connection to update.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow connections](fmeflow_connections.md)	 - Lists connections on FME Flow

