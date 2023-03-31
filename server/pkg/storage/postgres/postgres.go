package postgres

import (
	"database/sql"
	"fmt"
)

type VisitsDB struct {
	DB *sql.DB
}

func NewVisitsDB(dbUrl string) (*VisitsDB, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS visits (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			timestamp TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("unable to create table: %v", err)
	}

	return &VisitsDB{DB: db}, nil
}

func (vdb *VisitsDB) RecordVisit(name string) error {
	_, err := vdb.DB.Exec(`INSERT INTO visits (name, timestamp) VALUES ($1, NOW())`, name)
	if err != nil {
		return fmt.Errorf("unable to record visit: %w", err)
	}
	return nil
}

func (vdb *VisitsDB) CountVisits(name string) (int, error) {
	var count int
	err := vdb.DB.QueryRow(`SELECT COUNT(*) FROM visits WHERE name=$1`, name).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("unable to count visits: %v", err)
	}
	return count, nil
}
