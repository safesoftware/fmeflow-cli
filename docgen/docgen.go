package main

import (
	"log"

	"github.com/safesoftare/fmeserver-cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.NewRootCommand(), "../docs")
	if err != nil {
		log.Fatal(err)
	}
}
