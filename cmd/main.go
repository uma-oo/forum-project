// cmd/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/internal/auth"
	"forum/internal/database"
	"forum/internal/handlers"
	"forum/internal/middlewares"
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
	

	// Handlers for rendering :
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/create_post", handlers.Create_Post)
	http.HandleFunc("/static/", handlers.Serve_Static)
	// handlers that does some checking 
	http.Handle("/api/create_comment", middlewares.Auth_Middleware(http.HandlerFunc(handlers.CreateComment)))

	http.HandleFunc("/api/log_in", auth.Log_in)
	http.HandleFunc("/api/create_post", handlers.CreatePost)
	http.HandleFunc("/api/create_account", auth.Register)
	http.HandleFunc("/api/like_post", handlers.LikePost)
	http.HandleFunc("/api/like_comment", handlers.LikeComment)
	http.HandleFunc("/api/dislike_post", handlers.DislikePost)
	http.HandleFunc("/api/dislike_comment", handlers.DislikeComment)
	// fmt.Println("server is running on port 8080 ... http://localhost:8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", configuration.Port), nil))
}
