package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"forum/internal"
	"forum/internal/auth"
	"forum/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

var rateLimiter = auth.NewRateLimiter(5, time.Minute) // 5 requests per minute limit

// Custom middleware to validate the registration form
func Reg_Log_Middleware(next http.Handler) http.Handler {
	fmt.Println("inside the Reg_log_middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pages := internal.Pagess.All_Templates
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
			fmt.Println("refister condintion")
			if !utils.IsValidUsername(username) {
				auth.FormErrors.InvalidUserName = "Username is invalid."
				http.Redirect(w, r, "/register", http.StatusBadRequest)
			}
			_, exist := utils.IsExist("userName", "", username)
			if exist {
				auth.FormErrors.InvalidUserName = "Username is already taken."
				http.Redirect(w, r, "/register", http.StatusBadRequest)

			}
			if !utils.IsValidEmail(email) {
				auth.FormErrors.InvalidEmail = "Email is invalid."
				http.Redirect(w, r, "/register", http.StatusBadRequest)
				return
			}
			_, exist = utils.IsExist("userEmail", "", email)
			if exist {
				auth.FormErrors.InvalidEmail = "Email is already taken."
				http.Redirect(w, r, "/register", http.StatusBadRequest)
				return
			}
			// Validate password
			if !utils.IsStrongPassword(password) {
				auth.FormErrors.InvalidPassword = "Password is too weak."
				http.Redirect(w, r, "/register", http.StatusBadRequest)
				return
			}
		}
		// validate login form
		if r.URL.Path == "/auth/log_in" {
			fmt.Println("login condintion")
			// check if username is exist
			pass, exist := utils.IsExist("userName", " , userPassword", username)
			if !exist {
				auth.FormErrors.FormError = "Invalid Username or Password."
				http.Redirect(w, r, "/login", http.StatusBadRequest)
				return
			}

			//if !utils.IsExist("userEmail", email) { // this if you are using email in login
			//	pages.ExecuteTemplate(w, "login.html", "Invalid Email or Password.")
			//	return
			//}

			// lest check the pass
			if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password)); err != nil {
				auth.FormErrors.FormError = "Invalid Username or Password."
				http.Redirect(w, r, "/login", http.StatusUnauthorized)
				return
			}
		}
		fmt.Println("here")
		next.ServeHTTP(w, r)
	})
}
