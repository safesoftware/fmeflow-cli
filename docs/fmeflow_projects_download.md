## fmeflow projects download

Downloads an FME Server Project

### Synopsis

Downloads an FME Server Project to a local file. Useful for backing up or moving a project to another FME Server.

```
fmeflow projects download [flags]
```

### Examples

```

  # download a project named "Test Project" to a local file with default name
  fmeflow projects download --name "Test Project"
	
  # download a project named "Test Project" to a local file named MyProject.fsproject
  fmeflow projects download --name "Test Project" -f MyProject.fsproject
```

### Options

```
      --exclude-sensitive-info   Whether to exclude sensitive information from the exported package. Sensitive information will be excluded from connections, subscriptions, publications, schedule tasks, S3 resources, and user accounts. Other items in the project may still contain sensitive data, especially workspaces. Please be careful before sharing the project export pacakge with others.
  -f, --file string              Path to file to download the backup to. (default "ProjectPackage.fsproject")
  -h, --help                     help for download
      --name string              Name of the project to download.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow projects](fmeflow_projects.md)	 - List, Upload and Download projects on the FME Server

