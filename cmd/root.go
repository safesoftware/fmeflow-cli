package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var jsonOutput bool
var outputType string
var noHeaders bool

const notSet string = "not set"

// this information will be collected at build time, by `-ldflags "-X github.com/safesoftare/fmeserver-cli/cmd.appVersion=0.1"`
var appVersion = notSet

var ErrSilent = errors.New("ErrSilent")

// rootCmd represents the base command when called without any subcommands
var rootCmd = NewRootCommand()

func NewRootCommand() *cobra.Command {
	cmds := &cobra.Command{
		Use:           "fmeserver",
		Short:         "A command line interface for interacting with FME Server.",
		Long:          `A command line interface for interacting with FME Server.`,
		Version:       appVersion,
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// make sure the config file is set up correctly
			_, err := os.Stat(viper.ConfigFileUsed())
			if err != nil {
				return fmt.Errorf("could not open the config file " + viper.ConfigFileUsed() + ". Have you called the login command? ")
			}
			fmeserverUrl := viper.GetString("url")

			// check the fme server URL is valid
			_, err = url.ParseRequestURI(fmeserverUrl)
			if err != nil {
				return fmt.Errorf("invalid FME Server url in config file " + viper.ConfigFileUsed() + ". Have you called the login command? ")
			}

			// check there is a token to use for auth
			fmeserverToken := viper.GetString("token")
			if fmeserverToken == "" {
				return fmt.Errorf("no token found in config file " + viper.ConfigFileUsed() + ". Have you called the login command? ")
			}

			// check there is a build set in the config file
			fmeserverBuild := viper.GetString("build")
			if fmeserverBuild == "" {
				return fmt.Errorf("no build found in config file " + viper.ConfigFileUsed() + ". Have you called the login command? ")
			}

			return nil
		},
	}
	cmds.ResetFlags()
	cmds.AddCommand(newHealthcheckCmd())
	cmds.AddCommand(newBackupCmd())
	cmds.AddCommand(newEnginesCmd())
	cmds.AddCommand(newJobsCmd())
	cmds.AddCommand(newInfoCmd())
	cmds.AddCommand(newLicenseCmd())
	cmds.AddCommand(newLoginCmd())
	cmds.AddCommand(newMigrationCmd())
	cmds.AddCommand(newRestoreCmd())
	cmds.AddCommand(newRunCmd())
	cmds.AddCommand(newCancelCmd())
	cmds.AddCommand(newRepositoryCmd())
	cmds.AddCommand(newWorkspaceCmd())
	cmds.AddCommand(newProjectsCmd())
	cmds.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.PrintErrln(err)
		cmd.PrintErrln(cmd.UsageString())
		return ErrSilent
	})
	cobra.OnInitialize(initConfig)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	cmds.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/.fmeserver-cli.yaml)")
	cmds.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output JSON")

	return cmds
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	err := rootCmd.Execute()
	return err
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else if configFilePath := os.Getenv("FMESERVER_CLI_CONFIG"); configFilePath != "" {
		// use path from FMESERVER_CLI_CONFIG if set
		viper.SetConfigFile(configFilePath)
	} else {
		// check if XDG_CONFIG_HOME is set
		defaultConfigDirectory := os.Getenv("XDG_CONFIG_HOME")
		if defaultConfigDirectory == "" {
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			defaultConfigDirectory = filepath.Join(home, ".config")
		}

		viper.SetConfigFile(filepath.Join(defaultConfigDirectory, ".fmeserver-cli.yaml"))

	}
	//fmt.Println(viper.ConfigFileUsed())
	//viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()

}

// Function for commands that provide no arguments. This will turn usage back on
// so that it will be output if a user tries to pass in an argument when they should not
func NoArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		cmd.Usage()
		return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
	}
	return nil
}
