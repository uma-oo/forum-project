package auth

import (
	"net/http"

	"forum/internal/database"
)

func IsCookieSet(r *http.Request, cookieName string) bool {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return false
		}
		return false
	}
	if cookie.Value == "" {
		return false
	}
	// lets extract the token value from the cookie and compare it with the one we have in databse

	var tokenExist bool
	// lets extract the token from users table
	// be care full with  no token
	tokenErr := database.Database.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE token = $1)", cookie.Value).Scan(&tokenExist)
	if tokenErr != nil || !tokenExist {
		return false
	}

	return true
}
