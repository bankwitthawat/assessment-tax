package tax

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	url := os.Getenv("DATABASE_URL")
	log.Println("DATABASE_URL", url)
	DB, err = sql.Open("postgres", url)
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
