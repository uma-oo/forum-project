// internal/auth/auth.go
package auth

import (
	"database/sql"
	"net/http"

	"forum/internal/database"
	"forum/internal/handlers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	pages  handlers.Pages
	server handlers.Server
)

func Signup_treatment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	if r.URL.Path != "/create_account" || IsCookieSet(r, "session") {
		w.WriteHeader(http.StatusNotFound)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Page Not Found")
		return

	}
	User := r.FormValue("userName")
	Pass := r.FormValue("userPassword")
	Email := r.FormValue("userEmail")
	if User == "" || Pass == "" || Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Bad Request")
		return
	}
	Hach_pass, err := bcrypt.GenerateFromPassword([]byte(Pass), 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Internal Server Error")
		return
	}
	user_already_exist := database.Database.QueryRow("SELECT userName FROM users WHERE userName = $1 ", User).Scan(&User)
	email_already_exist := database.Database.QueryRow("SELECT userEmail FROM users WHERE  userEmail = $1", Email).Scan(&Email)

	if user_already_exist != nil && email_already_exist != nil {
		if user_already_exist == sql.ErrNoRows && email_already_exist == sql.ErrNoRows {
			server.Log = false

			token := uuid.New().String()
			database.Database.Exec("INSERT INTO users  (userName, userPassword, userEmail) VALUES ($1, $2, $3)", User, Hach_pass, Email)
			http.Redirect(w, r, "/", http.StatusFound)
			http.SetCookie(w,
				&http.Cookie{
					Name:   "token",
					Value:  token,
					MaxAge: 3600,
				})

			return
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		pages.All_Templates.ExecuteTemplate(w, "signup.html", "username or email already exist")
		return

	}
}
