## fmeflow migration tasks

Retrieves the records for migration tasks.

### Synopsis

Retrieves the records for migration tasks. Get all migration tasks or for a specific task by passing in the id.

```
fmeflow migration tasks [flags]
```

### Examples

```

  # Get all migration tasks
  fmeflow migration tasks
	
  # Get all migration tasks in json
  fmeflow migration tasks --json
	
  # Get the migration task for a given id
  fmeflow migration tasks --id 1
	
  # Output the migration log for a given id to the console
  fmeflow migration tasks --id 1 --log
	
  # Output the migration log for a given id to a local file
  fmeflow migration tasks --id 1 --log --file my-backup-log.txt
	
  # Output just the start and end time of the a given id
  fmeflow migration tasks --id 1 --output="custom-columns=Start Time:.startDate,End Time:.finishedDate"
```

### Options

```
      --file string     File to save the log to.
  -h, --help            help for tasks
      --id int          Retrieves the record for a migration task according to the given ID. (default -1)
      --log             Downloads the log file of a migration task.
      --no-headers      Don't print column headers
  -o, --output string   Specify the output type. Should be one of table, json, or custom-columns (default "table")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow migration](fmeflow_migration.md)	 - Returns information on migrations using the tasks subcommand.

