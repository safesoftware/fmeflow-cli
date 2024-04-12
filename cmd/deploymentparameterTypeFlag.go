package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type deploymentParameterTypeFlag string

const (
	deploymentParameterTypeFlagText     deploymentParameterTypeFlag = "text"
	deploymentParameterTypeFlagDatabase deploymentParameterTypeFlag = "database"
	deploymentParameterTypeFlagWeb      deploymentParameterTypeFlag = "web"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *deploymentParameterTypeFlag) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *deploymentParameterTypeFlag) Set(v string) error {
	switch v {
	case "text", "database", "web":
		*e = deploymentParameterTypeFlag(v)
		return nil
	default:
		return errors.New(`must be one of "text" or "database" or "web"`)
	}
}

// Type is only used in help text
func (e *deploymentParameterTypeFlag) Type() string {
	return "string"
}

// https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
var _ pflag.Value = (*apiVersionFlag)(nil)

// enable tab completion
func deploymentParameterTypeFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"text\ttext deployment parameter type",
		"database\tdatabase connection deployment parameter type",
		"web\tweb connection deployment parameter type",
	}, cobra.ShellCompDirectiveDefault
}
