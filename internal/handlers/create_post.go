package handlers

import (
	"fmt"
	"net/http"

	"forum/internal/database"
)

type Post struct {
	UserId        int
	Content       string
	Title         string
	TotalLikes    int
	TotalDislikes int
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	pages := Pagess.All_Templates
	post := NewPost()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	post.Content = r.FormValue("postBody")
	post.Title = r.FormValue("postTitle")
	// lets check for emptyness
	if post.Content == "" || post.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "bad request")
		return
	}
	post.UserId = 1

	// lets insert this data to our database
	_, err := database.Database.Exec("INSERT INTO posts (user_id,title,content,total_likes,total_dislikes) VALUES ( ?,?,?,?,?)", post.UserId, post.Title, post.Content, post.TotalLikes, post.TotalDislikes)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func NewPost() *Post {
	return &Post{
		UserId:        0,
		Content:       "",
		Title:         "",
		TotalLikes:    0,
		TotalDislikes: 0,
	}
}


func AllPosts(){
	
}