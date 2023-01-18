package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"k8s.io/client-go/util/jsonpath"
)

var jsonRegexp = regexp.MustCompile(`^\{\.?([^{}]+)\}$|^\.?([^{}]+)$`)

// RelaxedJSONPathExpression attempts to be flexible with JSONPath expressions, it accepts:
//   - metadata.name (no leading '.' or curly braces '{...}'
//   - {metadata.name} (no leading '.')
//   - .metadata.name (no curly braces '{...}')
//   - {.metadata.name} (complete expression)
//
// And transforms them all into a valid jsonpath expression:
//
//	{.metadata.name}
func RelaxedJSONPathExpression(pathExpression string) (string, error) {
	if len(pathExpression) == 0 {
		return pathExpression, nil
	}
	submatches := jsonRegexp.FindStringSubmatch(pathExpression)
	if submatches == nil {
		return "", fmt.Errorf("unexpected path string, expected a 'name1.name2' or '.name1.name2' or '{name1.name2}' or '{.name1.name2}'")
	}
	if len(submatches) != 3 {
		return "", fmt.Errorf("unexpected submatch list: %v", submatches)
	}
	var fieldSpec string
	if len(submatches[1]) != 0 {
		fieldSpec = submatches[1]
	} else {
		fieldSpec = submatches[2]
	}
	return fmt.Sprintf("{.%s}", fieldSpec), nil
}

// This will create a table object and return it with the columns and queries specified by columnsInput
// applied to the jsonItems array
func createTableFromCustomColumns(jsonItems [][]byte, columnsInput []string) (table.Writer, error) {
	headers := table.Row{}
	t := table.NewWriter()
	t.SetStyle(defaultStyle)
	// for each row
	first := true
	for _, element := range jsonItems {
		row := table.Row{}
		// for each column
		for _, column := range columnsInput {
			if !strings.Contains(column, ":") {
				return nil, errors.New("custom column \"" + column + "\" syntax invalid")
			}
			// split on the first instance of ":"
			columnHeader, columnQuery, _ := strings.Cut(column, ":")
			columnQuery, err := RelaxedJSONPathExpression(columnQuery)
			if err != nil {
				return nil, fmt.Errorf("error parsing JSON Query for custom column: %w", err)
			}
			if first {
				headers = append(headers, columnHeader)
			}

			v := interface{}(nil)
			json.Unmarshal(element, &v)

			j := jsonpath.New("Parser")
			if err := j.Parse(columnQuery); err != nil {
				return nil, err
			}
			valueString := new(bytes.Buffer)
			err = j.Execute(valueString, v)
			if err != nil {
				return nil, fmt.Errorf("error parsing JSON Query for custom column: %w", err)
			}
			row = append(row, valueString)
		}
		first = false
		t.AppendRow(row)

	}
	t.AppendHeader(headers)
	return t, nil
}
