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
		log.Fatal(err)
	}
	defer logger.Close()
	// lets load the configuration
	configuration := config.LoadConfig()

	// server static files
	http.HandleFunc("/web/", handlers.Serve_Files)

	// routes for pages handling and rendering
	http.HandleFunc("/", handlers.Home)
	// http.HandleFunc("//{post_id}", handlers.Post)
	http.Handle("/create_post", middlewares.Auth_Middleware(http.HandlerFunc(handlers.CreatePost)))
	http.Handle("/my_posts", middlewares.Auth_Middleware(http.HandlerFunc(handlers.MyPosts)))
	http.Handle("/liked_posts", middlewares.Auth_Middleware(http.HandlerFunc(handlers.LikedPosts)))
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/register", handlers.Register)
	// routes for auth handlers in auth package we need to add the auth middleware for login and register likly deferrant
	http.HandleFunc("/auth/register", auth.Register)
	http.HandleFunc("/auth/log_in", auth.LogIn)
	http.HandleFunc("/auth/logout", auth.LogOut)
	// routes for forms actions
	// http.HandleFunc("/filter_posts", handlers.FilterPosts)
	http.Handle("/api/add_post", middlewares.Auth_Middleware(http.HandlerFunc(handlers.AddPost)))
	http.Handle("/api/react_to_post", middlewares.Auth_Middleware(http.HandlerFunc(handlers.PostReactions)))
	http.Handle("/api/add_post_comment", middlewares.Auth_Middleware(http.HandlerFunc(handlers.AddPostComment)))
	http.Handle("/api/react_to_comment", middlewares.Auth_Middleware(http.HandlerFunc(handlers.LikeComment)))
	// http.HandleFunc("/api/dislike_comment", handlers.DislikeComment)
	// fmt.Println("server is running on port 8080 ... http://localhost:8080")
	fmt.Printf("Server starting on port: %d >>> http://localhost:8080\n", configuration.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", configuration.Port), nil))
}
