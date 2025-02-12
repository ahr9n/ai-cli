package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "0.1.0"
)

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", Version)
		},
	}

	return cmd
}
