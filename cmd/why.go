/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

// whyCmd represents the why command
var whyCmd = &cobra.Command{
	Use:   "why",
	Short: "Outputs Why",
	Long:  `Outputs Why`,
	Run: func(cmd *cobra.Command, args []string) {
		myFigure := figure.NewColorFigure("WHY?", "", "green", true)
		myFigure.Print()
		fmt.Println("")
		fmt.Println("- A good CLI is helpful for writing pipelines")
		fmt.Println("- REST API is powerful, but writing REST calls is annoying")
		fmt.Println("- Easier to run things locally to test before running in pipeline")
		fmt.Println("- Can implement more complicated things that make multiple different API calls")
	},
}

func init() {
	rootCmd.AddCommand(whyCmd)
	whyCmd.Hidden = true

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
