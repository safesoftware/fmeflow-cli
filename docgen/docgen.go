package main

import (
	"log"

	"github.com/safesoftware/fmeflow-cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmdRoot := cmd.NewRootCommand()
	cmdRoot.InitDefaultCompletionCmd()
	err := doc.GenMarkdownTree(cmdRoot, "../docs")
	if err != nil {
		log.Fatal(err)
	}
}
