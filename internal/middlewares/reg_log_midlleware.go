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
		auth.FormsData.UserNameInput= r.FormValue("userName")
		auth.FormsData.UserEmailInput= r.FormValue("userEmail")
		auth.FormsData.UserPasswordInput= r.FormValue("userPassword")
		// Validate the registration form
		invalidFormdata := false
		if r.URL.Path == "/auth/register" {
			fmt.Println("refister condintion")
			if !utils.IsValidUsername(auth.FormsData.UserNameInput) {
				auth.FormErrors.InvalidUserName = "Username is invalid."
				invalidFormdata = true
			}
			_, exist := utils.IsExist("userName", "", auth.FormsData.UserNameInput)
			if exist {
				auth.FormErrors.InvalidUserName = "Username is already taken."
				invalidFormdata = true
			}
			if !utils.IsValidEmail(auth.FormsData.UserEmailInput) {
				auth.FormErrors.InvalidEmail = "Email is invalid."
				invalidFormdata = true
			}
			_, exist = utils.IsExist("userEmail", "", auth.FormsData.UserEmailInput)
			if exist {
				auth.FormErrors.InvalidEmail = "Email is already taken."
				invalidFormdata = true
			}
			// Validate auth.FormsData.UserPasswordInput
			if !utils.IsStrongPassword(auth.FormsData.UserPasswordInput) {
				auth.FormErrors.InvalidPassword = "Password is too weak."
				invalidFormdata = true
			}
			if invalidFormdata {
				http.Redirect(w, r, "/register", http.StatusSeeOther)
				return
			}
		}
		// validate login form
		if r.URL.Path == "/auth/log_in" {
			fmt.Println("login condintion")
			// check if username is exist
			pass, exist := utils.IsExist("userName", " , userPassword", auth.FormsData.UserNameInput)
			if !exist {
				auth.FormErrors.FormError = "Invalid Username or Password."
				invalidFormdata = true
			}

			//if !utils.IsExist("userEmail", auth.FormsData.UserEmailInput) { // this if you are using auth.FormsData.UserEmailInput in login
			//	pages.ExecuteTemplate(w, "login.html", "Invalid Email or Password.")
			//	return
			//}

			// lest check the pass
			if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(auth.FormsData.UserPasswordInput)); err != nil {
				auth.FormErrors.FormError = "Invalid Username or Password."
			}

			if invalidFormdata {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
		fmt.Println("going next")
		next.ServeHTTP(w, r)
	})
}
