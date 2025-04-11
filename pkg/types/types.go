package types

import "time"

type HealAttemp struct {
	Timestamp time.Time
	Success   bool
	Error     string
}

type ServiceStatus struct {
	Name              string
	Healthy           bool
	Latency           int64
	LastCheck         time.Time
	ConsecutiveFails  int
	HealHistory       []HealAttemp
	HealingInProgress bool
}
