/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

// howCmd represents the how command
var howCmd = &cobra.Command{
	Use:   "how",
	Short: "Outputs how",
	Long:  `Outputs how`,
	Run: func(cmd *cobra.Command, args []string) {
		myFigure := figure.NewColorFigure("How?", "", "green", true)
		myFigure.Print()
		fmt.Println("")
		fmt.Println("- Written in Go")
		fmt.Println("- Uses a Go package called Cobra to help generate the CLI")
	},
}

func init() {
	rootCmd.AddCommand(howCmd)
	howCmd.Hidden = true

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// howCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// howCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
