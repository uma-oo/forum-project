// internal/middleware/middleware.go
package middlewares

import (
	"net/http"
	"time"

	"forum/internal"
	"forum/internal/auth"
	"forum/internal/utils"
)

var auth_rateLimiter = auth.NewRateLimiter(30, time.Second) // 30 requests per second limit
func Auth_Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request contains a token
		ip := r.RemoteAddr
		pages := internal.Templates
		// Check rate limit (applies to both login and registration)
		if auth_rateLimiter.CheckRateLimit(ip) {
			w.WriteHeader(http.StatusTooManyRequests)
			pages.ExecuteTemplate(w, "error.html", "Too many requests. Please try again later.")
			return
		}
		if !utils.IsCookieSet(r, "token") {
			http.Redirect(w, r, "/login", http.StatusFound) // Use StatusFound (302)
			return
		}
		next.ServeHTTP(w, r)
	})
}
