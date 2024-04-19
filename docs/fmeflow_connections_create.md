## fmeflow connections create

Create a connection

### Synopsis

Create a connection.

```
fmeflow connections create [flags]
```

### Examples

```

  # Create a PostgreSQL connection
  fmeflow connections create --name myPGSQLConnection --category database --type PostgreSQL --parameter HOST=myDBHost --parameter PORT=5432 --parameter DATASET=dbname --parameter USER_NAME=dbuser --parameter SSL_OPTIONS="" --parameter SSLMODE=prefer

  # Create a Google Drive connection (web service must already exist on FME Flow)
  fmeflow connections create --name googleDriveConn --category oauthV2 --type "Google Drive"

```

### Options

```
      --authentication-method string   Authentication method of the connection to create.
      --category string                Category of the connection to create. Typically it is one of: "basic", "database", "token", "oauthV1", "oauthV2".
  -h, --help                           help for create
      --name string                    Name of the connection to create.
      --parameter stringArray          Parameters of the connection to create. Must be of the form name=value. Can be specified multiple times.
      --password string                Password of the connection to create.
      --type string                    Type of connection.
      --username string                Username of the connection to create.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow connections](fmeflow_connections.md)	 - Lists connections on FME Flow

