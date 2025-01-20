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
		fmt.Println(0)
		log.Fatal(err)
	
	}

	_, usererr := Database.Exec( `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`)
	if usererr != nil {
		fmt.Println(1)
		fmt.Println(usererr)
	
	}
	_, posterr := Database.Exec(`CREATE TABLE IF NOT EXISTS posts ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "user_id" INTEGER NOT NULL, "title" TEXT, "content" TEXT,FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE CASCADE);`)
	if posterr != nil {
		fmt.Println(2)
		log.Fatal(posterr)
		
	}
	_, likeerr := Database.Exec(`CREATE TABLE IF NOT EXISTS likes ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,"post_id" INTEGER NOT NULL , "like" INTEGER "dislike" INTEGER, FOREIGN KEY ("post_id") REFERENCES posts("id") ON DELETE CASCADE);`)
	
	if likeerr != nil {
		fmt.Println(3)
		log.Fatal(likeerr)
		
	}
	_,commenterr := Database.Exec(`CREATE TABLE IF NOT EXISTS comments ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,"post_id" INTEGER NOT NULL , "comment"  TEXT, FOREIGN KEY ("post_id") REFERENCES posts("id") ON DELETE CASCADE);`)
	if commenterr != nil {
		fmt.Println(4)
		log.Fatal(commenterr)
		
	}
	_, gategoryerr := Database.Exec(`CREATE TABLE IF NOT EXISTS gategories ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,"post_id" INTEGER NOT NULL, "gategory" TEXT, FOREIGN KEY ("post_id") REFERENCES posts("id") ON DELETE CASCADE);`)

	if gategoryerr != nil {
		fmt.Println(5)
		log.Fatal(gategoryerr)
		
	}

}
