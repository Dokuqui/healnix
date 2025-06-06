package monitor

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Dokuqui/healnix/internal/healer"
	"github.com/Dokuqui/healnix/internal/history"
	"github.com/Dokuqui/healnix/pkg/config"
	"github.com/Dokuqui/healnix/pkg/types"
)

type Monitor struct {
	Services []config.Service
	Statuses map[string]*types.ServiceStatus
	History  history.History
	mu       sync.Mutex
}

func NewMonitor(cfg *config.Config, hist history.History) *Monitor {
	statuses := make(map[string]*types.ServiceStatus)
	for _, svc := range cfg.Services {
		statuses[svc.Name] = &types.ServiceStatus{Name: svc.Name, Healthy: true}
	}
	return &Monitor{Services: cfg.Services, Statuses: statuses, History: hist}
}

func (m *Monitor) Start() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	for _, svc := range m.Services {
		go m.checkService(svc)
	}

	<-stop
	log.Println("Shutting down monitor...")
}

func (m *Monitor) checkService(svc config.Service) {
	client := &http.Client{Timeout: 2 * time.Second}
	for {
		start := time.Now()
		resp, err := client.Get(svc.Endpoint)
		latency := time.Since(start).Milliseconds()

		m.mu.Lock()
		status := m.Statuses[svc.Name]
		status.LastCheck = time.Now()
		status.Latency = latency

		if err != nil {
			status.ConsecutiveFails++
			status.Healthy = false
			log.Printf("Service %s failed: %v (fail count: %d)", svc.Name, err, status.ConsecutiveFails)
		} else if latency > int64(svc.Threshold) {
			status.ConsecutiveFails++
			status.Healthy = false
			log.Printf("Service %s unhealthy (latency: %dms > %dms, fail count: %d)", svc.Name, latency, svc.Threshold, status.ConsecutiveFails)
		} else {
			status.ConsecutiveFails = 0
			status.Healthy = true
			log.Printf("Service %s healthy (latency: %dms)", svc.Name, latency)
		}

		if !status.Healthy && svc.Heal != "" && status.ConsecutiveFails >= svc.FailureThreshold && !status.HealingInProgress {
			status.HealingInProgress = true
			m.mu.Unlock()
			healer := healer.NewHealer()
			success := healer.Heal(status, svc.ContainerName)
			m.mu.Lock()
			status.HealingInProgress = false
			if success {
				status.ConsecutiveFails = 0
			}

			if len(status.HealHistory) > 0 {
				latestAttempt := status.HealHistory[len(status.HealHistory)-1]
				m.History.SaveHealAttempt(svc.Name, latestAttempt)
			}
			log.Printf("Service %s heal history: %d attempts", svc.Name, len(status.HealHistory))
		}
		m.mu.Unlock()

		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(5 * time.Second)
	}
}

func (m *Monitor) GetStatus(name string) *types.ServiceStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Statuses[name]
}

func (m *Monitor) PrintStatus() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, status := range m.Statuses {
		log.Printf("Service: %s, Healthy: %v, Latency: %dms, LastCheck: %s, Fails: %d/%d",
			status.Name, status.Healthy, status.Latency, status.LastCheck.Format(time.RFC3339),
			status.ConsecutiveFails, m.Services[0].FailureThreshold)
	}
}
