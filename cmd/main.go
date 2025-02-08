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
	http.Handle("/create_post", middlewares.Auth_Middleware(http.HandlerFunc(handlers.CreatePost)))
	http.Handle("/my_posts", middlewares.Auth_Middleware(http.HandlerFunc(handlers.MyPosts)))
	http.Handle("/liked_posts", middlewares.Auth_Middleware(http.HandlerFunc(handlers.LikedPosts)))
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/register", handlers.Register)

	// routes for auth handlers in auth package we need to add the auth middleware for login and register likly deferrant
	http.Handle("/auth/register", middlewares.Reg_Log_Middleware(http.HandlerFunc(auth.Register)))
	http.Handle("/auth/log_in", middlewares.Reg_Log_Middleware(http.HandlerFunc(auth.LogIn)))
	http.HandleFunc("/auth/logout", auth.LogOut)

	// routes for forms actions
	http.HandleFunc("/post", handlers.Post)
	http.HandleFunc("/filter_posts", handlers.FilterPosts)
	http.Handle("/api/add_post", middlewares.Auth_Middleware(http.HandlerFunc(handlers.AddPost)))
	http.Handle("/api/react_to_post", middlewares.Auth_Middleware(http.HandlerFunc(handlers.PostReactions)))
	http.Handle("/api/add_post_comment", middlewares.Auth_Middleware(http.HandlerFunc(handlers.AddPostComment)))
	http.Handle("/api/react_to_comment", middlewares.Auth_Middleware(http.HandlerFunc(handlers.LikeComment)))

	fmt.Printf("Server starting on port: %d >>> http://localhost:8080\n", configuration.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", configuration.Port), nil))
}
