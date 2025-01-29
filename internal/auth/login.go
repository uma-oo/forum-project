package auth

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"forum/internal/database"
	"forum/internal/handlers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Log_in(w http.ResponseWriter, r *http.Request) {
	pages := handlers.Pagess
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	if IsCookieSet(r, "token") {
		w.WriteHeader(http.StatusNotFound)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Page Not Found hhhhh")
		return
	}
	UserName := r.FormValue("userName")
	Password := r.FormValue("userPassword")

	if UserName == "" || Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Bad Request")
		return

	}
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
	// statement, err := database.Database.Prepare("UPDATE users SET token = ? where userName = ? ")
	// if err != nil {
	// 	fmt.Printf("err in statement of the database: %v\n", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	pages.All_Templates.ExecuteTemplate(w, "error.html", "Internal Server Error")
	// 	return
	// }
	// _, err = statement.Exec(Token, username)
	// if err != nil {
	// 	fmt.Printf("err in the exec of database: %v\n", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	pages.All_Templates.ExecuteTemplate(w, "error.html", "Internal Server Error")
	// 	return
	// }
	_, err = database.Database.Exec("UPDATE users SET token = ? WHERE userName = ?", Token, username)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	cookie := &http.Cookie{
		Name:   "token",
		Value:  Token,
		MaxAge: 3600,
		Path:   "/",
	}
	http.SetCookie(w, cookie)
	r.AddCookie(cookie)
	log.Println(UserName, "logged in")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
