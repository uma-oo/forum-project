package auth

import (
	"database/sql"
	"net/http"

	"forum/internal/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Log_in(w http.ResponseWriter, r *http.Request) {
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
	User := r.FormValue("username")
	Pass := r.FormValue("password")

	if User == "" || Pass == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Invalid Request")
		return

	}
	var Pas string

	err := database.Database.QueryRow("SELECT username , password  FROM users WHERE  username= $1 ", User).Scan(&User, &Pas)
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
	cookie := &http.Cookie{
		Name:   "token",
		Value:  Token,
		MaxAge: 3600,
	}
	http.SetCookie(w, cookie)
	server.Log = false

	http.Redirect(w, r, "/", http.StatusFound)
}
