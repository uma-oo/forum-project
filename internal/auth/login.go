package auth

import (
	"database/sql"
	"fmt"
	"net/http"

	"forum/internal/database"
	"forum/internal/handlers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Log_in(w http.ResponseWriter, r *http.Request) {
	pages := handlers.Pagess
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	if r.URL.Path != "/api/log_in" || IsCookieSet(r, "session") {
		w.WriteHeader(http.StatusNotFound)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Page Not Found")
		return
	}

	username := r.FormValue("userName")
	password_got := r.FormValue("userPassword")
	fmt.Printf("username: %v\n", username)
	fmt.Printf("password_got: %v\n", password_got)
	if username == "" || password_got == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Invalid Request")
		return

	}

	var password string

	err := database.Database.QueryRow("SELECT userName , userPassword  FROM users WHERE  userName= $1 ", username).Scan(&username, &password)
	fmt.Printf("password: %v\n", bcrypt.CompareHashAndPassword([]byte(password), []byte(password_got)))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			pages.All_Templates.ExecuteTemplate(w, "login.html", "Invalid Credentials")
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", err)
		return
		
	}
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(password_got)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Invalid Credentials")
		return
	}
	Token := uuid.New().String()
	statement, err := database.Database.Prepare("UPDATE users SET token = ? where userName = ? ")
	if err != nil {
		fmt.Printf("err in statement of the database: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Internal Server Error")
		return
	}
	_, err = statement.Exec(Token, username)
	if err != nil {
		fmt.Printf("err in the exec of database: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Internal Server Error")
		return
	}

	fmt.Println("Token:", Token)
	cookie := &http.Cookie{
		Name:   "token",
		Value:  Token,
		MaxAge: 3600,
	}
	http.SetCookie(w, cookie)
	// server.Log = false
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
