package auth

import (
	"log"
	"net/http"

	"forum/internal"
	"forum/internal/database"
	"forum/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	pages := internal.Pagess.All_Templates
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "405 method not allowed")
		return
	}
	userName := r.FormValue("userName")
	userPassword := r.FormValue("userPassword")
	Email := r.FormValue("userEmail")
	Token := uuid.New().String()

	Hach_pass, err := bcrypt.GenerateFromPassword([]byte(userPassword), 10)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	/*var userExist bool
	var emailExist bool
	stm, Err := database.Database.Prepare("SELECT EXISTS (SELECT 1 FROM users WHERE userEmail = ?)")
	if Err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}

	emailErr := stm.QueryRow(Email).Scan(&emailExist)
	stm, Err = database.Database.Prepare("SELECT EXISTS (SELECT 1 FROM users WHERE userName = ?)")
	if Err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	userErr := stm.QueryRow(userName).Scan(&userExist)

	if userErr != nil || emailErr != nil {
		logger.LogWithDetails(fmt.Errorf("%s, %s", userErr, emailErr))
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}

	if userExist || emailExist {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "User already exists")
		return
	} else {*/

	stm, err := database.Database.Prepare("INSERT INTO users (userName,userEmail,userPassword,token) VALUES (?, ?, ?, ? )")
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}

	_, err = stm.Exec(userName, Email, string(Hach_pass), Token)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "500 Internal Server Error")
		return
	}
	log.Printf("%s account has been created", userName)

	cookie := &http.Cookie{
		Name:   "token",
		Value:  Token,
		MaxAge: 3600,
		Path:   "/",
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}
