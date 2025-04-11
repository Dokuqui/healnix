package main

import (
	"log"

	"github.com/Dokuqui/healnix/internal/history"
	"github.com/Dokuqui/healnix/internal/monitor"
	"github.com/Dokuqui/healnix/pkg/config"
)

func main() {
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
}
