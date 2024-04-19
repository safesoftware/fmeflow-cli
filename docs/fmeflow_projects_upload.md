## fmeflow projects upload

Imports FME Flow Projects from a downloaded package.

### Synopsis

Imports FME Flow Projects from a downloaded package. The upload happens in two steps. The package is uploaded to the server, a preview is generated that contains the list of items, and then the import is run. This command can be run using a few different modes.
- Using the --get-selectable flag will just generate the preview and output the selectable items in the package and then delete the import
- Using the --quick flag will skip the preview and import everything in the package by default.
- Using the --interactive flag will prompt the user to select items to import from the list of selectable items if any exist
- Using the --selected-items flag will import only the items specified. The default is to import all items in the package.

```
fmeflow projects upload [flags]
```

### Examples

```

  # Upload a project and import all selectable items if any exist
  fmeflow projects upload --file ProjectPackage.fsproject

  # Upload a project without overwriting existing items
  fmeflow projects upload --file ProjectPackage.fsproject --overwrite=false
  
  # Upload a project and perform a quick import
  fmeflow projects upload --file ProjectPackage.fsproject --quick
  
  # Upload a project and be prompted for which items to import of the selectable items
  fmeflow projects upload --file ProjectPackage.fsproject --interactive 
 
  # Upload a project and get the list of selectable items
  fmeflow projects upload --file ProjectPackage.fsproject --get-selectable
  
  # Upload a project and import only the specified selectable items
  fme projects upload --file ProjectPackage.fsproject --selected-items="mysqldb:connection,slack con:connector"
```

### Options

```
      --disable-project-items   Whether to disable items in the imported FME Server Projects. If true, items that are new or overwritten will be imported but disabled. If false, project items are imported as defined in the import package.
      --failure-topic string    Topic to notify on failure of the backup. (default "MIGRATION_ASYNC_JOB_FAILURE")
  -f, --file string             Path to backup file to upload to restore. Can be a local file or the relative path inside the specified shared resource.
      --get-selectable          Output the selectable items in the import package.
  -h, --help                    help for upload
      --interactive             Prompt interactively for the selectable items to import (if any exist).
      --overwrite               If specified, the items in the project will overwrite existing items. (default true)
      --pause-notifications     Disable notifications for the duration of the restore. (default true)
      --quick                   Import everything in the package by default.
      --selected-items string   The items to import. Set to "all" to import all items, and "none" to omit selectable items. Otherwise, this should be a comma separated list of item ids type pairs separated by a colon. e.g. a:b,c:d (default "all")
      --success-topic string    Topic to notify on success of the backup. (default "MIGRATION_ASYNC_JOB_SUCCESS")
      --wait                    Wait for import to complete. Set to false to return immediately after the import is started. (default true)
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow projects](fmeflow_projects.md)	 - List, Upload and Download projects on FME Flow

