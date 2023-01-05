package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/jedib0t/go-pretty/v6/table"
)

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

			if first {
				headers = append(headers, columnHeader)
			}

			v := interface{}(nil)
			json.Unmarshal(element, &v)

			test, err := jsonpath.Get(columnQuery, v)
			if err != nil {
				return nil, fmt.Errorf("error parsing JSON Query for custom column: %w", err)
			}
			row = append(row, test)
		}
		first = false
		t.AppendRow(row)

	}
	t.AppendHeader(headers)
	return t, nil
}
