package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, os.Getenv("PG_USER"), os.Getenv("PG_PW"), os.Getenv("PG_DB"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Users")
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	fmt.Println("Users:")
	for rows.Next() {
		var user_id int
		var first_name, last_name, email, password, auth string
		var is_affiliated_with_rso bool

		err = rows.Scan(&user_id, &first_name, &last_name, &email, &password, &auth, &is_affiliated_with_rso)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		fmt.Printf("ID: %d, First Name: %s, Last Name: %s, Email: %s, Password: %s, Auth: %s, Affiliated with RSO: %t\n",
			user_id, first_name, last_name, email, password, auth, is_affiliated_with_rso)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Row iteration error: %v", err)
	}
}
