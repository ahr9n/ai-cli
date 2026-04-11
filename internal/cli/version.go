package cli

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Version is set via ldflags at build time (GoReleaser, Makefile).
var Version = "dev"

func resolveVersion() string {
	if Version != "dev" {
		return Version
	}
	// Fallback for `go install`: read the module version embedded by the Go toolchain.
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}
	return "dev"
}

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ai-cli version %s\n", resolveVersion())
		},
	}
}
