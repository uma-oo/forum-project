// internal/middleware/middleware.go
package middlewares

import (
	"log"
	"net/http"
	"strings"
)

func Auth_Midlleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// lets log the request
		log.Printf("%v recieved a request >>   Method: %s, URL: %s, Client IP: %s", r.Header.Get("Date"), r.Method, r.URL.Path, r.RemoteAddr)
		log.Print("------------------------")
		// now lets check for authorisation
		authorisationHrader := r.Header.Get("Authorization")
		if authorisationHrader == "" || !strings.HasPrefix(authorisationHrader, "aghyor ") || !validToken(authorisationHrader[7:]) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// now call the next hamdel in the chain
		next.ServeHTTP(w, r)
	})
}

func validToken(token string) bool {
	return token == "boda"
}
