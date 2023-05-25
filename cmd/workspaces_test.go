package cmd

import (
	"net/http"
	"testing"
)

func TestWorkspaces(t *testing.T) {
	responseV4 := `{   
		"items": [
		  {
			"name": "austinApartments.fmw",
			"title": "City of Austin: Apartments and other (SPATIALITE 2 KML)",
			"description": "",
			"type": "workspace",
			"repositoryName": "Samples",
			"lastSaveDate": "2022-06-15T10:59:24.000Z",
			"lastPublishDate": "2022-12-07T18:50:52.635Z",
			"lastPublishUserId": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"lastPublishUser": "admin",
			"totalFileSize": 6680576,
			"fileCount": 1,
			"averageCpuPercent": 92.85714285714286,
			"averagePeakMemoryUsage": 3123940,
			"averageCpuTime": 728,
			"totalRuns": 2,
			"averageElapsedTime": 784,
			"favorite": false
		  },
		  {
			"name": "austinDownload.fmw",
			"title": "City of Austin: Data Download",
			"description": "",
			"type": "workspace",
			"repositoryName": "Samples",
			"lastSaveDate": "2022-06-15T13:40:29.000Z",
			"lastPublishDate": "2022-12-07T18:50:52.833Z",
			"lastPublishUserId": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"lastPublishUser": "admin",
			"totalFileSize": 23596057,
			"fileCount": 30,
			"averageCpuPercent": 71.16968698517299,
			"averagePeakMemoryUsage": 45244856,
			"averageCpuTime": 1728,
			"totalRuns": 1,
			"averageElapsedTime": 2428,
			"favorite": false
		  },
		  {
			"name": "AverageRunningTime.fmw",
			"title": "",
			"description": "",
			"type": "workspace",
			"repositoryName": "Dashboards",
			"lastSaveDate": "2022-01-31T06:57:03.000Z",
			"lastPublishDate": "2022-12-07T18:50:51.836Z",
			"lastPublishUserId": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"lastPublishUser": "admin",
			"totalFileSize": 0,
			"fileCount": 0,
			"averageCpuPercent": 0.0,
			"averagePeakMemoryUsage": 0,
			"averageCpuTime": 0,
			"totalRuns": 0,
			"averageElapsedTime": 0,
			"favorite": false
		  },
		  {
			"name": "backupConfiguration.fmw",
			"title": "",
			"description": "",
			"type": "workspace",
			"repositoryName": "Utilities",
			"lastSaveDate": "2022-01-31T05:49:24.000Z",
			"lastPublishDate": "2022-12-07T18:50:51.365Z",
			"lastPublishUserId": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"lastPublishUser": "admin",
			"totalFileSize": 0,
			"fileCount": 0,
			"averageCpuPercent": 0.0,
			"averagePeakMemoryUsage": 0,
			"averageCpuTime": 0,
			"totalRuns": 0,
			"averageElapsedTime": 0,
			"favorite": false
		  },
		  {
			"name": "DailyAverageQueuedTime.fmw",
			"title": "",
			"description": "",
			"type": "workspace",
			"repositoryName": "Dashboards",
			"lastSaveDate": "2022-01-31T06:56:29.000Z",
			"lastPublishDate": "2022-12-07T18:50:52.228Z",
			"lastPublishUserId": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"lastPublishUser": "admin",
			"totalFileSize": 0,
			"fileCount": 0,
			"averageCpuPercent": 0.0,
			"averagePeakMemoryUsage": 0,
			"averageCpuTime": 0,
			"totalRuns": 0,
			"averageElapsedTime": 0,
			"favorite": false
		  },
		  {
			"name": "DailyTotalRunningTime.fmw",
			"title": "",
			"description": "",
			"type": "workspace",
			"repositoryName": "Dashboards",
			"lastSaveDate": "2022-01-31T06:55:16.000Z",
			"lastPublishDate": "2022-12-07T18:50:52.368Z",
			"lastPublishUserId": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"lastPublishUser": "admin",
			"totalFileSize": 0,
			"fileCount": 0,
			"averageCpuPercent": 0.0,
			"averagePeakMemoryUsage": 0,
			"averageCpuTime": 0,
			"totalRuns": 0,
			"averageElapsedTime": 0,
			"favorite": false
		  }
		],
		"totalCount": 6,
		"limit": 100,
		"offset": 0
	  }`

	responseSamplesV4 := `{
		"items": [
		  {
			"name": "austinApartments.fmw",
			"title": "City of Austin: Apartments and other (SPATIALITE 2 KML)",
			"description": "",
			"type": "workspace",
			"repositoryName": "Samples",
			"lastSaveDate": "2022-06-15T10:59:24.000Z",
			"lastPublishDate": "2022-12-07T18:50:52.635Z",
			"lastPublishUserId": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"lastPublishUser": "admin",
			"totalFileSize": 6680576,
			"fileCount": 1,
			"averageCpuPercent": 92.85714285714286,
			"averagePeakMemoryUsage": 3123940,
			"averageCpuTime": 728,
			"totalRuns": 2,
			"averageElapsedTime": 784,
			"favorite": false
		  },
		  {
			"name": "austinDownload.fmw",
			"title": "City of Austin: Data Download",
			"description": "",
			"type": "workspace",
			"repositoryName": "Samples",
			"lastSaveDate": "2022-06-15T13:40:29.000Z",
			"lastPublishDate": "2022-12-07T18:50:52.833Z",
			"lastPublishUserId": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"lastPublishUser": "admin",
			"totalFileSize": 23596057,
			"fileCount": 30,
			"averageCpuPercent": 71.16968698517299,
			"averagePeakMemoryUsage": 45244856,
			"averageCpuTime": 1728,
			"totalRuns": 1,
			"averageElapsedTime": 2428,
			"favorite": false
		  },
		  {
			"name": "earthquakesextrusion.fmw",
			"title": "Earthquakes: GeoJSON to KML Diagrams",
			"description": "",
			"type": "workspace",
			"repositoryName": "Samples",
			"lastSaveDate": "2022-06-15T15:52:07.000Z",
			"lastPublishDate": "2022-12-07T18:50:53.682Z",
			"lastPublishUserId": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"lastPublishUser": "admin",
			"totalFileSize": 0,
			"fileCount": 0,
			"averageCpuPercent": 63.57292393579902,
			"averagePeakMemoryUsage": 3828352,
			"averageCpuTime": 911,
			"totalRuns": 1,
			"averageElapsedTime": 1433,
			"favorite": false
		  },
		  {
			"name": "easyTranslator.fmw",
			"title": "Generic format translator",
			"description": "",
			"type": "workspace",
			"repositoryName": "Samples",
			"lastSaveDate": "2022-06-15T15:56:11.000Z",
			"lastPublishDate": "2022-12-07T18:50:53.840Z",
			"lastPublishUserId": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"lastPublishUser": "admin",
			"totalFileSize": 0,
			"fileCount": 0,
			"averageCpuPercent": 0.0,
			"averagePeakMemoryUsage": 0,
			"averageCpuTime": 0,
			"totalRuns": 0,
			"averageElapsedTime": 0,
			"favorite": false
		  }
		],
		"totalCount": 4,
		"limit": 100,
		"offset": 0
	  }`

	responseV4SingleWorkspace := `{
		"name": "austinApartments.fmw",
		"title": "City of Austin: Apartments and other (SPATIALITE 2 KML)",
		"type": "workspace",
		"buildNumber": 22337,
		"category": "",
		"history": "",
		"lastSaveDate": "2022-06-15T10:59:24.000Z",
		"lastPublishDate": "2022-12-07T18:50:52.635Z",
		"lastSaveBuild": "FME(R) 2022.0.0.0 (20220428 - Build 22337 - WIN64)",
		"legalTermsConditions": "",
		"properties": [
		  {
			"attributes": {},
			"category": "fmedatadownload_FMEUSERPROPDATA",
			"name": "ADVANCED",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmedatadownload_FMEUSERPROPDATA",
			"name": "FAILURE_TOPICS",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmedatadownload_FMEUSERPROPDATA",
			"name": "HTTP_DATASET",
			"value": "SPATIALITE_NATIVE_1"
		  },
		  {
			"attributes": {},
			"category": "fmedatadownload_FMEUSERPROPDATA",
			"name": "NOTIFICATION_WRITER",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmedatadownload_FMEUSERPROPDATA",
			"name": "OUTPUT_WRITER",
			"value": "OGCKML_1"
		  },
		  {
			"attributes": {},
			"category": "fmedatadownload_FMEUSERPROPDATA",
			"name": "SUCCESS_TOPICS",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmedatastreaming_FMEUSERPROPDATA",
			"name": "ADVANCED",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmedatastreaming_FMEUSERPROPDATA",
			"name": "FAILURE_TOPICS",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmedatastreaming_FMEUSERPROPDATA",
			"name": "HTTP_DATASET",
			"value": "SPATIALITE_NATIVE_1"
		  },
		  {
			"attributes": {},
			"category": "fmedatastreaming_FMEUSERPROPDATA",
			"name": "NOTIFICATION_WRITER",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmedatastreaming_FMEUSERPROPDATA",
			"name": "OUTPUT_WRITER",
			"value": "OGCKML_1"
		  },
		  {
			"attributes": {},
			"category": "fmedatastreaming_FMEUSERPROPDATA",
			"name": "SUCCESS_TOPICS",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmejobsubmitter_FMEUSERPROPDATA",
			"name": "ADVANCED",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmejobsubmitter_FMEUSERPROPDATA",
			"name": "FAILURE_TOPICS",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmejobsubmitter_FMEUSERPROPDATA",
			"name": "HTTP_DATASET",
			"value": "SPATIALITE_NATIVE_1"
		  },
		  {
			"attributes": {},
			"category": "fmejobsubmitter_FMEUSERPROPDATA",
			"name": "NOTIFICATION_WRITER",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmejobsubmitter_FMEUSERPROPDATA",
			"name": "SUCCESS_TOPICS",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/description",
			"value": "A KML Link document processed by FME Server"
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/Link/refreshInterval",
			"value": "60"
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/Link/refreshMode",
			"value": "onInterval"
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/Link/viewFormat",
			"value": "BBOX=[bboxWest],[bboxSouth],[bboxEast],[bboxNorth]"
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/Link/viewRefreshMode",
			"value": "onRegion"
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/Link/viewRefreshTime",
			"value": "7"
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/name",
			"value": "KML Link Translation"
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/Region/LatLonAltBox/east",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/Region/LatLonAltBox/north",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/Region/LatLonAltBox/south",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/Region/LatLonAltBox/west",
			"value": ""
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/Region/Lod/maxLodPixels",
			"value": "-1"
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/Region/Lod/minLodPixels",
			"value": "0"
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "NetworkLink/visibility",
			"value": "0"
		  },
		  {
			"attributes": {},
			"category": "fmekmllink_FMEUSERPROPDATA",
			"name": "OUTPUT_WRITER",
			"value": "OGCKML_1"
		  }
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
		"averageCpuPercent": 92.85714285714286,
		"averagePeakMemoryUsage": 3123940,
		"averageCpuTime": 728,
		"totalRuns": 2,
		"averageElapsedTime": 784,
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
					{
					  "decimals": 0,
					  "name": "FILE",
					  "type": "integer",
					  "width": 0
					},
					{
					  "decimals": 0,
					  "name": "LALAT",
					  "type": "integer",
					  "width": 0
					},
					{
					  "decimals": 0,
					  "name": "LALONG",
					  "type": "integer",
					  "width": 0
					},
					{
					  "decimals": 0,
					  "name": "LANAME",
					  "type": "varchar",
					  "width": 31
					},
					{
					  "decimals": 0,
					  "name": "LAND",
					  "type": "integer",
					  "width": 0
					},
					{
					  "decimals": 0,
					  "name": "MODULE",
					  "type": "varchar",
					  "width": 9
					},
					{
					  "decimals": 0,
					  "name": "OBJECTID",
					  "type": "integer",
					  "width": 0
					},
					{
					  "decimals": 0,
					  "name": "SOURCE",
					  "type": "varchar",
					  "width": 2
					}
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
				{
				  "attributes": {},
				  "category": "DATASET_TYPE",
				  "name": "SPATIALITE_NATIVE_1",
				  "value": "FILE"
				},
				{
				  "attributes": {},
				  "category": "DEFAULT_DATASET_PATH",
				  "name": "SPATIALITE_NATIVE_1",
				  "value": "$(FME_MF_DIR_USERTYPED)\\Landmarks.sqlite"
				},
				{
				  "attributes": {},
				  "category": "DEFAULT_OVERRIDE",
				  "name": "SPATIALITE_NATIVE_1",
				  "value": "-SPATIALITE_NATIVE_1_DATASET"
				},
				{
				  "attributes": {},
				  "category": "KEYWORD_SUFFIX",
				  "name": "SPATIALITE_NATIVE_1",
				  "value": "Landmarks.sqlite"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "ADVANCED_PARMS",
				  "value": "SPATIALITE_NATIVE_IN_BEGIN_SQL SPATIALITE_NATIVE_IN_END_SQL SPATIALITE_NATIVE_OUT_WRITER_MODE SPATIALITE_NATIVE_OUT_TRANSACTION_INTERVAL SPATIALITE_NATIVE_OUT_BEGIN_SQL{0} SPATIALITE_NATIVE_OUT_END_SQL{0}"      
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "ALLOW_DATASET_CONFLICT",
				  "value": "YES"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "ATTRIBUTE_READING",
				  "value": "DEFLINE"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "ATTRIBUTE_READING_HISTORIC",
				  "value": "ALL"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "DATASET_NAME",
				  "value": "spatialite file"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "FEATURE_TYPE_DEFAULT_NAME",
				  "value": "Table1"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "FEATURE_TYPE_NAME",
				  "value": "Table"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "READER_DATASET_HINT",
				  "value": "Select the SpatiaLite Database file(s)"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "SUPPORTS_SCHEMA_IN_FEATURE_TYPE_NAME",
				  "value": "NO"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "WRITER_DATASET_HINT",
				  "value": "Specify a name for the SpatiaLite Database file"
				}
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
					{
					  "decimals": 0,
					  "name": "FILE",
					  "type": "kml_char",
					  "width": 20
					},
					{
					  "decimals": 0,
					  "name": "LALAT",
					  "type": "kml_char",
					  "width": 20
					},
					{
					  "decimals": 0,
					  "name": "LALONG",
					  "type": "kml_char",
					  "width": 20
					},
					{
					  "decimals": 0,
					  "name": "LANAME",
					  "type": "kml_char",
					  "width": 31
					},
					{
					  "decimals": 0,
					  "name": "LAND",
					  "type": "kml_char",
					  "width": 20
					},
					{
					  "decimals": 0,
					  "name": "MODULE",
					  "type": "kml_char",
					  "width": 9
					},
					{
					  "decimals": 0,
					  "name": "OBJECTID",
					  "type": "kml_char",
					  "width": 20
					},
					{
					  "decimals": 0,
					  "name": "SOURCE",
					  "type": "kml_char",
					  "width": 2
					}
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
				{
				  "attributes": {},
				  "category": "DATASET_TYPE",
				  "name": "OGCKML_1",
				  "value": "FILE_OR_URL"
				},
				{
				  "attributes": {},
				  "category": "DEFAULT_DATASET_PATH",
				  "name": "OGCKML_1",
				  "value": "OGCKML_1/austinApartments.kml"
				},
				{
				  "attributes": {},
				  "category": "DEFAULT_OVERRIDE",
				  "name": "OGCKML_1",
				  "value": "-OGCKML_1_DATASET"
				},
				{
				  "attributes": {},
				  "category": "FANOUT_GROUP",
				  "name": "OGCKML_1",
				  "value": "NO"
				},
				{
				  "attributes": {},
				  "category": "KEYWORD_SUFFIX",
				  "name": "OGCKML_1",
				  "value": "austinApartments.kml"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "ADVANCED_PARMS",
				  "value": "OGCKML_OUT_OUTPUT_SCHEMA OGCKML_OUT_AUTO_CREATE_NETWORK_LINKS OGCKML_OUT_LOG_VERBOSE OGCKML_OUT_ORIENTATION OGCKML_OUT_DATASET_HINT OGCKML_OUT_STYLE_DOC OGCKML_OUT_SCHEMA_DOC OGCKML_OUT_DETECT_RASTERS OGCKML_OUT_RASTER_MODE OGCKML_OUT_RASTER_FORMAT OGCKML_OUT_TEXTURE_FORMAT OGCKML_OUT_COPY_ICON OGCKML_OUT_REGIONATE_DATA OGCKML_OUT_EXEC_GO_PIPELINE OGCKML_OUT_EXEC_PO_PIPELINE OGCKML_OUT_OMIT_DOCUMENT_ELEMENT OGCKML_OUT_CREATE_EMPTY_FOLDERS OGCKML_OUT_KML21_TARGET_HREF OGCKML_OUT_KML21_FANOUT_TYPE OGCKML_OUT_MOVE_TO_KML_LOCAL_COORDSYS OGCKML_OUT_WRITE_3D_GEOM_AS_POLYGONS OGCKML_OUT_WRITE_TEXTURES_TXT_FILE KML21_INFORMATION_POINT_ICON KML21_REGIONATOR_PIPELINE KML21_GO_PYRAMIDER_PIPELINE KML21_PO_PYRAMIDER_PIPELINE"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "DATASET_NAME",
				  "value": "kml file"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "MIME_TYPE",
				  "value": ".kml application/vnd.google-earth.kml+xml .kmz application/vnd.google-earth.kmz ADD_DISPOSITION"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "NETWORK_AUTHENTICATION",
				  "value": "ALWAYS"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "NETWORK_PROXY",
				  "value": "NO"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "READER_DATASET_HINT",
				  "value": "Select the Google Earth KML file(s)"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "SUPPORTS_ESTABLISHED_CACHE",
				  "value": "YES"
				},
				{
				  "attributes": {},
				  "category": "METAFILE_PARAMETER",
				  "name": "WRITER_DATASET_HINT",
				  "value": "Specify a name for the Google Earth KML File"
				}
			  ],
			  "source": false
			}
		  ]
		},
		"services": {
		  "dataDownload": {
			"registered": true,
			"reader": "SPATIALITE_NATIVE_1",
			"writers": [
			  "OGCKML_1"
			],
			"zipLayout": {}
		  },
		  "dataStreaming": {
			"registered": true,
			"reader": "SPATIALITE_NATIVE_1",
			"writers": [
			  "OGCKML_1"
			]
		  },
		  "jobSubmitter": {
			"registered": true,
			"reader": "SPATIALITE_NATIVE_1"
		  },
		  "kmlNetworkLink": {
			"registered": true,
			"writers": [
			  "OGCKML_1"
			],
			"name": "KML Link Translation",
			"visibility": "0",
			"description": "A KML Link document processed by FME Server",
			"link": {
			  "viewRefreshMode": "onRegion",
			  "viewRefreshTime": 7,
			  "viewFormat": "BBOX=[bboxWest],[bboxSouth],[bboxEast],[bboxNorth]",
			  "refreshMode": "onInterval",
			  "viewRefreshInterval": 60
			},
			"lod": {
			  "minLodPixels": 0,
			  "maxLodPixels": -1
			}
		  }
		},
		"parameters": []
	  }`

	repoNotFoundV4 := `{
		"message": "Unauthorized request by user admin due to lack of proper permissions or the object does not exist."
	  }`
	workspaceNotFoundV4 := `{
		"message": "Item test.fmw does not exist."
	  }`

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
			  "location": "/data/fmeflowdata/resources/system/temp/engineresults\\austinApartments.kml",
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
			name:        "repository does not exist V4",
			statusCode:  http.StatusNotFound,
			body:        repoNotFoundV4,
			wantErrText: "Unauthorized request by user admin due to lack of proper permissions or the object does not exist.",
			args:        []string{"workspaces", "--repository", "Samples123", "--api-version", "v4"},
		},
		{
			name:        "workspace does not exist V4",
			statusCode:  http.StatusNotFound,
			body:        workspaceNotFoundV4,
			wantErrText: "Item test.fmw does not exist.",
			args:        []string{"workspaces", "--repository", "Samples", "--name", "test.fmw", "--api-version", "v4"},
		},
		{
			name:            "get workspaces table output V4",
			statusCode:      http.StatusOK,
			args:            []string{"workspaces", "--api-version", "v4"},
			body:            responseV4,
			wantOutputRegex: "^[\\s]*NAME[\\s]*TITLE[\\s]*LAST SAVE DATE[\\s]*austinApartments.fmw[\\s]*City of Austin: Apartments and other \\(SPATIALITE 2 KML\\)[\\s]*2022-06-15 10:59:24 \\+0000 UTC[\\s]*austinDownload.fmw[\\s]*City of Austin: Data Download[\\s]*2022-06-15 13:40:29 \\+0000 UTC [\\s]*AverageRunningTime.fmw[\\s]*2022-01-31 06:57:03 \\+0000 UTC[\\s]*backupConfiguration.fmw[\\s]*2022-01-31 05:49:24 \\+0000 UTC[\\s]*DailyAverageQueuedTime.fmw[\\s]*2022-01-31 06:56:29 \\+0000 UTC[\\s]*DailyTotalRunningTime.fmw[\\s]*2022-01-31 06:55:16 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get workspaces single repo table output V4",
			statusCode:      http.StatusOK,
			args:            []string{"workspaces", "--repository", "Samples", "--api-version", "v4"},
			body:            responseSamplesV4,
			wantOutputRegex: "^[\\s]*NAME[\\s]*TITLE[\\s]*LAST SAVE DATE[\\s]*austinApartments.fmw[\\s]*City of Austin: Apartments and other \\(SPATIALITE 2 KML\\)[\\s]*2022-06-15 10:59:24 \\+0000 UTC[\\s]*austinDownload.fmw[\\s]*City of Austin: Data Download[\\s]*2022-06-15 13:40:29 \\+0000 UTC[\\s]*earthquakesextrusion.fmw[\\s]*Earthquakes: GeoJSON to KML Diagrams[\\s]*2022-06-15 15:52:07 \\+0000 UTC[\\s]*easyTranslator.fmw[\\s]*Generic format translator[\\s]*2022-06-15 15:56:11 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get single workspace V4",
			statusCode:      http.StatusOK,
			body:            responseV4SingleWorkspace,
			args:            []string{"workspaces", "--repository", "Samples", "--name", "austinApartments.fmw", "--api-version", "v4"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*TITLE[\\s]*LAST SAVE DATE[\\s]*austinApartments.fmw[\\s]*City of Austin: Apartments and other \\(SPATIALITE 2 KML\\)[\\s]*2022-06-15 10:59:24 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get single workspace custom columns V4",
			statusCode:      http.StatusOK,
			body:            responseV4SingleWorkspace,
			args:            []string{"workspaces", "--repository", "Samples", "--name", "austinApartments.fmw", "--output=custom-columns=NAME:.name,SOURCE:.datasets.source[*].format,DEST:.datasets.destination[*].format,BUILD:.buildNumber", "--api-version", "v4"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*SOURCE[\\s]*DEST[\\s]*BUILD[\\s]*austinApartments.fmw[\\s]*SPATIALITE_NATIVE[\\s]*OGCKML[\\s]*22337[\\s]*$",
		},
		{
			name:        "repository does not exist V3",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found: check that the specified repository exists",
			args:        []string{"workspaces", "--repository", "Samples123", "--api-version", "v3"},
		},
		{
			name:        "workspace does not exist V3",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found: check that the specified repository and workspace exists",
			args:        []string{"workspaces", "--repository", "Samples", "--name", "austinAprtmnts.fmw", "--api-version", "v3"},
		},
		{
			name:            "get workspaces table output V3",
			statusCode:      http.StatusOK,
			args:            []string{"workspaces", "--repository", "Samples", "--api-version", "v3"},
			body:            responseV3,
			wantOutputRegex: "^[\\s]*NAME[\\s]*TITLE[\\s]*LAST SAVE DATE[\\s]*austinApartments.fmw[\\s]*City of Austin: Apartments and other \\(SPATIALITE 2 KML\\)[\\s]*2022-06-15 10:59:24 \\+0000 UTC[\\s]*austinDownload.fmw[\\s]*City of Austin: Data Download[\\s]*2022-06-15 13:40:29 \\+0000 UTC[\\s]*earthquakesextrusion.fmw[\\s]*Earthquakes: GeoJSON to KML Diagrams[\\s]*2022-06-15 15:52:07 \\+0000 UTC[\\s]*easyTranslator.fmw[\\s]*Generic format translator[\\s]*2022-06-15 15:56:11 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get single workspace V3",
			statusCode:      http.StatusOK,
			body:            responseV3SingleWorkspace,
			args:            []string{"workspaces", "--repository", "Samples", "--name", "austinApartmets.fmw", "--api-version", "v3"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*TITLE[\\s]*LAST SAVE DATE[\\s]*austinApartments.fmw[\\s]*City of Austin: Apartments and other \\(SPATIALITE 2 KML\\)[\\s]*2022-06-15 10:59:24 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get single workspace custom columns V3",
			statusCode:      http.StatusOK,
			body:            responseV3SingleWorkspace,
			args:            []string{"workspaces", "--repository", "Samples", "--name", "austinApartments.fmw", "--output=custom-columns=NAME:.name,SOURCE:.datasets.source[*].format,DEST:.datasets.destination[*].format,SERVICES:.services[*].name", "--api-version", "v3"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*SOURCE[\\s]*DEST[\\s]*SERVICES[\\s]*austinApartments.fmw[\\s]*SPATIALITE_NATIVE[\\s]*OGCKML[\\s]*fmedatadownload fmedatastreaming fmejobsubmitter fmekmllink[\\s]*$",
		},
	}

	runTests(cases, t)

}
