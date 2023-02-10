## fmeserver jobs

Lists jobs on FME Server

### Synopsis

Lists jobs on FME Server

```
fmeserver jobs [flags]
```

### Examples

```

  # List all jobs (currently limited to the most recent 1000)
  fmeserver jobs --all
	
  # List all running jobs
  fmeserver jobs --running
	
  # List all jobs from a given repository
  fmeserver jobs --repository Samples
	
  # List all jobs that ran a given workspace
  fmeserver jobs --repository Samples --workspace austinApartments.fmw
	
  # List all jobs in JSON format
  fmeserver jobs --json
	
  # List the workspace, CPU time and peak memory usage for a given repository
  fmeserver jobs --repository Samples --output="custom-columns=WORKSPACE:.workspace,CPU Time:.cpuTime,Peak Memory:.peakMemUsage"
	
```

### Options

```
      --active               Retrieve active jobs
      --all                  Retrieve all jobs
      --completed            Retrieve completed jobs
  -h, --help                 help for jobs
      --id int               Specify the job id to display (default -1)
      --no-headers           Don't print column headers
  -o, --output string        Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --queued               Retrieve queued jobs
      --repository string    If specified, only jobs from the specified repository will be returned.
      --running              Retrieve running jobs
      --source-id string     If specified along with source type, only jobs from the specified type with the specified id will be returned. For Automations, the source id is the automation id. For WorkspaceSubscriber, the source id is the id of the subscription. For Scheduler, the source id is the category and name of the schedule separated by '/'. For example, 'Category/Name'.
      --source-type string   If specified, only jobs run by this source type will be returned.
      --user-name string     If specified, only jobs run by the specified user will be returned.
      --workspace string     If specified along with repository, only jobs from the specified repository and workspace will be returned.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeserver-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeserver](fmeserver.md)	 - A command line interface for interacting with FME Server.

