## Custom Columns

Many of the commands in the CLI support an output type of `custom-columns`. This allows you to output specific fields to the output table that exist in the raw JSON that the command returns. This can be helpful to get specialized information, or if a pipeline requires a certain field to pass into another command.

### General Syntax

The `custom-columns` output type will take a comma delimited list of header and value pairs that will make up the output table. The general syntax is:
```
HEADERNAME:.jsonPathForValue
```
`HEADERNAME` will be the heading of the column in the table, and `.jsonPathForValue` will be the JSONPath expression to retrieve the data for that column from the returned JSON. This allows you to dynamically set column values to whatever data you need from the json.

If a specified json query references a field that is missing on the json that it is querying, it will return an empty value instead of erroring.

### JSONPath

JSONPath is a query language for JSON. The implementation of it in the FME Server CLI is heavily based on the implementation of it for the `kubectl` CLI for Kubernetes. More detailed documentation on how to use it can be found in the [Kubernetes documentation](https://kubernetes.io/docs/reference/kubectl/jsonpath/) for things that are not covered here.

### How to use it

When crafting a `custom-column` output for a given command that supports it in the CLI, a good place to start is outputing the raw JSON for that command to see what fields are available.

For example, let's output all the workspaces in the `Samples` repository.

```
> fmeflow workspaces --repository Samples
 NAME                      TITLE                                                    LAST SAVE DATE                
 austinApartments.fmw      City of Austin: Apartments and other (SPATIALITE 2 KML)  2022-06-15 10:59:24 +0000 UTC 
 austinDownload.fmw        City of Austin: Data Download                            2022-06-15 13:40:29 +0000 UTC 
 earthquakesextrusion.fmw  Earthquakes: GeoJSON to KML Diagrams                     2022-06-15 15:52:07 +0000 UTC 
 easyTranslator.fmw        Generic format translator                                2022-06-15 15:56:11 +0000 UTC
```
In order to see what data we can output in our `custom-columns` output, we should take a look at the raw JSON:
```
> fmeflow workspaces --repository Samples --json
{
  "items": [
    {
      "name": "austinApartments.fmw",
      "title": "City of Austin: Apartments and other (SPATIALITE 2 KML)",
      "description": "",
      "type": "workspace",
      "repositoryName": "Samples",
      "lastSaveDate": "2022-06-15T10:59:24.000Z",
      "lastPublishDate": "2023-02-03T01:03:26.513Z",
      "lastPublishUserId": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
      "lastPublishUser": "admin",
      "totalFileSize": 6680576,
      "fileCount": 1,
      "averageCpuPercent": 0.0,
      "averagePeakMemoryUsage": 0,
      "averageCpuTime": 0,
      "totalRuns": 0,
      "averageElapsedTime": 0,
      "favorite": false
    },
   ...
  ],
  "totalCount": 4,
  "limit": 100,
  "offset": 0
}
```
The output has been truncated here to organize it better, but we can now see what is available for each workspace. In the FME Server CLI, the JSONPath query is applied to each item in the `items` array automatically. This means to access a field, we would do it from the item level. We also can omit the `$` from the start of the JSONPath query. For example, to access the field `averageCpuPercent`, the JSONPath query we need for the columns is simply `.averageCpuPercent`.

This means to output the workspace name along with the average CPU percent and average elapsed time, our CLI command looks like this:
```
> fmeflow workspaces --repository Samples --output custom-columns="NAME:.name,CPU PERCENT:.averageCpuPercent,ELAPSED TIME:.averageElapsedTime"
 NAME                      CPU PERCENT         ELAPSED TIME 
 austinApartments.fmw      13.242375601926163  1246         
 austinDownload.fmw        63.6535552193646    2644         
 earthquakesextrusion.fmw  37.996127783155856  2066         
 easyTranslator.fmw        61.24852767962309   849
```

For a slightly more complicated example, specifying a specific workspace gives more data:
```
> fmeflow workspaces --repository Samples --name "austinApartments.fmw" --json
{
  "name": "austinApartments.fmw",
  "title": "City of Austin: Apartments and other (SPATIALITE 2 KML)",
  "type": "workspace",
  "buildNumber": 22337,
  "category": "",
  "history": "",
  "lastSaveDate": "2022-06-15T10:59:24.000Z",
  "lastPublishDate": "2023-02-03T01:03:26.513Z",
  "lastSaveBuild": "FME(R) 2022.0.0.0 (20220428 - Build 22337 - WIN64)",
  "legalTermsConditions": "",
  "properties": [
    {
      "attributes": {},
      "category": "fmedatadownload_FMEUSERPROPDATA",
      "name": "ADVANCED",
      "value": ""
    },
    ...
  ],
  "requirements": "",
  "requirementsKeyword": "",
  "resources": [
    {
      "name": "Landmarks.sqlite",
      "size": 6680576
    }
  ],
  "usage": "",
  "userName": "admin",
  "fileSize": 136390,
  "favorite": false,
  "averageCpuPercent": 13.242375601926163,
  "averagePeakMemoryUsage": 5736760,
  "averageCpuTime": 165,
  "totalRuns": 1,
  "averageElapsedTime": 1246,
  "description": "",
  "datasets": {
    "source": [
      {
        "featureTypes": [
          {
            "attributes": [
              {
                "decimals": 0,
                "name": "CFCC",
                "type": "varchar",
                "width": 4
              },
              ...
            ],
            "description": "",
            "name": "Landmarks",
            "properties": []
          }
        ],
        "format": "SPATIALITE_NATIVE",
        "location": "$(FME_MF_DIR_USERTYPED)\\Landmarks.sqlite",
        "name": "SPATIALITE_NATIVE_1",
        "properties": [
          {
            "attributes": {},
            "category": "DATASET_DESCRIPTION",
            "name": "SPATIALITE_NATIVE_1",
            "value": "Landmarks.sqlite [SPATIALITE_NATIVE]"
          },
          ...
        ],
        "source": true
      }
    ],
    "destination": [
      {
        "featureTypes": [
          {
            "attributes": [
              {
                "decimals": 0,
                "name": "CFCC",
                "type": "kml_char",
                "width": 4
              },
              ...
            ],
            "description": "",
            "name": "Landmarks",
            "properties": []
          }
        ],
        "format": "OGCKML",
        "location": "/data/fmeflowdata/resources/system/temp/engineresults\\austinApartments.kml",
        "name": "OGCKML_1",
        "properties": [
          {
            "attributes": {},
            "category": "DATASET_DESCRIPTION",
            "name": "OGCKML_1",
            "value": "austinApartments.kml [OGCKML]"
          },
          ...
        ],
        "source": false
      }
    ]
  },
  "services": {
    ...
    }
  },
  "parameters": []
}
```

The JSON about a single workspace includes information on source and destination datasets, that are stored in a JSON list. In JSONPath, lists are accessed using `[]`, with a specific index being specified as a number such as `[0]`, or all results from that list being denoted as `[*]`. For example, if we want to get all the source formats for this workspace, we would use the JSonPath query `.datasets.source[*].format`. A full example:

```
> fmeflow workspaces --repository Samples --name "austinApartments.fmw" --output custom-columns="NAME:.name,SOURCE:.datasets.source[*].format,DEST:datasets.destination[*].format"
 NAME                  SOURCE             DEST   
 austinApartments.fmw  SPATIALITE_NATIVE  OGCKML
```

### SEE ALSO

* [fmeflow](fmeflow.md)	 - A command line interface for interacting with FME Server.

