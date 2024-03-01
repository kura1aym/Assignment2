package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	db, err := sql.Open("postgres", "user=username password=password host=localhost dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Fatalf("Error: Unable to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		rows.Scan(&id, &name)
		fmt.Printf("User ID: %d, Name: %s\n", id, name)
	}

	var id int64
	var name string
	row := db.QueryRow("SELECT id, name FROM users WHERE id = $1", 1)

	if err := row.Scan(&id, &name); err == sql.ErrNoRows {
		fmt.Println("User not found")
	} else if err != nil {
		log.Fatalf("Error: Unable to execute query: %v", err)
	} else {
		fmt.Printf("User ID: %d, Name: %s\n", id, name)
	}
}
