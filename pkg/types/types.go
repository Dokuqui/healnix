package types

import "time"

type ServiceStatus struct {
	Name      string
	Healthy   bool
	Latency   int64
	LastCheck time.Time
}
