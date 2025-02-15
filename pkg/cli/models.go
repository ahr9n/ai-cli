package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/ahr9n/ollama-cli/pkg/client/ollama"
	"github.com/ahr9n/ollama-cli/pkg/config"
	"github.com/spf13/cobra"
)

func newModelsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "models",
		Short: "List available models",
		Long:  "List all available models in your Ollama installation",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			client := ollama.NewClient(cfg)
			models, err := client.ListModels()
			if err != nil {
				return fmt.Errorf("failed to list models: %w", err)
			}

			if len(models) == 0 {
				fmt.Println("No models found. Try pulling one with: ollama pull <model>")
				return nil
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "NAME\tSIZE\tMODIFIED\tFAMILY")
			fmt.Fprintln(w, "----\t----\t--------\t------")

			for _, model := range models {
				size := formatSize(model.Size)
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					model.Name,
					size,
					model.Modified,
					model.Details.Family,
				)
			}
			return w.Flush()
		},
	}

	return cmd
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
