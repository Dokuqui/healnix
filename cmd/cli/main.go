package main

import (
	"log"

	"github.com/Dokuqui/healnix/cmd/cli/root"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "healnix"}
	rootCmd.AddCommand(root.MonitorCmd(), root.StatusCmd(), root.HistoryCmd(), root.PredictCmd())
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}
