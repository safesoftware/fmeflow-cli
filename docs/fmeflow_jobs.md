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
  fmeflow jobs --repository Samples --output="custom-columns=WORKSPACE:.workspace,CPU Time:.cpuTime,Peak Memory:.peakMemoryUsage"
	
```

### Options

```
      --active               Retrieve active jobs, in V4 it is equivalent to [queued; running]
      --all                  Retrieve all jobs
      --cancelled            V4 Retrieve all cancelled jobs
      --completed            Retrieve completed jobs, in V4 it is equivalent to [success; failure; cancelled]
      --failure              V4 Retrieve failed jobs 
  -h, --help                 help for jobs
      --id int               Specify the job id to display (default -1)
      --no-headers           Don't print column headers
  -o, --output string        Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --queued               V4 Retrieve queued jobs
      --repository string    If specified, only jobs from the specified repository will be returned.
      --running              V4 Retrieve running jobs
      --source-id string     If specified along with source type, only jobs from the specified type with the specified id will be returned. For Automations, the source id is the automation id. For WorkspaceSubscriber, the source id is the id of the subscription. For Scheduler, the source id is the category and name of the schedule separated by '/'. For example, 'Category/Name'.
      --source-type string   If specified, only jobs run by this source type will be returned.
      --success              V4 Retrieve successful jobs
      --user-name string     If specified, only jobs run by the specified user will be returned.
      --workspace string     If specified along with repository, only jobs from the specified repository and workspace will be returned.
      --status               V4 Jobs matching any of the specified statuses will be returned. [queued, running] are mutually exclusive with [success, failure, cancelled]
      --engine-name          V4 If specified, only jobs run by the specified engine will be returned. Queued jobs cannot be filtered by engine.
      --queue                V4 If specified, only jobs routed through the specified queue will be returned.
      --sort                 V4 Specifies the sort order of the result. Append _asc to the property name to sort in ascending order or _desc to sort in descending order. For example, workspace_asc to sort by workspace in ascending order. If the suffix is unspecified, the property will be sorted in ascending order. If no sorting parameter is specified, the result will be sorted by timeFinished in descending order by default. Sorting is only supported on workspace, timeFinished, timeStarted, and status.

```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Flow.

