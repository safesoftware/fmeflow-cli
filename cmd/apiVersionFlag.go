package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type apiVersionFlag string

const (
	apiVersionFlagV3 apiVersionFlag = "v3"
	apiVersionFlagV4 apiVersionFlag = "v4"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *apiVersionFlag) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *apiVersionFlag) Set(v string) error {
	switch v {
	case "v3", "v4":
		*e = apiVersionFlag(v)
		return nil
	default:
		return errors.New(`must be one of "v3" or "v4"`)
	}
}

// Type is only used in help text
func (e *apiVersionFlag) Type() string {
	return "string"
}

// https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
var _ pflag.Value = (*apiVersionFlag)(nil)

// enable tab completion
func apiVersionFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"v3\tuse the v3 API",
		"v4\tuse the v4 API",
	}, cobra.ShellCompDirectiveDefault
}
