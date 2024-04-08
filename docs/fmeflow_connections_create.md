## fmeflow connections create

Create a connection

### Synopsis

Create a connection.

```
fmeflow connections create [flags]
```

### Examples

```

	Examples:
	# Create a connection with the name "myConnection" and the category "PostgreSQL" and the type "database" with username "myUser" and password "myPassword"
	fmeflow connections create --name myConnection --category database  --type PostgreSQL --username myUser --password myPassword

```

### Options

```
      --authenticationMethod string   Authentication method of the connection to create.
      --category string               Category of the connection to create. Typically it is one of: "basic", "database", "token", "oauthV1", "oauthV2".
  -h, --help                          help for create
      --name string                   Name of the connection to create.
      --parameter stringArray         Parameters of the connection to create. Must be of the form name=value. Can be specified multiple times.
      --password string               Password of the connection to create.
      --type string                   Type of connection.
      --username string               Username of the connection to create.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow connections](fmeflow_connections.md)	 - Lists connections on FME Flow

