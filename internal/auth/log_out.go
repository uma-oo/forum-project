package auth

import (
	"log"
	"net/http"

	"forum/internal/handlers"
)

func LogOut(w http.ResponseWriter, r *http.Request) {
	pages := handlers.Pagess
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	// lets check in first that is already have a session
	if IsCookieSet(r, "token") {
		http.SetCookie(w, &http.Cookie{
			Name:   "token", // name of the cookie
			Value:  "",      // clear the cookie value
			MaxAge: -1,      // set expiration time to a time in the past
			Path:   "/",     // scope of the cookie
			// HttpOnly: true,            // prevent JavaScript access
			// Secure:   true,            // ensure cookie is only sent over HTTPS
		})
		log.Print("A User logged out")
		http.Redirect(w, r, "/", http.StatusFound)
	}
	return
}
