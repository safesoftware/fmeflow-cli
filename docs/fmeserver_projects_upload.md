## fmeflow projects upload

Imports FME Server Projects from a downloaded package.

### Synopsis

Imports FME Server Projects from a downloaded package. Useful for moving a project from one FME Server to another.

```
fmeflow projects upload [flags]
```

### Examples

```

  # Restore from a backup in a local file
  fmeflow projects upload --file ProjectPackage.fsproject

  # Restore from a backup in a local file using UPDATE mode
  fmeflow projects upload --file ProjectPackage.fsproject --import-mode UPDATE
```

### Options

```
      --disable-project-items         Whether to disable items in the imported FME Server Projects. If true, items that are new or overwritten will be imported but disabled. If false, project items are imported as defined in the import package.
  -f, --file string                   Path to backup file to upload to restore. Can be a local file or the relative path inside the specified shared resource.
  -h, --help                          help for upload
      --import-mode string            To import only items in the import package that do not exist on the current instance, specify INSERT. To overwrite items on the current instance with those in the import package, specify UPDATE. Default is INSERT. (default "INSERT")
      --pause-notifications           Disable notifications for the duration of the restore. (default true)
      --projects-import-mode string   Import mode for projects. To import only projects in the import package that do not exist on the current instance, specify INSERT. To overwrite projects on the current instance with those in the import package, specify UPDATE. If not supplied, importMode will be used.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow projects](fmeflow_projects.md)	 - Lists projects on the FME Server

