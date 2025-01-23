// cmd/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/internal/auth"
	"forum/internal/database"
	"forum/internal/handlers"
)

func init() {
	database.Create_database()
	handlers.ParseTemplates()
}

func main() {
	// Get the current working directory
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/create_account", auth.Register)
	http.HandleFunc("/log_in", auth.Log_in)
	http.HandleFunc("/static/", handlers.Serve_Static)
	http.HandleFunc("/create_comment", handlers.CreateComment)
	fmt.Println("server is running on port 8080 ... http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
