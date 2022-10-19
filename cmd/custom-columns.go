package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/jedib0t/go-pretty/v6/table"
)

var jsonRegexp = regexp.MustCompile(`^\{\.?([^{}]+)\}$|^\.?([^{}]+)$`)

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
			headerQueryArr := strings.Split(column, ":")
			columnHeader := headerQueryArr[0]
			columnQuery, err := massageQuery(headerQueryArr[1])
			if err != nil {
				return nil, err
			}
			if first {
				headers = append(headers, columnHeader)
			}

			v := interface{}(nil)
			json.Unmarshal(element, &v)

			test, err := jsonpath.Get(columnQuery, v)
			if err != nil {
				return nil, err
			}
			row = append(row, test)
		}
		first = false
		t.AppendRow(row)

	}
	t.AppendHeader(headers)
	return t, nil
}
func massageQuery(q string) (string, error) {
	submatches := jsonRegexp.FindStringSubmatch(q)
	if submatches == nil {
		return "", errors.New("unexpected path string, expected a 'name1.name2' or '.name1.name2' or '{name1.name2}' or '{.name1.name2}'")
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
	return fieldSpec, nil
}
