// internal/middleware/middleware.go
package middlewares

import (
	"net/http"

	"forum/internal/utils"
)

func Auth_Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request contains a token
		if !utils.IsCookieSet(r, "token") {
			http.Redirect(w, r, "/login", http.StatusFound) // Use StatusFound (302)
			return
		}
		next.ServeHTTP(w, r)
	})
}
