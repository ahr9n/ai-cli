package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "0.1.0"
)

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ai-cli version %s\n", Version)
		},
	}
}
