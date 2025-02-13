// internal/middleware/middleware.go
package middlewares

import (
	"net/http"
	"time"

	"forum/internal"
	"forum/internal/auth"
	"forum/internal/models"
	"forum/internal/utils"
	"forum/pkg/logger"
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
		user, err := utils.UserData(r, "token", "")
		if err != nil {
			logger.LogWithDetails(err)
			utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
			return
		}
		ok, err := utils.CheckTokenExpired(user)
		if err != nil {
			logger.LogWithDetails(err)
			utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
			return
		}
		if ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
