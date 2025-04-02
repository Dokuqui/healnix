package monitor

import (
	"github.com/Dokuqui/healnix/pkg/config"
	"log"
	"net/http"
	"time"
)

type Monitor struct {
	Services []config.Service
}

func NewMonitor(cfg *config.Config) *Monitor {
	return &Monitor{Services: cfg.Services}
}

func (m *Monitor) Start() {
	for _, svc := range m.Services {
		go m.checkService(svc)
	}
	select {}
}

func (m *Monitor) checkService(svc config.Service) {
	for {
		start := time.Now()
		resp, err := http.Get(svc.Endpoint)
		latency := time.Since(start).Milliseconds()

		if err != nil {
			log.Printf("Service %s failed: %v", svc.Name, err)
		} else if latency > int64(svc.Threshold) {
			log.Printf("Service %s unhealthy (latency: %dms > %dms)", svc.Name, latency, svc.Threshold)
		} else {
			log.Printf("Service %s healthy (latency: %dms)", svc.Name, latency)
		}
		if resp != nil {
			err := resp.Body.Close()
			if err != nil {
				return
			}
		}
		time.Sleep(5 * time.Second)
	}
}
