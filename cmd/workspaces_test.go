package cmd

import (
	"net/http"
	"testing"
)

func TestWorkspaces(t *testing.T) {
	responseV3 := `{
		"offset": -1,
		"limit": -1,
		"totalCount": 4,
		"items": [
		  {
			"lastSaveDate": "2022-06-15T10:59:24Z",
			"avgCpuPct": 92.85714285714286,
			"avgPeakMemUsage": 3123940,
			"description": "",
			"repositoryName": "Samples",
			"title": "City of Austin: Apartments and other (SPATIALITE 2 KML)",
			"type": "WORKSPACE",
			"userName": "admin",
			"fileCount": 1,
			"avgCpuTime": 728,
			"lastPublishDate": "2022-12-07T18:50:52Z",
			"name": "austinApartments.fmw",
			"totalFileSize": 6680576,
			"totalRuns": 2,
			"avgElapsedTime": 784
		  },
		  {
			"lastSaveDate": "2022-06-15T13:40:29Z",
			"avgCpuPct": 71.16968698517299,
			"avgPeakMemUsage": 45244856,
			"description": "",
			"repositoryName": "Samples",
			"title": "City of Austin: Data Download",
			"type": "WORKSPACE",
			"userName": "admin",
			"fileCount": 30,
			"avgCpuTime": 1728,
			"lastPublishDate": "2022-12-07T18:50:52Z",
			"name": "austinDownload.fmw",
			"totalFileSize": 23596057,
			"totalRuns": 1,
			"avgElapsedTime": 2428
		  },
		  {
			"lastSaveDate": "2022-06-15T15:52:07Z",
			"avgCpuPct": 63.57292393579902,
			"avgPeakMemUsage": 3828352,
			"description": "",
			"repositoryName": "Samples",
			"title": "Earthquakes: GeoJSON to KML Diagrams",
			"type": "WORKSPACE",
			"userName": "admin",
			"fileCount": 0,
			"avgCpuTime": 911,
			"lastPublishDate": "2022-12-07T18:50:53Z",
			"name": "earthquakesextrusion.fmw",
			"totalFileSize": 0,
			"totalRuns": 1,
			"avgElapsedTime": 1433
		  },
		  {
			"lastSaveDate": "2022-06-15T15:56:11Z",
			"avgCpuPct": 0,
			"avgPeakMemUsage": 0,
			"description": "",
			"repositoryName": "Samples",
			"title": "Generic format translator",
			"type": "WORKSPACE",
			"userName": "admin",
			"fileCount": 0,
			"avgCpuTime": 0,
			"lastPublishDate": "2022-12-07T18:50:53Z",
			"name": "easyTranslator.fmw",
			"totalFileSize": 0,
			"totalRuns": 0,
			"avgElapsedTime": 0
		  }
		]
	  }`

	responseV3SingleWorkspace := `{
		"legalTermsConditions": "",
		"avgCpuPct": 92.85714285714286,
		"usage": "",
		"avgPeakMemUsage": 3123940,
		"description": "",
		"datasets": {
		  "destination": [
			{
			  "format": "OGCKML",
			  "name": "OGCKML_1",
			  "location": "/data/fmeserverdata/resources/system/temp/engineresults\\austinApartments.kml",
			  "source": false,
			  "featuretypes": [
				{
				  "name": "Landmarks",
				  "description": "",
				  "attributes": [
					{
					  "decimals": 0,
					  "name": "CFCC",
					  "width": 4,
					  "type": "kml_char"
					},
					{
					  "decimals": 0,
					  "name": "FILE",
					  "width": 20,
					  "type": "kml_char"
					},
					{
					  "decimals": 0,
					  "name": "LALAT",
					  "width": 20,
					  "type": "kml_char"
					},
					{
					  "decimals": 0,
					  "name": "LALONG",
					  "width": 20,
					  "type": "kml_char"
					},
					{
					  "decimals": 0,
					  "name": "LANAME",
					  "width": 31,
					  "type": "kml_char"
					},
					{
					  "decimals": 0,
					  "name": "LAND",
					  "width": 20,
					  "type": "kml_char"
					},
					{
					  "decimals": 0,
					  "name": "MODULE",
					  "width": 9,
					  "type": "kml_char"
					},
					{
					  "decimals": 0,
					  "name": "OBJECTID",
					  "width": 20,
					  "type": "kml_char"
					},
					{
					  "decimals": 0,
					  "name": "SOURCE",
					  "width": 2,
					  "type": "kml_char"
					}
				  ],
				  "properties": []
				}
			  ],
			  "properties": [
				{
				  "name": "OGCKML_1",
				  "attributes": {},
				  "category": "DEFAULT_DATASET_PATH",
				  "value": "OGCKML_1/austinApartments.kml"
				},
				{
				  "name": "OGCKML_1",
				  "attributes": {},
				  "category": "DATASET_DESCRIPTION",
				  "value": "austinApartments.kml [OGCKML]"
				},
				{
				  "name": "MIME_TYPE",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": ".kml application/vnd.google-earth.kml+xml .kmz application/vnd.google-earth.kmz ADD_DISPOSITION"
				},
				{
				  "name": "READER_DATASET_HINT",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "Select the Google Earth KML file(s)"
				},
				{
				  "name": "ADVANCED_PARMS",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "OGCKML_OUT_OUTPUT_SCHEMA OGCKML_OUT_AUTO_CREATE_NETWORK_LINKS OGCKML_OUT_LOG_VERBOSE OGCKML_OUT_ORIENTATION OGCKML_OUT_DATASET_HINT OGCKML_OUT_STYLE_DOC OGCKML_OUT_SCHEMA_DOC OGCKML_OUT_DETECT_RASTERS OGCKML_OUT_RASTER_MODE OGCKML_OUT_RASTER_FORMAT OGCKML_OUT_TEXTURE_FORMAT OGCKML_OUT_COPY_ICON OGCKML_OUT_REGIONATE_DATA OGCKML_OUT_EXEC_GO_PIPELINE OGCKML_OUT_EXEC_PO_PIPELINE OGCKML_OUT_OMIT_DOCUMENT_ELEMENT OGCKML_OUT_CREATE_EMPTY_FOLDERS OGCKML_OUT_KML21_TARGET_HREF OGCKML_OUT_KML21_FANOUT_TYPE OGCKML_OUT_MOVE_TO_KML_LOCAL_COORDSYS OGCKML_OUT_WRITE_3D_GEOM_AS_POLYGONS OGCKML_OUT_WRITE_TEXTURES_TXT_FILE KML21_INFORMATION_POINT_ICON KML21_REGIONATOR_PIPELINE KML21_GO_PYRAMIDER_PIPELINE KML21_PO_PYRAMIDER_PIPELINE"
				},
				{
				  "name": "OGCKML_1",
				  "attributes": {},
				  "category": "KEYWORD_SUFFIX",
				  "value": "austinApartments.kml"
				},
				{
				  "name": "OGCKML_1",
				  "attributes": {},
				  "category": "DEFAULT_OVERRIDE",
				  "value": "-OGCKML_1_DATASET"
				},
				{
				  "name": "NETWORK_PROXY",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "NO"
				},
				{
				  "name": "WRITER_DATASET_HINT",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "Specify a name for the Google Earth KML File"
				},
				{
				  "name": "DATASET_NAME",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "kml file"
				},
				{
				  "name": "OGCKML_1",
				  "attributes": {},
				  "category": "FANOUT_GROUP",
				  "value": "NO"
				},
				{
				  "name": "NETWORK_AUTHENTICATION",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "ALWAYS"
				},
				{
				  "name": "OGCKML_1",
				  "attributes": {},
				  "category": "DATASET_TYPE",
				  "value": "FILE_OR_URL"
				},
				{
				  "name": "SUPPORTS_ESTABLISHED_CACHE",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "YES"
				}
			  ]
			}
		  ],
		  "source": [
			{
			  "format": "SPATIALITE_NATIVE",
			  "name": "SPATIALITE_NATIVE_1",
			  "location": "$(FME_MF_DIR_USERTYPED)\\Landmarks.sqlite",
			  "source": true,
			  "featuretypes": [
				{
				  "name": "Landmarks",
				  "description": "",
				  "attributes": [
					{
					  "decimals": 0,
					  "name": "CFCC",
					  "width": 4,
					  "type": "varchar"
					},
					{
					  "decimals": 0,
					  "name": "FILE",
					  "width": 0,
					  "type": "integer"
					},
					{
					  "decimals": 0,
					  "name": "LALAT",
					  "width": 0,
					  "type": "integer"
					},
					{
					  "decimals": 0,
					  "name": "LALONG",
					  "width": 0,
					  "type": "integer"
					},
					{
					  "decimals": 0,
					  "name": "LANAME",
					  "width": 31,
					  "type": "varchar"
					},
					{
					  "decimals": 0,
					  "name": "LAND",
					  "width": 0,
					  "type": "integer"
					},
					{
					  "decimals": 0,
					  "name": "MODULE",
					  "width": 9,
					  "type": "varchar"
					},
					{
					  "decimals": 0,
					  "name": "OBJECTID",
					  "width": 0,
					  "type": "integer"
					},
					{
					  "decimals": 0,
					  "name": "SOURCE",
					  "width": 2,
					  "type": "varchar"
					}
				  ],
				  "properties": []
				}
			  ],
			  "properties": [
				{
				  "name": "SUPPORTS_SCHEMA_IN_FEATURE_TYPE_NAME",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "NO"
				},
				{
				  "name": "FEATURE_TYPE_NAME",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "Table"
				},
				{
				  "name": "SPATIALITE_NATIVE_1",
				  "attributes": {},
				  "category": "DEFAULT_DATASET_PATH",
				  "value": "$(FME_MF_DIR_USERTYPED)\\Landmarks.sqlite"
				},
				{
				  "name": "DATASET_NAME",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "spatialite file"
				},
				{
				  "name": "SPATIALITE_NATIVE_1",
				  "attributes": {},
				  "category": "DATASET_DESCRIPTION",
				  "value": "Landmarks.sqlite [SPATIALITE_NATIVE]"
				},
				{
				  "name": "ADVANCED_PARMS",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "SPATIALITE_NATIVE_IN_BEGIN_SQL SPATIALITE_NATIVE_IN_END_SQL SPATIALITE_NATIVE_OUT_WRITER_MODE SPATIALITE_NATIVE_OUT_TRANSACTION_INTERVAL SPATIALITE_NATIVE_OUT_BEGIN_SQL{0} SPATIALITE_NATIVE_OUT_END_SQL{0}"
				},
				{
				  "name": "WRITER_DATASET_HINT",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "Specify a name for the SpatiaLite Database file"
				},
				{
				  "name": "ALLOW_DATASET_CONFLICT",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "YES"
				},
				{
				  "name": "ATTRIBUTE_READING_HISTORIC",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "ALL"
				},
				{
				  "name": "FEATURE_TYPE_DEFAULT_NAME",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "Table1"
				},
				{
				  "name": "SPATIALITE_NATIVE_1",
				  "attributes": {},
				  "category": "DATASET_TYPE",
				  "value": "FILE"
				},
				{
				  "name": "ATTRIBUTE_READING",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "DEFLINE"
				},
				{
				  "name": "SPATIALITE_NATIVE_1",
				  "attributes": {},
				  "category": "KEYWORD_SUFFIX",
				  "value": "Landmarks.sqlite"
				},
				{
				  "name": "READER_DATASET_HINT",
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "value": "Select the SpatiaLite Database file(s)"
				},
				{
				  "name": "SPATIALITE_NATIVE_1",
				  "attributes": {},
				  "category": "DEFAULT_OVERRIDE",
				  "value": "-SPATIALITE_NATIVE_1_DATASET"
				}
			  ]
			}
		  ]
		},
		"title": "City of Austin: Apartments and other (SPATIALITE 2 KML)",
		"type": "WORKSPACE",
		"buildNumber": 22337,
		"enabled": true,
		"avgCpuTime": 728,
		"lastPublishDate": "2022-12-07T18:50:52Z",
		"lastSaveBuild": "FME(R) 2022.0.0.0 (20220428 - Build 22337 - WIN64)",
		"avgElapsedTime": 784,
		"lastSaveDate": "2022-06-15T10:59:24Z",
		"requirements": "",
		"resources": [
		  {
			"size": 6680576,
			"name": "Landmarks.sqlite",
			"description": ""
		  }
		],
		"history": "",
		"services": [
		  {
			"displayName": "Data Download",
			"name": "fmedatadownload"
		  },
		  {
			"displayName": "Data Streaming",
			"name": "fmedatastreaming"
		  },
		  {
			"displayName": "Job Submitter",
			"name": "fmejobsubmitter"
		  },
		  {
			"displayName": "KML Network Link",
			"name": "fmekmllink"
		  }
		],
		"userName": "admin",
		"requirementsKeyword": "",
		"fileSize": 136390,
		"name": "austinApartments.fmw",
		"totalRuns": 2,
		"category": "",
		"parameters": [],
		"properties": [
		  {
			"name": "NOTIFICATION_WRITER",
			"attributes": {},
			"category": "fmejobsubmitter_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "OUTPUT_WRITER",
			"attributes": {},
			"category": "fmedatastreaming_FMEUSERPROPDATA",
			"value": "OGCKML_1"
		  },
		  {
			"name": "FAILURE_TOPICS",
			"attributes": {},
			"category": "fmedatadownload_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "FAILURE_TOPICS",
			"attributes": {},
			"category": "fmedatastreaming_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "NetworkLink/Region/Lod/minLodPixels",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": "0"
		  },
		  {
			"name": "NetworkLink/Link/refreshMode",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": "onInterval"
		  },
		  {
			"name": "NetworkLink/Region/LatLonAltBox/east",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "HTTP_DATASET",
			"attributes": {},
			"category": "fmedatastreaming_FMEUSERPROPDATA",
			"value": "SPATIALITE_NATIVE_1"
		  },
		  {
			"name": "NetworkLink/Link/refreshInterval",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": "60"
		  },
		  {
			"name": "NetworkLink/Region/Lod/maxLodPixels",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": "-1"
		  },
		  {
			"name": "NOTIFICATION_WRITER",
			"attributes": {},
			"category": "fmedatadownload_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "ADVANCED",
			"attributes": {},
			"category": "fmejobsubmitter_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "NetworkLink/Link/viewRefreshTime",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": "7"
		  },
		  {
			"name": "SUCCESS_TOPICS",
			"attributes": {},
			"category": "fmedatastreaming_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "OUTPUT_WRITER",
			"attributes": {},
			"category": "fmedatadownload_FMEUSERPROPDATA",
			"value": "OGCKML_1"
		  },
		  {
			"name": "ADVANCED",
			"attributes": {},
			"category": "fmedatastreaming_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "NetworkLink/visibility",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": "0"
		  },
		  {
			"name": "ADVANCED",
			"attributes": {},
			"category": "fmedatadownload_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "NetworkLink/Link/viewRefreshMode",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": "onRegion"
		  },
		  {
			"name": "OUTPUT_WRITER",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": "OGCKML_1"
		  },
		  {
			"name": "NetworkLink/description",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": "A KML Link document processed by FME Server"
		  },
		  {
			"name": "NetworkLink/Link/viewFormat",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": "BBOX=[bboxWest],[bboxSouth],[bboxEast],[bboxNorth]"
		  },
		  {
			"name": "HTTP_DATASET",
			"attributes": {},
			"category": "fmejobsubmitter_FMEUSERPROPDATA",
			"value": "SPATIALITE_NATIVE_1"
		  },
		  {
			"name": "HTTP_DATASET",
			"attributes": {},
			"category": "fmedatadownload_FMEUSERPROPDATA",
			"value": "SPATIALITE_NATIVE_1"
		  },
		  {
			"name": "NOTIFICATION_WRITER",
			"attributes": {},
			"category": "fmedatastreaming_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "NetworkLink/Region/LatLonAltBox/west",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "NetworkLink/name",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": "KML Link Translation"
		  },
		  {
			"name": "NetworkLink/Region/LatLonAltBox/north",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "NetworkLink/Region/LatLonAltBox/south",
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "SUCCESS_TOPICS",
			"attributes": {},
			"category": "fmejobsubmitter_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "FAILURE_TOPICS",
			"attributes": {},
			"category": "fmejobsubmitter_FMEUSERPROPDATA",
			"value": ""
		  },
		  {
			"name": "SUCCESS_TOPICS",
			"attributes": {},
			"category": "fmedatadownload_FMEUSERPROPDATA",
			"value": ""
		  }
		]
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"workspaces", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"workspaces", "--repository", "Samples"},
		},
		{
			name:        "required flag not set",
			wantErrText: "required flag(s) \"repository\" not set",
			args:        []string{"workspaces", "--name", "austinApartments.fmw"},
		},
		{
			name:        "repository does not exist",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found: check that the specified repository exists",
			args:        []string{"workspaces", "--repository", "Samples123"},
		},
		{
			name:        "workspace does not exist",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found: check that the specified repository and workspace exists",
			args:        []string{"workspaces", "--repository", "Samples", "--name", "austinAprtmnts.fmw"},
		},
		{
			name:            "get workspaces table output",
			statusCode:      http.StatusOK,
			args:            []string{"workspaces", "--repository", "Samples"},
			body:            responseV3,
			wantOutputRegex: "^[\\s]*NAME[\\s]*TITLE[\\s]*LAST SAVE DATE[\\s]*austinApartments.fmw[\\s]*City of Austin: Apartments and other \\(SPATIALITE 2 KML\\)[\\s]*2022-06-15 10:59:24 \\+0000 UTC[\\s]*austinDownload.fmw[\\s]*City of Austin: Data Download[\\s]*2022-06-15 13:40:29 \\+0000 UTC[\\s]*earthquakesextrusion.fmw[\\s]*Earthquakes: GeoJSON to KML Diagrams[\\s]*2022-06-15 15:52:07 \\+0000 UTC[\\s]*easyTranslator.fmw[\\s]*Generic format translator[\\s]*2022-06-15 15:56:11 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get single workspace",
			statusCode:      http.StatusOK,
			body:            responseV3SingleWorkspace,
			args:            []string{"workspaces", "--repository", "Samples", "--name", "austinApartmets.fmw"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*TITLE[\\s]*LAST SAVE DATE[\\s]*austinApartments.fmw[\\s]*City of Austin: Apartments and other \\(SPATIALITE 2 KML\\)[\\s]*2022-06-15 10:59:24 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get single workspace custom columns",
			statusCode:      http.StatusOK,
			body:            responseV3SingleWorkspace,
			args:            []string{"workspaces", "--repository", "Samples", "--name", "austinApartments.fmw", "--output=custom-columns=NAME:$.name,SOURCE:$.datasets.source[*].format,DEST:$.datasets.destination[*].format,SERVICES:$.services[*].name"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*SOURCE[\\s]*DEST[\\s]*SERVICES[\\s]*austinApartments.fmw[\\s]*\\[SPATIALITE_NATIVE\\][\\s]*\\[OGCKML\\][\\s]*\\[fmedatadownload fmedatastreaming fmejobsubmitter fmekmllink\\][\\s]*$",
		},
	}

	runTests(cases, t)

}
