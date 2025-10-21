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

  # Mix and match multiple statuses in V4
  fmeflow jobs --active
  fmeflow jobs --running --success
  fmeflow jobs --queued --failure
  fmeflow jobs --running --success --cancelled
  etc.
	
  # List all jobs from a given repository
  fmeflow jobs --repository Samples
	
  # List all jobs that ran a given workspace
  fmeflow jobs --repository Samples --workspace austinApartments.fmw
	
  # List all jobs in JSON format
  fmeflow jobs --json
	
  # List the workspace, CPU time and peak memory usage for a given repository
  fmeflow jobs --repository Samples --output="custom-columns=WORKSPACE:.workspace,CPU Time:.cpuTime"
	
```

### Options

```
      --active               Retrieve active jobs
      --all                  Retrieve all jobs
      --api-version string   The api version to use when contacting FME Server. Must be one of v3 or v4
      --cancelled            V4 only - Retrieve cancelled jobs
      --completed            Retrieve completed jobs
      --engine-name string   V4 only - If specified, only jobs run by the specified engine will be returned. Queued jobs cannot be filtered by engine. For v4 API only
      --failure              V4 only - Retrieve failed jobs
  -h, --help                 help for jobs
      --id int               Specify the job id to display (default -1)
      --no-headers           Don't print column headers
  -o, --output string        Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --queue string         V4 only - If specified, only jobs routed through the specified queue will be returned. For v4 API only
      --queued               Retrieve queued jobs
      --repository string    If specified, only jobs from the specified repository will be returned.
      --running              Retrieve running jobs
      --sort string          V4 only - Sort jobs by one of: workspace, timeFinished, timeStarted, status. Append _asc or _desc to specify ascending or descending order. For example: workspace_asc. For v4 API only
      --source-id string     If specified along with source type, only jobs from the specified type with the specified id will be returned. For Automations, the source id is the automation id. For WorkspaceSubscriber, the source id is the id of the subscription. For Scheduler, the source id is the category and name of the schedule separated by '/'. For example, 'Category/Name'.
      --source-type string   If specified, only jobs run by this source type will be returned.
      --success              V4 only - Retrieve succeeded jobs
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

