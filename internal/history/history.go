package history

import (
	"database/sql"
	"log"
	"time"

	"github.com/Dokuqui/healnix/pkg/types"
	_ "modernc.org/sqlite"
)

type History interface {
	SaveHealAttempt(serviceName string, attempt types.HealAttemp) error
	PrintHistory(serviceName string)
	Predict(serviceName string, lookback time.Duration) []string
	Close() error
}

type SQLiteHistory struct {
	db *sql.DB
}

func NewHistory(dbPath string) (*SQLiteHistory, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Printf("Failed to open SQLite database: %v", err)
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
		log.Printf("Failed to create heal_history table: %v", err)
		return nil, err
	}
	return &SQLiteHistory{db: db}, nil
}

func (h *SQLiteHistory) SaveHealAttempt(serviceName string, attempt types.HealAttemp) error {
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

func (h *SQLiteHistory) PrintHistory(serviceName string) {
	rows, err := h.db.Query("SELECT timestamp, success, error FROM heal_history WHERE service_name = ?", serviceName)
	if err != nil {
		log.Printf("Failed to query history for %s: %v", serviceName, err)
		return
	}
	defer rows.Close()

	log.Printf("Healing history for %s:", serviceName)
	for rows.Next() {
		var timestamp time.Time
		var success bool
		var errorMsg string
		if err := rows.Scan(&timestamp, &success, &errorMsg); err != nil {
			log.Printf("Failed to scan history row: %v", err)
			continue
		}
		if success {
			log.Printf("  %s: Success", timestamp.Format(time.RFC3339))
		} else {
			log.Printf("  %s: Failed (%s)", timestamp.Format(time.RFC3339), errorMsg)
		}
	}
}

func (h *SQLiteHistory) Predict(serviceName string, lookback time.Duration) []string {
	end := time.Now()
	start := end.Add(-lookback)
	rows, err := h.db.Query(`
        SELECT timestamp, success
        FROM heal_history
        WHERE service_name = ? AND timestamp >= ? AND timestamp <= ?`,
		serviceName, start, end,
	)
	if err != nil {
		log.Printf("Failed to query history for %s: %v", serviceName, err)
		return nil
	}
	defer rows.Close()

	var failures int
	var total int
	for rows.Next() {
		var timestamp time.Time
		var success bool
		if err := rows.Scan(&timestamp, &success); err != nil {
			continue
		}
		total++
		if !success {
			failures++
		}
	}

	suggestions := []string{}
	if total > 0 {
		failureRate := float64(failures) / float64(total)
		if failureRate > 0.5 {
			suggestions = append(suggestions, "High failure rate detected. Consider scaling up or checking service configuration.")
		}
		if total > 5 {
			suggestions = append(suggestions, "Frequent healing attempts. Investigate resource limits or external dependencies.")
		}
	}
	return suggestions
}

func (h *SQLiteHistory) Close() error {
	return h.db.Close()
}
