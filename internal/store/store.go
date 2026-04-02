package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct { db *sql.DB }

type Embargo struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	ReleaseAt    string   `json:"release_at"`
	Status       string   `json:"status"`
	CreatedAt    string   `json:"created_at"`
}

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	dsn := filepath.Join(dataDir, "embargo.db") + "?_journal_mode=WAL&_busy_timeout=5000"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS embargos (
			id TEXT PRIMARY KEY,\n\t\t\ttitle TEXT DEFAULT '',\n\t\t\tcontent TEXT DEFAULT '',\n\t\t\trelease_at TEXT DEFAULT '',\n\t\t\tstatus TEXT DEFAULT 'pending',
			created_at TEXT DEFAULT (datetime('now'))
		)`)
	if err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }

func (d *DB) Create(e *Embargo) error {
	e.ID = genID()
	e.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err := d.db.Exec(`INSERT INTO embargos (id, title, content, release_at, status, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		e.ID, e.Title, e.Content, e.ReleaseAt, e.Status, e.CreatedAt)
	return err
}

func (d *DB) Get(id string) *Embargo {
	row := d.db.QueryRow(`SELECT id, title, content, release_at, status, created_at FROM embargos WHERE id=?`, id)
	var e Embargo
	if err := row.Scan(&e.ID, &e.Title, &e.Content, &e.ReleaseAt, &e.Status, &e.CreatedAt); err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Embargo {
	rows, err := d.db.Query(`SELECT id, title, content, release_at, status, created_at FROM embargos ORDER BY created_at DESC`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []Embargo
	for rows.Next() {
		var e Embargo
		if err := rows.Scan(&e.ID, &e.Title, &e.Content, &e.ReleaseAt, &e.Status, &e.CreatedAt); err != nil {
			continue
		}
		result = append(result, e)
	}
	return result
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM embargos WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM embargos`).Scan(&n)
	return n
}
