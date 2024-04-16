## fmeflow projects

List, Upload and Download projects on the FME Server

### Synopsis

List projects on the FME Server. Pass in a name to retrieve information on a single project. Use the upload and download sub commands to upload and download projects.

```
fmeflow projects [flags]
```

### Examples

```

  # List all projects
  fmeflow projects

  # List all projects owned by the user admin
  fmeflow projects --owner admin
  
  # Get a single project by name
  fmeflow projects --name "My Project"
  
  # Get a single project by id
  fmeflow projects --id a64297e7-a119-4e10-ac37-5d0bba12194b
  
  # Get a single project by name and output as JSON
  fmeflow projects --name "My Project" --output json
  
  # Get all projects and output as custom columns
  fmeflow projects --output=custom-columns=ID:.id,NAME:.name
```

### Options

```
      --api-version string   The api version to use when contacting FME Server. Must be one of v3 or v4
  -h, --help                 help for projects
      --id string            Return a single project with the given id. (v4 only)
      --name string          Return a single project with the given name.
      --no-headers           Don't print column headers
  -o, --output string        Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --owner string         If specified, only projects owned by the specified user will be returned.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Flow.
* [fmeflow projects delete](fmeflow_projects_delete.md)	 - Deletes an FME Flow Project
* [fmeflow projects download](fmeflow_projects_download.md)	 - Downloads an FME Server Project
* [fmeflow projects items](fmeflow_projects_items.md)	 - Lists the items for the specified project
* [fmeflow projects upload](fmeflow_projects_upload.md)	 - Imports FME Flow Projects from a downloaded package.

