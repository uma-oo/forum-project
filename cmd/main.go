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
}

func main() {
	// Get the current working directory
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/create_account", auth.Signup_treatment)
	http.HandleFunc("/log_in", auth.Log_in)
	http.HandleFunc("/create_post", handlers.CreatePost)
	http.HandleFunc("/filterPost", handlers.FilterPosts)
	http.HandleFunc("/myposts", handlers.MyPosts)
	http.HandleFunc("/likedposts", handlers.LikedPosts)
	http.HandleFunc("/categorizePost", handlers.CategorizePosts)
	http.HandleFunc("/settings", handlers.Settings)
	http.HandleFunc("/logout", handlers.Logout)
	http.HandleFunc("/static/", handlers.Serve_Static)
	fmt.Println("server is running on port 8080 ... http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
