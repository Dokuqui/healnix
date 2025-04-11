package history

import (
	"database/sql"
	"log"

	"github.com/Dokuqui/healnix/pkg/types"
	_ "modernc.org/sqlite"
)

type History struct {
	db *sql.DB
}

func NewHistory(dbPath string) (*History, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS heal_history (
			service_name TEXT,
			timestamp DATETIME,
			success BOOLEAN,
			error TEXT
		)
	`)
	if err != nil {
		return nil, err
	}
	return &History{db: db}, nil
}

func (h *History) SaveHealAttempt(serviceName string, attempt types.HealAttemp) error {
	_, err := h.db.Exec(
		"INSERT INTO heal_history (service_name, timestamp, success, error) VALUES (?, ?, ?, ?)",
		serviceName, attempt.Timestamp, attempt.Success, attempt.Error,
	)
	if err != nil {
		log.Printf("Failed to save heal attempt for %s: %v", serviceName, err)
		return err
	}
	log.Printf("Saved heal attempt for %s: success=%v", serviceName, attempt.Success)
	return nil
}

func (h *History) Close() {
	h.db.Close()
}
