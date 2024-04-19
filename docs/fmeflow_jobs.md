## fmeflow jobs

Lists jobs on FME Server

### Synopsis

Lists running, queued, and/or queued jobs on FME Server. Pass in a job id to get information on a specific job.

```
fmeflow jobs [flags]
```

### Examples

```

  # List all jobs (currently limited to the most recent 1000)
  fmeflow jobs --all
	
  # List all running jobs
  fmeflow jobs --running
	
  # List all jobs from a given repository
  fmeflow jobs --repository Samples
	
  # List all jobs that ran a given workspace
  fmeflow jobs --repository Samples --workspace austinApartments.fmw
	
  # List all jobs in JSON format
  fmeflow jobs --json
	
  # List the workspace, CPU time and peak memory usage for a given repository
  fmeflow jobs --repository Samples --output="custom-columns=WORKSPACE:.workspace,CPU Time:.cpuTime,Peak Memory:.peakMemUsage"
	
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
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Flow.

