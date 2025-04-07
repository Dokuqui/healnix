package healer

import (
	"log"
	"os/exec"

	"github.com/Dokuqui/healnix/pkg/types"
)

type Healer struct{}

func NewHealer() *Healer {
	return &Healer{}
}

func (h *Healer) Heal(status types.ServiceStatus) {
	log.Printf("Attempting to heal %s...", status.Name)
	cmd := exec.Command("echo", "Restarting", status.Name)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to heal %s: %v", status.Name, err)
	} else {
		log.Printf("Successfully healed %s", status.Name)
	}
}
