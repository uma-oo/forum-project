package auth

import (
	"log"
	"net/http"

	"forum/internal/database"
	"forum/internal/handlers"
	"forum/internal/models"
	"forum/pkg/logger"

	"github.com/google/uuid"
)

var LogRegFormsErrors models.FormErrors

func LogIn(w http.ResponseWriter, r *http.Request) {
	pages := handlers.Pagess
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "405 Method Not Allowed")
		return
	}

	UserName := r.FormValue("userName")
	Token := uuid.New().String()
	stm, err := database.Database.Prepare("UPDATE users SET token = ? where userName = ? ")
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}

	_, err = stm.Exec(Token, UserName)
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
