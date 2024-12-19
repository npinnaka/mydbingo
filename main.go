package main

import (
	"fmt"
	"log"
	"time"

	"github.com/npinnaka/mydbingo/connection" // Replace with your package name
)

func main() {
	connStr := "postgres://dbuser:password@localhost:5432/mydb?sslmode=disable"

	db, err := connection.ConnectToDB(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	connection.StartConnectionHealthCheck(db, 5*time.Minute)

	// ... (rest of your code, using the `executeQueriesInTransaction` function)
}
