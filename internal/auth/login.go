package auth

import (
	"database/sql"
	"log"
	"net/http"

	"forum/internal/database"
	"forum/internal/handlers"
	"forum/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	pages := handlers.Pagess
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}

	UserName := r.FormValue("userName")
	Password := r.FormValue("userPassword")

	// if UserName == "" || Password == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	pages.All_Templates.ExecuteTemplate(w, "error.html", "Bad Request") should be removed and apliuc this logic auth middleware
	// 	return
	// }

	var pasword string
	var username string

	err := database.Database.QueryRow("SELECT userName , userPassword  FROM users WHERE  username= $1 ", UserName).Scan(&username, &pasword)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			pages.All_Templates.ExecuteTemplate(w, "error.html", "user not exist") // should execute login page here for no rows err
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", err)
		return

	}
	if err := bcrypt.CompareHashAndPassword([]byte(pasword), []byte(Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Invalid Password or username or exist")
		return
	}
	Token := uuid.New().String()
	stm, err := database.Database.Prepare("UPDATE users SET token = ? where userName = ? ")
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}

	_, err = stm.Exec(Token, username)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "500 internal server error")
		return
	}
	cookie := &http.Cookie{
		Name:   "token",
		Value:  Token,
		MaxAge: 3600,
		Path:   "/",
	}
	http.SetCookie(w, cookie)

	log.Println(UserName, "logged in")
	http.Redirect(w, r, "/", http.StatusFound)
}
