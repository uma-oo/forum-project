package middlewares

import (
	"net/http"
	"time"

	"forum/internal/auth"
	"forum/internal/handlers"
	"forum/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

var rateLimiter = auth.NewRateLimiter(5, time.Minute) // 5 requests per minute limit

// Custom middleware to validate the registration form
func Reg_Log_Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pages := handlers.Pagess.All_Templates
		ip := r.RemoteAddr
		// Check rate limit (applies to both login and registration)
		if rateLimiter.CheckRateLimit(ip) {
			w.WriteHeader(http.StatusTooManyRequests)
			pages.ExecuteTemplate(w, "error.html", "Too many requests. Please try again later.")
			return
		}
		username := r.FormValue("userName")
		password := r.FormValue("userPassword")
		email := r.FormValue("userEmail")
		// Validate the registration form
		if r.URL.Path == "/auth/register" {
			if !utils.IsValidUsername(username) {
				w.WriteHeader(http.StatusBadRequest)
				pages.ExecuteTemplate(w, "register.html", "Username is invalid.")
				return
			}
			_, exist := utils.IsExist("userName", "", username)
			if exist {
				w.WriteHeader(http.StatusBadRequest)
				pages.ExecuteTemplate(w, "register.html", "Username is already taken.")
				return
			}

			if !utils.IsValidEmail(email) {
				w.WriteHeader(http.StatusBadRequest)
				pages.ExecuteTemplate(w, "register.html", "Email is invalid.")
				return
			}
			_, exist = utils.IsExist("userEmail", "", email)
			if exist {
				w.WriteHeader(http.StatusBadRequest)
				pages.ExecuteTemplate(w, "register.html", "Email is already taken.")
				return
			}
			// Validate password
			if !utils.IsStrongPassword(password) {
				w.WriteHeader(http.StatusBadRequest)
				pages.ExecuteTemplate(w, "register.html", "Password is too weak.")
				return
			}

		}
		// validate login form
		if r.URL.Path == "/auth/log_in" {
			// check if username is exist
			pass, exist := utils.IsExist("userName", " , userPassword", username)
			if !exist {
				w.WriteHeader(http.StatusBadRequest)
				pages.ExecuteTemplate(w, "login.html", "Invalid Username or Password.")
				return
			}
			//if !utils.IsExist("userEmail", email) { // this if you are using email in login
			//	w.WriteHeader(http.StatusBadRequest)
			//	pages.ExecuteTemplate(w, "login.html", "Invalid Email or Password.")
			//	return
			//}
			// lest check the pass
			if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password)); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				pages.ExecuteTemplate(w, "login.html", "Invalid Username or Password.")
				return
			}

		}

		next.ServeHTTP(w, r)
	})
}
