package cli

import (
	"os"

	"github.com/spf13/cobra"
)

func completionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for ai-cli.

To load completions:

Bash:
  $ source <(ai-cli completion bash)
  # Or add to ~/.bashrc:
  $ ai-cli completion bash > /etc/bash_completion.d/ai-cli

Zsh:
  $ ai-cli completion zsh > "${fpath[1]}/_ai-cli"

Fish:
  $ ai-cli completion fish | source
  $ ai-cli completion fish > ~/.config/fish/completions/ai-cli.fish

PowerShell:
  PS> ai-cli completion powershell | Out-String | Invoke-Expression`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletionV2(os.Stdout, true)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}
	return cmd
}
