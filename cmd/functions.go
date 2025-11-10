package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"unicode"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/viper"
)

// define the default style for all tables that are output
var defaultStyle = table.Style{
	Name: "Borderless",
	Box:  table.StyleBoxLight,
	Format: table.FormatOptions{
		Footer: text.FormatUpper,
		Header: text.FormatUpper,
		Row:    text.FormatDefault,
	},
	Options: table.Options{
		DrawBorder:      false,
		SeparateColumns: false,
		SeparateFooter:  false,
		SeparateHeader:  false,
		SeparateRows:    false,
	},
}

func buildFmeFlowRequest(endpoint string, method string, body io.Reader) (http.Request, error) {
	// retrieve url and token
	fmeflowUrl := viper.GetString("url")
	fmeflowToken := viper.GetString("token")

	req, err := http.NewRequest(method, fmeflowUrl+endpoint, body)
	if fmeflowToken != "" {
		req.Header.Set("Authorization", "fmetoken token="+fmeflowToken)
	}
	return *req, err
}

// since the JSON for published parameters has subtypes, we need to implement this ourselves
func (f *JobRequestV3) UnmarshalJSON(b []byte) error {
	type job JobRequestV3
	err := json.Unmarshal(b, (*job)(f))
	if err != nil {
		return err
	}

	for _, raw := range f.RawPublishedParameters {
		data := make(map[string]json.RawMessage)
		err = json.Unmarshal(raw, &data)
		if err != nil {
			return err
		}

		var i interface{}
		for k, v := range data {
			if k == "value" {
				if strings.HasPrefix(string(v), "[") {
					i = &ListParameter{}
				} else {
					i = &SimpleParameter{}
				}

			}
		}

		if i != nil {
			err = json.Unmarshal(raw, i)
			if err != nil {
				return err
			}
			f.PublishedParameters = append(f.PublishedParameters, i)
		}
	}
	return nil
}

func (f *JobRequestV3) MarshalJSON() ([]byte, error) {

	type job JobRequestV3
	if f.PublishedParameters != nil {
		for _, v := range f.PublishedParameters {
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			f.RawPublishedParameters = append(f.RawPublishedParameters, b)
		}
	}
	return json.Marshal((*job)(f))
}

func prettyPrintJSON(s []byte) (string, error) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, s, "", "  ")
	if err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func convertCamelCaseToTitleCase(s string) string {
	result := ""
	for i, c := range s {
		if unicode.IsUpper(c) && i != 0 {
			result += " " + string(c)
		} else {
			if i == 0 {
				result += strings.ToUpper(string(c))
			} else {
				result += string(c)
			}
		}
	}
	return result
}

// Pass in a struct that represents a JSON result and return a single row table
// with column headers set to the JSON attribute name
func createTableWithDefaultColumns(s any) table.Writer {

	v := reflect.ValueOf(s)
	typeOfS := v.Type()
	header := table.Row{}
	row := table.Row{}
	for i := 0; i < v.NumField(); i++ {
		header = append(header, convertCamelCaseToTitleCase(typeOfS.Field(i).Name))
		row = append(row, v.Field(i).Interface())
	}

	t := table.NewWriter()
	t.SetStyle(defaultStyle)

	t.AppendHeader(header)
	t.AppendRow(row)

	return t
}

// helper function to check if a struct is empty
func isEmpty(object interface{}) bool {
	//First check normal definitions of empty
	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	//Then see if it's a struct
	if reflect.ValueOf(object).Kind() == reflect.Struct {
		// and create an empty copy of the struct object to compare against
		empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()
		if reflect.DeepEqual(object, empty) {
			return true
		}
	}
	return false
}

func checkConfigFile(requireToken bool) error {
	// make sure the config file is set up correctly
	_, err := os.Stat(viper.ConfigFileUsed())
	if err != nil {
		return fmt.Errorf("could not open the config file " + viper.ConfigFileUsed() + ". Have you called the login command? ")
	}
	fmeflowUrl := viper.GetString("url")

	// check the fme server URL is valid
	_, err = url.ParseRequestURI(fmeflowUrl)
	if err != nil {
		return fmt.Errorf("invalid FME Flow url in config file " + viper.ConfigFileUsed() + ". Have you called the login command? ")
	}

	if requireToken {
		// check there is a token to use for auth
		fmeflowToken := viper.GetString("token")
		if fmeflowToken == "" {
			return fmt.Errorf("no token found in config file " + viper.ConfigFileUsed() + ". Have you called the login command? ")
		}
	}

	// check there is a build set in the config file
	fmeflowBuild := viper.GetString("build")
	if fmeflowBuild == "" {
		return fmt.Errorf("no build found in config file " + viper.ConfigFileUsed() + ". Have you called the login command? ")
	}
	return nil

}
