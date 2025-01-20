// internal/auth/auth.go
package auth

import (
	"database/sql"
	"forum/internal/database"
	"forum/internal/handlers"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)
var pages handlers.Pages
var server handlers.Server
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
	User := r.FormValue("username")
	Pass := r.FormValue("password")
	Email := r.FormValue("email")
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
	user_already_exist := database.Database.QueryRow("SELECT username FROM users WHERE username = $1 ", User).Scan(&User)
	email_already_exist :=database.Database.QueryRow("SELECT email FROM users WHERE  email = $1", Email).Scan(&Email)

	if user_already_exist != nil && email_already_exist != nil {
		if user_already_exist == sql.ErrNoRows && email_already_exist == sql.ErrNoRows {
			server.Log = false
		
			token := uuid.New().String()
			database.Database.Exec("INSERT INTO users  (username, password, email , token) VALUES ($1, $2, $3, $4)", User, Hach_pass, Email, token)
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
