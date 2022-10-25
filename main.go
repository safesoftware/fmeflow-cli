package main

import (
	"fmt"
	"os"

	"github.com/safesoftare/fmeserver-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if err != cmd.ErrSilent {
			fmt.Fprintln(os.Stderr, fmt.Errorf("ERROR: %w", err))
		}
		os.Exit(1)
	}
}
