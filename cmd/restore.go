package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ResourceRequestV4 struct {
	ResourceName       string `json:"resourceName"`
	PackagePath        string `json:"packagePath"`
	Overwrite          bool   `json:"overwrite"`
	PauseNotifications bool   `json:"pauseNotifications"`
	SuccessTopic       string `json:"successTopic,omitempty"`
	FailureTopic       string `json:"failureTopic,omitempty"`
}

type restoreFlags struct {
	file               string
	importMode         string
	pauseNotifications bool
	projectsImportMode string
	resource           bool
	resourceName       string
	failureTopic       string
	successTopic       string
	overwrite          bool
	apiVersion         apiVersionFlag
}

type RestoreResource struct {
	Id int `json:"id"`
}

var restoreV4BuildThreshold = 25208

func newRestoreCmd() *cobra.Command {
	f := restoreFlags{}
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restores the FME Server configuration from an import package",
		Long:  "Restores the FME Server configuration from an import package",
		PreRunE: func(cmd *cobra.Command, args []string) error {

			if f.apiVersion == "" {
				if viper.GetInt("build") < restoreV4BuildThreshold {
					f.apiVersion = apiVersionFlagV3
				} else {
					f.apiVersion = apiVersionFlagV4
				}
			}

			// ensure one of --file or --resource is set
			if f.file == "" && !f.resource {
				return errors.New("required flag \"file\" or \"resource\" not set")
			}

			// verify import mode is valid
			if f.apiVersion == apiVersionFlagV3 && f.importMode != "UPDATE" && f.importMode != "INSERT" {
				return errors.New("invalid import-mode. Must be either UPDATE or INSERT")
			}

			// verify projects import mode is valid
			if f.apiVersion == apiVersionFlagV3 && f.projectsImportMode != "UPDATE" && f.projectsImportMode != "INSERT" && f.projectsImportMode != "" {
				return errors.New("invalid projects-import-mode. Must be either UPDATE or INSERT")
			}

			// if restoring from the shared resource and no file is set, default to ServerConfigPackage.fsconfig
			if f.resource && f.file == "" {
				f.file = "ServerConfigPackage.fsconfig"
			}

			// in V3, if a failure topic or success topic is set, the restore needs to be of type "resource" as the upload endpoint doesn't support success and failure topics
			if (f.failureTopic != "" || f.successTopic != "") && !f.resource && f.apiVersion == apiVersionFlagV3 {
				return errors.New("in V3, setting a failure and/or success topic is only supported if restoring from a shared resource")
			}
			return nil
		},
		Example: `
  # Restore from a backup in a local file
  fmeflow restore --file ServerConfigPackage.fsconfig
  
  # Restore from a backup file stored in the Backup resource folder (FME_SHAREDRESOURCE_BACKUP) named ServerConfigPackage.fsconfig
  fmeflow restore --resource --file ServerConfigPackage.fsconfig
  
  # Restore from a backup file stored in the Data resource folder (FME_SHAREDRESOURCE_DATA) named ServerConfigPackage.fsconfig and set a failure and success topic to notify, overwrite items if they already exist
  fmeflow restore --resource --resource-name FME_SHAREDRESOURCE_DATA --file ServerConfigPackage.fsconfig --failure-topic MY_FAILURE_TOPIC --success-topic MY_SUCCESS_TOPIC --overwrite
  `,
		Args: NoArgs,
		RunE: restoreRun(&f),
	}

	cmd.Flags().StringVarP(&f.file, "file", "f", "", "Path to backup file to upload to restore. Can be a local file or the relative path inside the specified shared resource.")
	cmd.Flags().StringVar(&f.importMode, "import-mode", "INSERT", "To import only items in the import package that do not exist on the current instance, specify INSERT. To overwrite items on the current instance with those in the import package, specify UPDATE. Default is INSERT.")
	cmd.Flags().BoolVar(&f.pauseNotifications, "pause-notifications", true, "Disable notifications for the duration of the restore.")
	cmd.Flags().StringVar(&f.projectsImportMode, "projects-import-mode", "", "Import mode for projects. To import only projects in the import package that do not exist on the current instance, specify INSERT. To overwrite projects on the current instance with those in the import package, specify UPDATE. If not supplied, importMode will be used.")
	cmd.Flags().BoolVar(&f.resource, "resource", false, "Restore from a shared resource location instead of a local file.")
	cmd.Flags().StringVar(&f.resourceName, "resource-name", "FME_SHAREDRESOURCE_BACKUP", "Resource containing the import package. Default value is FME_SHAREDRESOURCE_BACKUP.")
	cmd.Flags().StringVar(&f.failureTopic, "failure-topic", "", "Topic to notify on failure of the import. Default is MIGRATION_ASYNC_JOB_FAILURE. Not supported when restoring from downloaded package in v3.")
	cmd.Flags().StringVar(&f.successTopic, "success-topic", "", "Topic to notify on success of the import. Default is MIGRATION_ASYNC_JOB_SUCCESS. Not supported when restoring from downloaded package in v3.")
	cmd.Flags().BoolVar(&f.overwrite, "overwrite", false, "Whether the system restore should overwrite items if they already exist.")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")

	return cmd
}
func restoreRun(f *restoreFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := &http.Client{}

		url := ""
		var request http.Request

		if f.apiVersion == apiVersionFlagV4 {
			if !f.resource {
				file, err := os.Open(f.file)
				if err != nil {
					return err
				}
				defer file.Close()

				// Create multipart form data
				var requestBody bytes.Buffer
				multiPartWriter := multipart.NewWriter(&requestBody)

				filename := filepath.Base(f.file)
				filePart, err := multiPartWriter.CreateFormFile("file", filename)
				if err != nil {
					return err
				}
				_, err = io.Copy(filePart, file)
				if err != nil {
					return err
				}

				requestData := map[string]interface{}{
					"overwrite":          f.overwrite,
					"pauseNotifications": f.pauseNotifications,
				}
				if f.successTopic != "" {
					requestData["successTopic"] = f.successTopic
				}
				if f.failureTopic != "" {
					requestData["failureTopic"] = f.failureTopic
				}

				header := make(map[string][]string)
				header["Content-Disposition"] = []string{`form-data; name="request"`}
				header["Content-Type"] = []string{"application/json"}
				requestPartWriter, err := multiPartWriter.CreatePart(header)
				if err != nil {
					return err
				}
				requestJson, err := json.Marshal(requestData)
				if err != nil {
					return err
				}
				_, err = requestPartWriter.Write(requestJson)
				if err != nil {
					return err
				}

				err = multiPartWriter.Close()
				if err != nil {
					return err
				}

				url = "/fmeapiv4/migrations/restore/upload"
				request, err = buildFmeFlowRequest(url, "POST", &requestBody)
				if err != nil {
					return err
				}
				request.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

				response, err := client.Do(&request)
				if err != nil {
					return err
				} else if response.StatusCode != 202 && response.StatusCode != 200 {
					return errors.New(response.Status)
				}
				defer response.Body.Close()

				responseData, err := io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				var result RestoreResource
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				} else {
					if !jsonOutput {
						fmt.Fprintln(cmd.OutOrStdout(), "Restore task submitted with id: "+strconv.Itoa(result.Id))
					} else {
						prettyJSON, err := prettyPrintJSON(responseData)
						if err != nil {
							return err
						}
						fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
					}
				}
			} else {
				var responseData []byte

				resourceRequest := ResourceRequestV4{}
				resourceRequest.ResourceName = f.resourceName
				if f.file != "" {
					resourceRequest.PackagePath = f.file
				}

				resourceRequest.Overwrite = f.overwrite
				resourceRequest.PauseNotifications = f.pauseNotifications

				if f.successTopic != "" {
					resourceRequest.SuccessTopic = f.successTopic
				}
				if f.failureTopic != "" {
					resourceRequest.FailureTopic = f.failureTopic
				}
				requestBodyJson, err := json.Marshal(resourceRequest)
				if err != nil {
					return err
				}

				url = "/fmeapiv4/migrations/restore/resource"
				request, err = buildFmeFlowRequest(url, "POST", bytes.NewBuffer(requestBodyJson))
				if err != nil {
					return err
				}

				request.Header.Set("Content-Type", "application/json")

				response, err := client.Do(&request)
				if err != nil {
					return err
				} else if response.StatusCode != 202 && response.StatusCode != 200 {
					return errors.New(response.Status)
				}
				defer response.Body.Close()

				responseData, err = io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				var result RestoreResource
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				} else {
					if !jsonOutput {
						fmt.Fprintln(cmd.OutOrStdout(), "Restore task submitted with id: "+strconv.Itoa(result.Id))
					} else {
						prettyJSON, err := prettyPrintJSON(responseData)
						if err != nil {
							return err
						}
						fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
					}
				}

			}

		} else if f.apiVersion == apiVersionFlagV3 {
			if !f.resource {
				file, err := os.Open(f.file)
				if err != nil {
					return err
				}
				defer file.Close()

				url = "/fmerest/v3/migration/restore/upload"
				request, err = buildFmeFlowRequest(url, "POST", file)
				if err != nil {
					return err
				}
				request.Header.Set("Content-Type", "application/octet-stream")
			} else {
				url = "/fmerest/v3/migration/restore/resource"
				var err error
				request, err = buildFmeFlowRequest(url, "POST", nil)
				if err != nil {
					return err
				}
				request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			}

			q := request.URL.Query()

			if f.resourceName != "" {
				q.Add("resourceName", f.resourceName)
			}

			if f.resource {
				q.Add("importPackage", f.file)
			}

			q.Add("pauseNotifications", strconv.FormatBool(f.pauseNotifications))

			if f.importMode != "" {
				q.Add("importMode", f.importMode)
			}

			if f.projectsImportMode != "" {
				q.Add("projectsImportMode", f.projectsImportMode)
			}

			if f.successTopic != "" {
				q.Add("successTopic", f.successTopic)
			}

			if f.failureTopic != "" {
				q.Add("failureTopic", f.failureTopic)
			}

			request.URL.RawQuery = q.Encode()

			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if !f.resource && response.StatusCode != http.StatusOK {
				if response.StatusCode == http.StatusInternalServerError {
					return fmt.Errorf("%w: check that the file specified is a valid backup file", errors.New(response.Status))
				} else {
					return errors.New(response.Status)
				}

			} else if f.resource && response.StatusCode != http.StatusAccepted {
				if response.StatusCode == http.StatusUnprocessableEntity {
					return fmt.Errorf("%w: check that the specified shared resource and file exist", errors.New(response.Status))
				}

				return errors.New(response.Status)
			}
			defer response.Body.Close()

			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			var result RestoreResource
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			} else {
				if !jsonOutput {
					fmt.Fprintln(cmd.OutOrStdout(), "Restore task submitted with id: "+strconv.Itoa(result.Id))
				} else {
					prettyJSON, err := prettyPrintJSON(responseData)
					if err != nil {
						return err
					}
					fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
				}
			}

		}
		return nil
	}
}
