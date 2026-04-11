package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahr9n/ai-cli/internal/cli"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := cli.NewRootCommand().ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
