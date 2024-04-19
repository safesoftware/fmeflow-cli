## fmeflow connections update

Update a connection

### Synopsis

Update a connection. Only things that need to be modified need to be specified.

```
fmeflow connections update [flags]
```

### Examples

```

  # Update a PostgreSQL connection with the name "myPGSQLConnection" and modify the host to "myDBHost"
  fmeflow connections update --name myPGSQLConnection --parameter HOST=myDBHost

```

### Options

```
      --authentication-method string   Authentication method of the connection to update.
  -h, --help                           help for update
      --name string                    Name of the connection to update.
      --parameter stringArray          Parameters of the connection to update. Must be of the form name=value. Can be specified multiple times.
      --password string                Password of the connection to update.
      --username string                Username of the connection to update.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow connections](fmeflow_connections.md)	 - Lists connections on FME Flow

