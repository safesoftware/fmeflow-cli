/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

// whatCmd represents the what command
var whatCmd = &cobra.Command{
	Use:   "what",
	Short: "Outputs what this thing is",
	Long:  `Outputs what this thing is`,
	Run: func(cmd *cobra.Command, args []string) {
		myFigure := figure.NewColorFigure("What is this?", "", "green", true)
		myFigure.Print()
		fmt.Println("")
		fmt.Println("- A Command Line Interface for FME Server")
		fmt.Println("- Makes calls to FME Server REST API")
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(whatCmd)
	whatCmd.Hidden = true

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
