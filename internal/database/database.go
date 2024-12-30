// internal/database/database.go
// internal/database/db.go
package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func init() {
	var err error
	Database, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	createTable := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"username" TEXT,
		"password" TEXT
	);`
	_, err = Database.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}
