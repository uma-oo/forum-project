// cmd/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/internal/auth"
	"forum/internal/database"
	"forum/internal/handlers"
	"forum/pkg/config"
	"forum/pkg/logger"
)

func init() {
	database.Create_database()
	handlers.ParseTemplates()
}

func main() {
	// Get the current working directory
	logger, err := logger.Create_Logger()
	if err != nil {
		fmt.Println("here")
		log.Fatal(err)

	}
	defer logger.Close()

	// lets load the configuration
	configuration := config.LoadConfig()
	fmt.Printf("Server starting on port: %d >>> http://localhost:8080\n", configuration.Port)

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/create_account", auth.Register)
	http.HandleFunc("/log_in", auth.Log_in)
	http.HandleFunc("/create_Post", handlers.Craete_post)
	http.HandleFunc("/static/", handlers.Serve_Static)
	//fmt.Println("server is running on port 8080 ... http://localhost:8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", configuration.Port), nil))
}
