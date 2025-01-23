package auth

import (
	"log"
	"net/http"
	"time"

	"forum/internal/handlers"
)

func Log_out(w http.ResponseWriter, r *http.Request) {
	pages := handlers.Pagess
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	// lets check in first that is already have a session
	if IsCookieSet(r, "token") {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",         // name of the cookie
			Value:    "",              // clear the cookie value
			Expires:  time.Unix(0, 0), // set expiration time to a time in the past
			Path:     "/",             // scope of the cookie
			HttpOnly: true,            // prevent JavaScript access
			Secure:   true,            // ensure cookie is only sent over HTTPS
		})
		log.Print("A User logged out")
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		w.WriteHeader(http.StatusNotFound)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "page not fount")
		return
	}
}
