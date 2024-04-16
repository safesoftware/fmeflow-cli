## fmeflow projects items

Lists the items for the specified project

### Synopsis

Lists the items contained in the specified project.

```
fmeflow projects items [flags]
```

### Examples

```

  # Get all items for a project via id
  fmeflow projects items --id a64297e7-a119-4e10-ac37-5d0bba12194b

  # Get all items for a project via name
  fmeflow projects items --name test_project

  # Get items with type workspace for a project via name
  fmeflow projects items --name test_project --type workspace
  
  # Get all items for a project via name without dependencies
  fmeflow projects items --name test_project --include-dependencies=false
  
  # Get all items for a project via name with a filter on name
  fmeflow projects items --name test_project --filter-string "test_name" --filter-properties "name"
```

### Options

```
      --filter-property stringArray   Property to filter by. Should be one of "name" or "owner". Can only be set if filter-string is also set
      --filter-string string          String to filter items by
  -h, --help                          help for items
      --id string                     Id of project to get items for 
      --include-dependencies          Include dependencies in the output (default true)
      --name string                   Name of project to get items for
      --no-headers                    Don't print column headers
  -o, --output string                 Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --type stringArray              Type of items to get
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow projects](fmeflow_projects.md)	 - List, Upload and Download projects on the FME Server

