package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/securiter/linear-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if strings.Contains(err.Error(), "LINEAR_API_KEY") {
			fmt.Fprintln(os.Stderr, "\nSet LINEAR_API_KEY env var or use: psst --global LINEAR_API_KEY -- linear ...")
		}
		os.Exit(1)
	}
}
