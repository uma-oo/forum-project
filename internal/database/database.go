// internal/database/database.go
// internal/database/db.go
package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func init() {
	var err error
	Database, err = sql.Open("sqlite3", "./internal/database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	
	_, usererr := Database.Exec(`CREATE TABLE IF NOT EXISTS users ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "username" TEXT, "email" TEXT "password" TEXT);`)
	if usererr != nil {
		log.Fatal(usererr)
	}
	_, posterr := Database.Exec(`CREATE TABLE IF NOT EXISTS posts ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "user_id" INTEGER NOT NULL, "title" TEXT, "content" TEXT,FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE CASCADE);`)
	if posterr != nil {
		log.Fatal(posterr)

	}
	_,likeerr := Database.Exec(`CREATE TABLE IF NOT EXISTS likes ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,"post_id" INTEGER NOT NULL , "like" INTEGER "dislike" INTEGER, FOREIGN KEY ("post_id") REFERENCES posts("id") ON DELETE CASCADE);`)
	if likeerr != nil {
		log.Fatal(likeerr)
	}
	_, gategoryerr := Database.Exec(`CREATE TABLE IF NOT EXISTS gategories ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,"post_id" INTEGER NOT NULL, "gategory" TEXT, FOREIGN KEY ("post_id") REFERENCES posts("id") ON DELETE CASCADE);`)
	
	if gategoryerr != nil {
	
		log.Fatal(gategoryerr)
	}

}
