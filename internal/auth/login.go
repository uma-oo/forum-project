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
	if r.URL.Path != "/log_in" || IsCookieSet(r, "session") {
		w.WriteHeader(http.StatusNotFound)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Page Not Found")
		return
	}
	User := r.FormValue("userName")
	Pass := r.FormValue("userPassword")
	fmt.Printf("username%v\n", User)
	fmt.Printf("Pass: %v\n", Pass)
	if User == "" || Pass == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Invalid Request")
		return

	}
	var Pas string

	err := database.Database.QueryRow("SELECT userName , userPassword  FROM users WHERE  userName= $1 ", User).Scan(&User, &Pas)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			pages.All_Templates.ExecuteTemplate(w, "login.html", "Invalid Password or username")
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", err)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(Pas), []byte(Pass)); err != nil {

		w.WriteHeader(http.StatusUnauthorized)
		pages.All_Templates.ExecuteTemplate(w, "login.html", "Invalid Password or username")
		return
	}
	Token := uuid.New().String()
	_, err = database.Database.Exec("UPDATE users set token = $1 where userName = $2", Token, User)
	if err != nil {
		fmt.Println(err)
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

	http.Redirect(w, r, "/", http.StatusFound)
}
