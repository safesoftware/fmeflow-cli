/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

// futureCmd represents the future command
var futureCmd = &cobra.Command{
	Use:   "future",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		myFigure := figure.NewColorFigure("Future...", "", "green", true)
		myFigure.Print()
		fmt.Println("")
		fmt.Println("- A good CLI is valuable for FME Server in my opinon")
		fmt.Println("- This will be another thing we need to maintain")
		fmt.Println("- Would need to think deeply about what we implement and how the CLI would work so we are consistent")
		fmt.Println("- Would need to think about versioning as the API changes and how we handle that in this CLI")
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(futureCmd)
	futureCmd.Hidden = true
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// futureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// futureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
