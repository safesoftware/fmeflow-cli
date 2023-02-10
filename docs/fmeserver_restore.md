## fmeserver restore

Restores the FME Server configuration from an import package

### Synopsis

Restores the FME Server configuration from an import package

```
fmeserver restore [flags]
```

### Examples

```

  # Restore from a backup in a local file
  fmeserver restore --file ServerConfigPackage.fsconfig

  # Restore from a backup in a local file using UPDATE mode
  fmeserver restore --file ServerConfigPackage.fsconfig --import-mode UPDATE
  
  # Restore from a backup file stored in the Backup resource folder (FME_SHAREDRESOURCE_BACKUP) named ServerConfigPackage.fsconfig
  fmeserver restore --resource --file ServerConfigPackage.fsconfig
  
  # Restore from a backup file stored in the Data resource folder (FME_SHAREDRESOURCE_DATA) named ServerConfigPackage.fsconfig and set a failure and success topic to notify
  fmeserver restore --resource --resource-name FME_SHAREDRESOURCE_DATA --file ServerConfigPackage.fsconfig --failure-topic MY_FAILURE_TOPIC --success-topic MY_SUCCESS_TOPIC
  
```

### Options

```
      --failure-topic string          Topic to notify on failure of the import. Default is MIGRATION_ASYNC_JOB_FAILURE. Only supported when restoring from a shared resource.
  -f, --file string                   Path to backup file to upload to restore. Can be a local file or the relative path inside the specified shared resource.
  -h, --help                          help for restore
      --import-mode string            To import only items in the import package that do not exist on the current instance, specify INSERT. To overwrite items on the current instance with those in the import package, specify UPDATE. Default is INSERT. (default "INSERT")
      --pause-notifications           Disable notifications for the duration of the restore. (default true)
      --projects-import-mode string   Import mode for projects. To import only projects in the import package that do not exist on the current instance, specify INSERT. To overwrite projects on the current instance with those in the import package, specify UPDATE. If not supplied, importMode will be used.
      --resource                      Restore from a shared resource location instead of a local file.
      --resource-name string          Resource containing the import package. (default "FME_SHAREDRESOURCE_BACKUP")
      --success-topic string          Topic to notify on success of the import. Default is MIGRATION_ASYNC_JOB_SUCCESS. Only supported when restoring from a shared resource.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver](fmeserver.md)	 - A command line interface for interacting with FME Server.

