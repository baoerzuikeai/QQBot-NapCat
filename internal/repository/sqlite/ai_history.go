package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createTableSQL = `
	CREATE TABLE IF NOT EXISTS ai_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		role TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		session_id TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_user_id ON ai_history (user_id);
	CREATE INDEX IF NOT EXISTS idx_session_id ON ai_history (session_id); 
	`
)

type sqliteAIHistoryRepository struct {
	db *sql.DB
}

func NewSqliteAIHistoryRepository(dateSourceName string) (*sqliteAIHistoryRepository, error) {
	db, err := sql.Open("sqlite3", dateSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %v", err)
	}
	log.Println("AI history table created successfully or already exist")
	return &sqliteAIHistoryRepository{db: db}, nil

}

func (s *sqliteAIHistoryRepository) SaveAIhistory(userID int64, content, role, sessionID string) error {
	stmt, err := s.db.Prepare("INSERT INTO ai_history (user_id, content, role, timestamp, session_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()
	timestamp := time.Now().Unix()
	_, err = stmt.Exec(userID, content, role, timestamp, sessionID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}
	log.Printf("AI history saved for user ID %d at timestamp %d", userID, timestamp)
	return nil
}

func (s *sqliteAIHistoryRepository) GetAIHistoryByUserID(userID int64) ([]AIHistory, error) {
	stmt, err := s.db.Prepare("SELECT id, user_id, content, role, timestamp, session_id FROM ai_history WHERE user_id = ? ORDER BY timestamp DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()
	var history []AIHistory
	for rows.Next() {
		var h AIHistory
		err := rows.Scan(&h.ID, &h.UserID, &h.Content, &h.Role, &h.Timestamp, &h.SessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		history = append(history, h)
	}
	return history, nil
}
