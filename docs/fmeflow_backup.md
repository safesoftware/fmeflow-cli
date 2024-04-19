## fmeflow backup

Backs up the FME Server configuration

### Synopsis

Backs up the FME Server configuration to a local file or to a shared resource location on the FME Server.

```
fmeflow backup [flags]
```

### Examples

```

  # back up to a local file
  fmeflow backup -f my_local_backup.fsconfig
	
  # back up to the "Backup" folder in the FME Server Shared Resources with the file name my_fme_backup.fsconfig
  fmeflow backup --resource --export-package my_fme_backup.fsconfig
```

### Options

```
      --export-package string   Path and name of the export package. (default "ServerConfigPackage.fsconfig")
      --failure-topic string    Topic to notify on failure of the backup. Default is MIGRATION_ASYNC_JOB_FAILURE
  -f, --file string             Path to file to download the backup to. (default "ServerConfigPackage.fsconfig")
  -h, --help                    help for backup
      --resource                Backup to a shared resource instead of downloading.
      --resource-name string    Shared Resource Name where the exported package is saved. (default "FME_SHAREDRESOURCE_BACKUP")
      --success-topic string    Topic to notify on success of the backup. Default is MIGRATION_ASYNC_JOB_SUCCESS
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Flow.

