// cmd/main.go
package main

import (
	"fmt"
	"forum/internal/handlers"
	"log"
	"net/http"
)

func main() {

	// Get the current working directory
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/sign_in", handlers.Sign_In)
	http.HandleFunc("/createpost", handlers.CreatePost)
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
