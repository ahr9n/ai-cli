package main

import (
	"log"
	"os"

	"github.com/ahr9n/ai-cli/pkg/cli"
)

func main() {
	cmd := cli.NewRootCommand()
	if err := cmd.Execute(); err != nil {
		log.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}
