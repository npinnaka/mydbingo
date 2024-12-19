package connection

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Query struct {
	Query  string
	Params []interface{}
}

func ConnectToDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Set maximum number of open connections and idle connections in the pool
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	// Test the connection before using it
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ExecuteQueriesInTransaction(db *sql.DB, queries []Query) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v", r)
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	for _, query := range queries {
		_, err = tx.Exec(query.Query, query.Params...)
		if err != nil {
			return err
		}
	}

	return nil
}

func StartConnectionHealthCheck(db *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		err := db.Ping()
		if err != nil {
			log.Printf("Error pinging database: %v", err)
			// Handle the error, e.g., retry connection, alert, etc.
			// Consider using a backoff strategy for retries
			time.Sleep(5 * time.Second) // Adjust the backoff time as needed
			err = db.Ping()
			if err != nil {
				log.Printf("Failed to recover connection: %v", err)
				// Trigger alerts or other recovery actions
			}
		} else {
			log.Println("Database connection is healthy")
		}
	}
}
