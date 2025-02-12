package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	// execute a query
	_, err = db.Exec("CRAETE TABLE IF NOT EXISTS users (username TEXT, password TEXT)")
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/login", loginHandler)
	fmt.Printf("Server starting on port: 8080 >>> http://localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	// execute a query
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		fmt.Println(err)
	}

	var user User
	err = db.QueryRow("SELECT * FROM users WHERE username = '$username' AND password = '$password'", username, password).Scan(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Welcome, %s!", user.Username)
}

type User struct {
	Username string
	Password string
}
