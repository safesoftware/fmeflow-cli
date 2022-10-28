package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var jsonOutput bool
var outputType string
var noHeaders bool

var ErrSilent = errors.New("ErrSilent")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "fmeserver",
	Short:         "A command line interface for interacting with FME Server.",
	Long:          `A command line interface for interacting with FME Server.`,
	Version:       "0.4",
	SilenceErrors: true,
	SilenceUsage:  true,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	err := rootCmd.Execute()
	return err
}

func init() {
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return ErrSilent
	})
	cobra.OnInitialize(initConfig)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fmeserver-cli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output JSON")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".fmeserver-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".fmeserver-cli")
		viper.SetConfigFile(path.Join(home, ".fmeserver-cli.yaml"))

	}
	//fmt.Println(viper.ConfigFileUsed())
	//viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}

// Function for commands that provide no arguments. This will turn usage back on
// so that it will be output if a user tries to pass in an argument when they should not
func NoArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		rootCmd.SilenceUsage = false
		return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
	}
	return nil
}
