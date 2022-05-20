/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// oldCmd represents the old command
var oldCmd = &cobra.Command{
	Use:   "old",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//myFigure := figure.NewColorFigure("fmeserverconsole?", "", "green", true)
		//myFigure.Print()
		fmt.Println("")
		fmt.Println("What about fmeserverconsole?")
		fmt.Println("")
		fmt.Println("- Not really maintained")
		fmt.Println("- Old and clunky")
		fmt.Println("- Very limited functionality")
		fmt.Println("- Not very portable")
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(oldCmd)
	oldCmd.Hidden = true

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// oldCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// oldCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
