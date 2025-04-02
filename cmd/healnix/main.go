package main

import (
	"github.com/Dokuqui/healnix/internal/monitor"
	"github.com/Dokuqui/healnix/pkg/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	mon := monitor.NewMonitor(cfg)
	mon.Start()
}
