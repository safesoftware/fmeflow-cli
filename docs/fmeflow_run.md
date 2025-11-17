## fmeflow run

Run a workspace on FME Server.

### Synopsis

Run a workspace on FME Server.

```
fmeflow run [flags]
```

### Examples

```

  # Submit a job asynchronously
  fmeflow run --repository Samples --workspace austinApartments.fmw
	
  # Submit a job and wait for it to complete
  fmeflow run --repository Samples --workspace austinApartments.fmw --wait
	
  # Submit a job to a specific queue and set a time to live in the queue
  fmeflow run --repository Samples --workspace austinApartments.fmw --queue Queue1 --max-time-in-queue 120
	
  # Submit a job and pass in a few published parameters
  fmeflow run --repository Samples --workspace austinDownload.fmw --published-parameter-list THEMES=railroad,airports --published-parameter COORDSYS=TX83-CF
	
  # Submit a job, wait for it to complete, and customize the output
  fmeflow run --repository Samples --workspace austinApartments.fmw --wait --output="custom-columns=Time Requested:.timeRequested,Time Started:.timeStarted,Time Finished:.timeFinished"
	
  # Upload a local file to use as the source data for the translation
  fmeflow run --repository Samples --workspace austinApartments.fmw --file Landmarks-edited.sqlite --wait
```

### Options

```
      --repository string                      The name of the repository containing the workspace to run.
      --workspace string                       The name of the workspace to run.
      --wait                                   Submit job and wait for it to finish.
      --published-parameter stringArray        Published parameters defined for this workspace. Specify as Key=Value. Can be passed in multiple times. For list parameters, use the --list-published-parameter flag.
      --published-parameter-list stringArray   A List-type published parameters defined for this workspace. Specify as Key=Value1,Value2. Can be passed in multiple times.
      --file string                            Upload a local file Source dataset to use to run the workspace. Note this causes the translation to run in synchonous mode whether the --wait flag is passed in or not. For v3 API only.
      --run-until-canceled                     Runs a job until it is explicitly canceled. The job will run again regardless of whether the job completed successfully, failed, or the server crashed or was shut down. For v3 API only.
      --description string                     Description of the request. For v3 API only.
      --success-topic stringArray              Topics to notify when the job succeeds. Can be specified more than once.
      --failure-topic stringArray              Topics to notify when the job fails. Can be specified more than once.
      --node-manager-directive stringArray     Additional NM Directives, which can include client-configured keys, to pass to the notification service for custom use by subscriptions. Specify as Key=Value Can be passed in multiple times. For v3 API only.
      --directive stringArray                  Additional directives to pass to the job submission. Specify as Key=Value. Can be passed in multiple times. For v4 API only.
      --queue string                           Queue of the job to submit. Equavalent to --tag (deprecated).
      --max-job-runtime int                    Time (in seconds) elapsed for a running job before it's cancelled. The minimum value is 1 second, values less than 1 second are ignored. Equavalent to --time-until-canceled (deprecated). (default -1)
      --max-time-in-queue int                  Time to live in the job queue (in seconds). Equavalent to --time-to-live (deprecated). (default -1)
      --max-total-life-time int                Time to live including both time in the queue and run time (in seconds). The maximum value is 86400 and the minimum value is 1. For v4 API only. (default -1)
  -o, --output string                          Specify the output type. Should be one of table, json, or custom-columns (default "table")
      --no-headers                             Don't print column headers
  -h, --help                                   help for run
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/.fmeflow-cli.yaml)
      --json            Output JSON
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Flow.

