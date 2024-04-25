package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	log.Println("database connected.")

	// createTb := `
	// CREATE TABLE IF NOT EXISTS users ( id SERIAL PRIMARY KEY, name TEXT, age INT );
	// `
	// _, err = db.Exec(createTb)

	// if err != nil {
	// 	log.Fatal("can't create table", err)
	// }
}
