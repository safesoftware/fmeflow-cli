## fmeflow projects

Lists projects on the FME Server

### Synopsis

Lists projects on the FME Server. Pass in a name to retrieve information on a single project.

```
fmeflow projects [flags]
```

### Examples

```

  # List all projects
  fmeflow projects

  # List all projects owned by the user admin
  fmeflow projects --owner admin
```

### Options

```
  -h, --help            help for projects
      --name string     Return a single project with the given name.
      --no-headers      Don't print column headers
  -o, --output string   Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --owner string    If specified, only projects owned by the specified user will be returned.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Server.
* [fmeflow projects download](fmeflow_projects_download.md)	 - Downloads an FME Server Project
* [fmeflow projects upload](fmeflow_projects_upload.md)	 - Imports FME Flow Projects from a downloaded package.

