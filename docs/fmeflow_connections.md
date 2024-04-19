## fmeflow connections

Lists connections on FME Flow

### Synopsis

Lists connections on FME Flow. Pass in a name to retrieve information on a single connection.

```
fmeflow connections [flags]
```

### Examples

```

  # List all connections
  fmeflow connections

  # Get a single connection with the name "myConnection"
  fmeflow connections --name myConnection
  
  # List all connections of type "Google Drive"
  fmeflow connections --type "Google Drive"
  
  # List all database connections
  fmeflow connections --category database
  
  # List the PostgreSQL connections with custom columns showing the name and host of the database connections
  fmeflow connections --category "database" --type "PostgreSQL" --output=custom-columns="NAME:.name,HOST:.parameters.HOST" 
```

### Options

```
      --category stringArray        The categories of connections to return. Can be passed in multiple times
      --excluded-type stringArray   The types of connections to exclude. Can be passed in multiple times
  -h, --help                        help for connections
      --name string                 Return a single project with the given name.
      --no-headers                  Don't print column headers
  -o, --output string               Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --type stringArray            The types of connections to return. Can be passed in multiple times
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Flow.
* [fmeflow connections create](fmeflow_connections_create.md)	 - Create a connection
* [fmeflow connections delete](fmeflow_connections_delete.md)	 - Delete a connection
* [fmeflow connections update](fmeflow_connections_update.md)	 - Update a connection

