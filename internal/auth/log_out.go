package auth

import (
	"log"
	"net/http"

	"forum/internal/models"
	"forum/internal/utils"
)

func LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RenderTemplate(w, "error.html",models.MethodNotAllowed,http.StatusMethodNotAllowed)
		return
	}

	// lets check in first that is already have a session
	if utils.IsCookieSet(r, "token") {
		http.SetCookie(w, &http.Cookie{
			Name:   "token", // name of the cookie
			Value:  "",      // clear the cookie value
			MaxAge: -1,      // set expiration time to a time in the past
			Path:   "/",     // scope of the cookie
		})
		log.Print("A User logged out")

	} else {
		utils.RenderTemplate(w, "error.html",models.PageNotFound,http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
