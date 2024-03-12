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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type projectUploadFlags struct {
	file                string
	overwrite           bool
	importMode          string
	pauseNotifications  bool
	projectsImportMode  string
	disableProjectItems bool
	getSelectable       bool
	selectedItems       string
	interactive         bool
	quick               bool
	wait                bool
	backupFailureTopic  string
	backupSuccessTopic  string
	apiVersion          apiVersionFlag
}

type ProjectUploadTask struct {
	Id int `json:"id"`
}

type ProjectItems struct {
	Items      []ProjectItemV4 `json:"items"`
	TotalCount int             `json:"totalCount"`
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
}

type ProjectItemV4 struct {
	ID            string `json:"id"`
	JobID         int    `json:"jobId"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	OwnerID       string `json:"ownerId"`
	OwnerName     string `json:"ownerName"`
	OwnerStatus   string `json:"ownerStatus"`
	OriginalOwner string `json:"originalOwner"`
	Selected      bool   `json:"selected"`
	Existing      bool   `json:"existing"`
	PreviewAction string `json:"previewAction"`
	Action        string `json:"action"`
	Source        string `json:"source"`
}

type ProjectImportRun struct {
	FallbackOwnerID    string                 `json:"fallbackOwnerID,omitempty"`
	Overwrite          bool                   `json:"overwrite"`
	PauseNotifications bool                   `json:"pauseNotifications"`
	DisableItems       bool                   `json:"disableItems"`
	Notification       *ProjectNotification   `json:"notification,omitempty"`
	SelectedItems      []ProjectSelectedItems `json:"selectedItems"`
}

type ProjectNotification struct {
	Type         string `json:"type,omitempty"`
	SuccessTopic string `json:"successTopic,omitempty"`
	FailureTopic string `json:"failureTopic,omitempty"`
}

type ProjectSelectedItems struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type ProjectUploadV4 struct {
	JobID     int              `json:"jobId"`
	Status    string           `json:"status"`
	Owner     string           `json:"owner"`
	OwnerID   string           `json:"ownerID"`
	Requested time.Time        `json:"requested"`
	Generated time.Time        `json:"generated"`
	FileName  string           `json:"fileName"`
	Request   ProjectImportRun `json:"request"`
}

var projectUploadV4BuildThreshold = 23766

func newProjectUploadCmd() *cobra.Command {
	f := projectUploadFlags{}
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Imports FME Server Projects from a downloaded package.",
		Long: `Imports FME Server Projects from a downloaded package. The upload happens in two steps. The package is uploaded to the server, a preview is generated that contains the list of items, and then the import is run. This command can be run using a few different modes.
- Using the --get-selectable flag will just generate the preview and output the selectable items in the package and then delete the import
- Using the --quick flag will skip the preview and import everything in the package by default.
- Using the --interactive flag will prompt the user to select items to import from the list of selectable items if any exist
- Using the --selected-items flag will import only the items specified. The default is to import all items in the package.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// get build to decide if we should use v3 or v4
			// FME Server 2022.0 and later can use v4. Otherwise fall back to v3
			if f.apiVersion == "" {
				fmeflowBuild := viper.GetInt("build")
				if fmeflowBuild < projectUploadV4BuildThreshold {
					f.apiVersion = apiVersionFlagV3
				} else {
					f.apiVersion = apiVersionFlagV4
				}
			}

			// verify import mode is valid
			if f.importMode != "UPDATE" && f.importMode != "INSERT" && f.importMode != "" {
				return errors.New("invalid import-mode. Must be either UPDATE or INSERT")
			}

			// verify projects import mode is valid
			if f.projectsImportMode != "UPDATE" && f.projectsImportMode != "INSERT" && f.projectsImportMode != "" {
				return errors.New("invalid projects-import-mode. Must be either UPDATE or INSERT")
			}

			if f.apiVersion == apiVersionFlagV4 {
				if f.importMode != "" {
					return errors.New("cannot set the importMode flag when using the V4 API")
				}
				if f.projectsImportMode != "" {
					return errors.New("cannot set the projectsImportMode flag when using the V4 API")
				}
			}

			if f.apiVersion == apiVersionFlagV3 {
				if f.importMode == "" && f.projectsImportMode == "" {
					if f.overwrite {
						f.importMode = "UPDATE"
					} else {
						f.importMode = "INSERT"
					}
				}
			}

			return nil
		},
		Example: `
  # Upload a project and import all selectable items if any exist
  fmeflow projects upload --file ProjectPackage.fsproject

  # Upload a project without overwriting existing items
  fmeflow projects upload --file ProjectPackage.fsproject --overwrite=false
  
  # Upload a project and perform a quick import
  fmeflow projects upload --file ProjectPackage.fsproject --quick
  
  # Upload a project and be prompted for which items to import of the selectable items
  fmeflow projects upload --file ProjectPackage.fsproject --interactive 
 
  # Upload a project and get the list of selectable items
  fmeflow projects upload --file ProjectPackage.fsproject --get-selectable
  
  # Upload a project and import only the specified selectable items
  fme projects upload --file ProjectPackage.fsproject --selected-items="mysqldb:connection,slack con:connector"`,
		Args: NoArgs,
		RunE: projectUploadRun(&f),
	}

	cmd.Flags().StringVarP(&f.file, "file", "f", "", "Path to backup file to upload to restore. Can be a local file or the relative path inside the specified shared resource.")
	cmd.Flags().BoolVar(&f.overwrite, "overwrite", true, "If specified, the items in the project will overwrite existing items.")
	cmd.Flags().BoolVar(&f.pauseNotifications, "pause-notifications", true, "Disable notifications for the duration of the restore.")
	cmd.Flags().BoolVar(&f.disableProjectItems, "disable-project-items", false, "Whether to disable items in the imported FME Server Projects. If true, items that are new or overwritten will be imported but disabled. If false, project items are imported as defined in the import package.")
	cmd.Flags().BoolVar(&f.getSelectable, "get-selectable", false, "Output the selectable items in the import package.")
	cmd.Flags().StringVar(&f.selectedItems, "selected-items", "all", "The items to import. Set to \"all\" to import all items, and \"none\" to omit selectable items. Otherwise, this should be a comma separated list of item ids type pairs separated by a colon. e.g. a:b,c:d")
	cmd.Flags().BoolVar(&f.interactive, "interactive", false, "Prompt interactively for the selectable items to import (if any exist).")
	cmd.Flags().BoolVar(&f.quick, "quick", false, "Import everything in the package by default.")
	cmd.Flags().BoolVar(&f.wait, "wait", true, "Wait for import to complete. Set to false to return immediately after the import is started.")
	cmd.Flags().StringVar(&f.backupFailureTopic, "failure-topic", "MIGRATION_ASYNC_JOB_FAILURE", "Topic to notify on failure of the backup.")
	cmd.Flags().StringVar(&f.backupSuccessTopic, "success-topic", "MIGRATION_ASYNC_JOB_SUCCESS", "Topic to notify on success of the backup.")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	// these flags are only for v3
	cmd.Flags().StringVar(&f.importMode, "import-mode", "", "To import only items in the import package that do not exist on the current instance, specify INSERT. To overwrite items on the current instance with those in the import package, specify UPDATE. Default is INSERT.")
	cmd.Flags().StringVar(&f.projectsImportMode, "projects-import-mode", "", "Import mode for projects. To import only projects in the import package that do not exist on the current instance, specify INSERT. To overwrite projects on the current instance with those in the import package, specify UPDATE. If not supplied, importMode will be used.")

	cmd.Flags().MarkHidden("api-version")
	cmd.Flags().MarkHidden("projects-import-mode")
	cmd.Flags().MarkHidden("import-mode")
	cmd.MarkFlagsMutuallyExclusive("overwrite", "projects-import-mode")
	cmd.MarkFlagsMutuallyExclusive("overwrite", "import-mode")
	cmd.MarkFlagsMutuallyExclusive("quick", "get-selectable")
	cmd.MarkFlagsMutuallyExclusive("interactive", "get-selectable")
	cmd.MarkFlagsMutuallyExclusive("get-selectable", "overwrite")
	cmd.MarkFlagsMutuallyExclusive("get-selectable", "pause-notifications")
	cmd.MarkFlagsMutuallyExclusive("selected-items", "interactive")
	cmd.MarkFlagsMutuallyExclusive("selected-items", "get-selectable")
	cmd.MarkFlagsMutuallyExclusive("selected-items", "quick")
	cmd.MarkFlagsMutuallyExclusive("quick", "interactive")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
	cmd.MarkFlagRequired("file")

	return cmd
}

func projectUploadRun(f *projectUploadFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := &http.Client{}

		url := ""
		var request http.Request
		file, err := os.Open(f.file)
		if err != nil {
			return err
		}
		defer file.Close()

		if f.apiVersion == "v4" {
			// Create a buffer to store our request body as bytes
			var requestBody bytes.Buffer

			// Create a multipart writer
			multiPartWriter := multipart.NewWriter(&requestBody)

			// Create a form file writer for the file field
			fileWriter, err := multiPartWriter.CreateFormFile("file", f.file)
			if err != nil {
				return err
			}

			// Copy the file data to the form file writer
			if _, err = io.Copy(fileWriter, file); err != nil {
				return err
			}

			// Close the multipart writer to get the terminating boundary
			if err = multiPartWriter.Close(); err != nil {
				return err
			}

			url = "/fmeapiv4/migrations/imports/upload"
			request, err = buildFmeFlowRequest(url, "POST", &requestBody)
			if err != nil {
				return err
			}
			// body as multipart form
			request.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

			if f.quick {
				// add the skip preview query parameter for quick import
				q := request.URL.Query()
				q.Add("skipPreview", strconv.FormatBool(f.quick))
				request.URL.RawQuery = q.Encode()
			}

			// execute the upload of the package
			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != http.StatusCreated {
				if response.StatusCode == http.StatusInternalServerError {
					return fmt.Errorf("%w: check that the file specified is a valid project file", errors.New(response.Status))
				} else {
					return errors.New(response.Status)
				}
			}

			// read the Location header to get the task id
			location := response.Header.Get("Location")
			// parse the task id from the location header
			// the task id is an integer at the end of the location header
			taskId := location[strings.LastIndex(location, "/")+1:]

			var selectedItemsStruct []ProjectSelectedItems
			// if this isn't a quick import, we need to get the selectable items by making another rest call.
			// also, if it isn't a quick import, it takes a bit of time for the preview to be ready, so we have to wait for it
			if !f.quick {

				// We have to do a get on the import to see if the status is ready
				ready := false
				tries := 0
				url = "/fmeapiv4/migrations/imports/" + taskId
				request, err = buildFmeFlowRequest(url, "GET", nil)
				if err != nil {
					return err
				}

				if !jsonOutput && !f.getSelectable {
					fmt.Fprint(cmd.OutOrStdout(), "Waiting for preview generation..")
				}

				// we have to loop until the preview is done generating
				for !ready {
					// get the status of the import
					response, err = client.Do(&request)
					if err != nil {
						return err
					} else if response.StatusCode != http.StatusOK {
						return errors.New(response.Status)
					}

					responseData, err := io.ReadAll(response.Body)
					if err != nil {
						return err
					}

					var importStatus ProjectUploadV4
					if err := json.Unmarshal(responseData, &importStatus); err != nil {
						return err
					}
					// check if it is ready
					if importStatus.Status == "ready" {
						ready = true
					} else if importStatus.Status == "generating_preview" {
						// if it is still generating the preview, wait a second and try again
						if !jsonOutput && !f.getSelectable {
							fmt.Fprint(cmd.OutOrStdout(), ".")
						}
						time.Sleep(1 * time.Second)
						tries++
					} else {
						return errors.New("import task did not complete successfully. Status is \"" + importStatus.Status + "\". Please check the FME Flow web interface for the status of the import task")
					}
				}

				// output a newline to cap off the waiting message
				if !jsonOutput && !f.getSelectable {
					fmt.Fprint(cmd.OutOrStdout(), "\n")
				}

				// get the selectable items from the preview
				url = "/fmeapiv4/migrations/imports/" + taskId + "/items"
				// set up the URL to query
				request, err := buildFmeFlowRequest(url, "GET", nil)
				if err != nil {
					return err
				}

				q := request.URL.Query()
				q.Add("selectable", "true")

				request.URL.RawQuery = q.Encode()

				response, err = client.Do(&request)
				if err != nil {
					return err
				} else if response.StatusCode != http.StatusOK {
					if response.StatusCode == http.StatusInternalServerError {
						return fmt.Errorf("%w: check that the file specified is a valid project file", errors.New(response.Status))
					} else {
						return errors.New(response.Status)
					}
				}

				responseData, err := io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				// store the selectable items in a struct
				var selectableItems ProjectItems
				if err := json.Unmarshal(responseData, &selectableItems); err != nil {
					return err
				}

				// if we are just outputing the selectable items for this package, just output them, delete the import and return
				if f.getSelectable {
					// delete the import since we are just getting the selectable items
					url = "/fmeapiv4/migrations/imports/" + taskId
					request, err = buildFmeFlowRequest(url, "DELETE", nil)
					if err != nil {
						return err
					}
					response, err = client.Do(&request)
					if err != nil {
						return err
					} else if response.StatusCode != http.StatusNoContent {
						fmt.Fprintln(cmd.OutOrStdout(), "Failed to delete the import task with id "+taskId+". You may need to delete it manually.")
					}

					// output the selectable items
					if jsonOutput {
						prettyJSON, err := prettyPrintJSON(responseData)
						if err != nil {
							return err
						}
						fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
					} else {
						t := table.NewWriter()
						t.SetStyle(defaultStyle)

						t.AppendHeader(table.Row{"Id", "Type"})

						for _, element := range selectableItems.Items {
							t.AppendRow(table.Row{element.ID, element.Type})
						}
						//if f.noHeaders {
						//	t.ResetHeaders()
						//}
						fmt.Fprintln(cmd.OutOrStdout(), t.Render())
					}
					return nil
				}

				// if we are interactive, we want to prompt the user to select items from the list of selectable ones
				if f.interactive {
					// store the item ids and types in a string array so that we can prompt the user to select items
					var items []string
					for _, element := range selectableItems.Items {
						items = append(items, element.ID+" ("+element.Type+")")
					}

					// prompt the user to select items
					// if items is empty, the prompt will automatically be skipped
					var selectedItems []string
					prompt := &survey.MultiSelect{
						Message: "Select items to import",
						Options: items,
					}
					survey.AskOne(prompt, &selectedItems)

					for _, element := range selectedItems {
						// split the string on the space to get the id and type
						split := strings.Split(element, " ")
						id := split[0]
						itemType := strings.Trim(split[1], "()")
						selectedItemsStruct = append(selectedItemsStruct, ProjectSelectedItems{ID: id, Type: itemType})
					}
				} else {
					// if we are not interactive, check the selected items flag to see what to import
					if f.selectedItems == "all" {
						// loop through all items and format them like we like our input so we can process them in the same way
						selectedItemsString := ""
						for _, element := range selectableItems.Items {
							selectedItemsString += element.ID + ":" + element.Type + ","
						}
						selectedItemsString = strings.TrimSuffix(selectedItemsString, ",")
						f.selectedItems = selectedItemsString
					} else if f.selectedItems == "none" {
						// if we don't want to select anything, set selectedItemsStruct to an empty list
						selectedItemsStruct = []ProjectSelectedItems{}
					}

					// parse the selected items and add them to the selectedItemsStruct
					if f.selectedItems != "" && f.selectedItems != "none" {
						// validate that the selected items are the correct format
						match, _ := regexp.MatchString(`^([^:,]+:[^:,]+,)*([^:,]+:[^:,]+)$`, f.selectedItems)
						if !match {
							return errors.New("invalid selected items. Must be a comma separated list of item ids and types. e.g. item1:itemtype1,item2:itemtype2")
						}

						// split the selected items on the comma
						split := strings.Split(f.selectedItems, ",")
						for _, element := range split {
							// split the string on the colon to get the id and type
							split := strings.Split(element, ":")
							id := split[0]
							itemType := split[1]
							// verify that the selected items are in the list of selectable items
							found := false
							for _, selectableItem := range selectableItems.Items {
								if selectableItem.ID == id && selectableItem.Type == itemType {
									found = true
									break
								}
							}
							if !found {
								return errors.New("selected item " + id + " (" + itemType + ") is not in the list of selectable items")
							}

							selectedItemsStruct = append(selectedItemsStruct, ProjectSelectedItems{ID: id, Type: itemType})
						}
					}
				}

			}

			if !f.getSelectable {
				// finally, we can run the import
				url = "/fmeapiv4/migrations/imports/" + taskId + "/run"
				var run ProjectImportRun
				// set the run struct
				run.Overwrite = f.overwrite
				run.PauseNotifications = f.pauseNotifications
				run.DisableItems = f.disableProjectItems

				// if we have selected items, set them here. If we haven't set it, use the default
				if selectedItemsStruct != nil {
					run.SelectedItems = selectedItemsStruct
				}

				// if a topic is specified, add that
				if f.backupFailureTopic != "" || f.backupSuccessTopic != "" {
					run.Notification = new(ProjectNotification)
					run.Notification.Type = "TOPIC"
					if f.backupSuccessTopic != "" {
						run.Notification.SuccessTopic = f.backupSuccessTopic
					}
					if f.backupFailureTopic != "" {
						run.Notification.FailureTopic = f.backupFailureTopic
					}
				} else {
					run.Notification = nil
				}

				// marshal the run struct to json
				runJson, err := json.Marshal(run)
				if err != nil {
					return err
				}

				request, err = buildFmeFlowRequest(url, "POST", bytes.NewReader(runJson))
				if err != nil {
					return err
				}
				request.Header.Set("Content-Type", "application/json")

				response, err = client.Do(&request)
				if err != nil {
					return err
				}
				if response.StatusCode != http.StatusAccepted {
					return errors.New(response.Status)
				} else {
					if !jsonOutput {
						fmt.Fprintln(cmd.OutOrStdout(), "Project Upload task submitted with id: "+taskId)
					} else if !f.wait {
						// if we are outputting json and not waiting, do a get on the task and output that
						url = "/fmeapiv4/migrations/imports/" + taskId
						request, err = buildFmeFlowRequest(url, "GET", nil)
						if err != nil {
							return err
						}

						response, err = client.Do(&request)
						if err != nil {
							return err
						} else if response.StatusCode != http.StatusOK {
							return errors.New(response.Status)
						}

						responseData, err := io.ReadAll(response.Body)
						if err != nil {
							return err
						}

						prettyJSON, err := prettyPrintJSON(responseData)
						if err != nil {
							return err
						}
						fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
					}
				}

				// if we are waiting for the import to complete, we have to loop until it is done
				if f.wait {
					finished := false
					url = "/fmeapiv4/migrations/imports/" + taskId
					request, err = buildFmeFlowRequest(url, "GET", nil)
					if err != nil {
						return err
					}

					if !jsonOutput {
						fmt.Fprint(cmd.OutOrStdout(), "Waiting for project to finish importing..")
					}
					var importStatus ProjectUploadV4
					for !finished {

						response, err = client.Do(&request)
						if err != nil {
							return err
						} else if response.StatusCode != http.StatusOK {
							return errors.New(response.Status)
						}

						responseData, err := io.ReadAll(response.Body)
						if err != nil {
							return err
						}

						if err := json.Unmarshal(responseData, &importStatus); err != nil {
							return err
						}

						if importStatus.Status == "imported" {
							finished = true
						} else if importStatus.Status != "importing" {
							return errors.New("import task did not complete successfully. Please check the FME Flow web interface for the status of the import task")
						} else {
							if !jsonOutput {
								fmt.Fprint(cmd.OutOrStdout(), ".")
							}

							time.Sleep(1 * time.Second)
						}
					}
					if !jsonOutput {
						fmt.Fprint(cmd.OutOrStdout(), "\n")
						fmt.Fprintln(cmd.OutOrStdout(), "Project import complete.")
					} else {
						jsonData, err := json.Marshal(importStatus)
						if err != nil {
							return err
						}
						prettyJSON, err := prettyPrintJSON(jsonData)
						if err != nil {
							return err
						}
						fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
					}
				}

			}

		} else if f.apiVersion == "v3" {

			url = "/fmerest/v3/projects/import/upload"
			request, err = buildFmeFlowRequest(url, "POST", file)
			if err != nil {
				return err
			}
			request.Header.Set("Content-Type", "application/octet-stream")

			q := request.URL.Query()

			if f.pauseNotifications {
				q.Add("pauseNotifications", strconv.FormatBool(f.pauseNotifications))
			}

			if f.importMode != "" {
				q.Add("importMode", f.importMode)
			}

			if f.projectsImportMode != "" {
				q.Add("projectsImportMode", f.projectsImportMode)
			}

			if f.disableProjectItems {
				q.Add("disableProjectItems", strconv.FormatBool(f.disableProjectItems))
			}

			request.URL.RawQuery = q.Encode()

			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != http.StatusOK {
				if response.StatusCode == http.StatusInternalServerError {
					return fmt.Errorf("%w: check that the file specified is a valid project file", errors.New(response.Status))
				} else {
					return errors.New(response.Status)
				}
			}

			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			var result ProjectUploadTask
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			} else {
				if !jsonOutput {
					fmt.Fprintln(cmd.OutOrStdout(), "Project Upload task submitted with id: "+strconv.Itoa(result.Id))
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
