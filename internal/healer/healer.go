package healer

import (
	"context"
	"log"
	"time"

	"github.com/Dokuqui/healnix/pkg/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Healer struct {
	dockerCli *client.Client
}

func NewHealer() *Healer {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Printf("Failed to create Docker client: %v", err)
		return &Healer{}
	}
	return &Healer{dockerCli: cli}
}

func (h *Healer) Heal(status *types.ServiceStatus, containerName string) bool {
	if h.dockerCli == nil {
		log.Printf("No Docker client available, skipping heal for %s", status.Name)
		return false
	}

	log.Printf("Attempting to heal %s by restarting container %s...", status.Name, containerName)

	ctx := context.Background()
	timeout := int(30 * time.Second / time.Second)
	stopOptions := container.StopOptions{
		Timeout: &timeout,
	}
	err := h.dockerCli.ContainerRestart(ctx, containerName, stopOptions)
	attempt := types.HealAttemp{Timestamp: time.Now()}

	if err != nil {
		log.Printf("Failed to heal %s (container %s): %v", status.Name, containerName, err)
		attempt.Success = false
		attempt.Error = err.Error()
	} else {
		log.Printf("Successfully healed %s by restarting container %s", status.Name, containerName)
		attempt.Success = true
	}

	status.HealHistory = append(status.HealHistory, attempt)
	return attempt.Success
}
