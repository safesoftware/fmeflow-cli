## fmeflow projects delete

Deletes an FME Flow Project

### Synopsis

Deletes an FME Flow Project from the FME Server. Can optionally also delete the project contents and its dependencies.

```
fmeflow projects delete [flags]
```

### Examples

```

  # Delete a project by id
  fmeflow projects delete --id 123

  # Delete a project by name
  fmeflow projects delete --name "My Project"
  
  # Delete a project by name and all its contents
  fmeflow projects delete --name "My Project" --all
  
  # Delete a project by name and all its contents and dependencies
  fmeflow projects delete --name "My Project" --all --dependencies
```

### Options

```
      --all            Delete the project and its contents
      --dependencies   Delete the project and its contents and dependencies. Can only be specified if all is also specified
  -h, --help           help for delete
      --id string      The id of the project to delete. Either id or name must be specified
      --name string    The name of the project to delete. Either id or name must be specified
  -y, --no-prompt      Do not prompt for confirmation.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow projects](fmeflow_projects.md)	 - List, Upload and Download projects on the FME Flow

