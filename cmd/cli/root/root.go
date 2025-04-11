package root

import (
	"log"
	"time"

	"github.com/Dokuqui/healnix/internal/history"
	"github.com/Dokuqui/healnix/internal/monitor"
	"github.com/Dokuqui/healnix/pkg/config"
	"github.com/spf13/cobra"
)

func MonitorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "monitor",
		Short: "Start monitoring services",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadConfig("config.yaml")
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}
			hist, err := history.NewHistory("healnix.db")
			if err != nil {
				log.Fatalf("Failed to initialize history: %v", err)
			}
			defer hist.Close()
			mon := monitor.NewMonitor(cfg, hist)
			mon.Start()
		},
	}
}

func StatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current service statuses",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadConfig("config.yaml")
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}
			hist, err := history.NewHistory("healnix.db")
			if err != nil {
				log.Fatalf("Failed to initialize history: %v", err)
			}
			defer hist.Close()
			mon := monitor.NewMonitor(cfg, hist)
			mon.PrintStatus()
		},
	}
}

func HistoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "history <service>",
		Short: "Show healing history for a service",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hist, err := history.NewHistory("healnix.db")
			if err != nil {
				log.Fatalf("Failed to initialize history: %v", err)
			}
			defer hist.Close()
			hist.PrintHistory(args[0])
		},
	}
}

func PredictCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "predict <service>",
		Short: "Predict issues based on healing history",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hist, err := history.NewHistory("healnix.db")
			if err != nil {
				log.Fatalf("Failed to initialize history: %v", err)
			}
			defer hist.Close()
			suggestions := hist.Predict(args[0], 24*time.Hour)
			if len(suggestions) == 0 {
				log.Printf("No issues detected for %s", args[0])
			} else {
				log.Printf("Suggestions for %s:", args[0])
				for _, s := range suggestions {
					log.Printf("  - %s", s)
				}
			}
		},
	}
}
