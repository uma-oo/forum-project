// internal/utils/sweet.go
package utils

import (
	"database/sql"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"forum/internal/database"
)

func IsValidUsername(username string) bool {
	// Example: Allow only alphanumeric characters and underscores
	match, _ := regexp.MatchString("^[a-zA-Z0-9_]{3,15}$", username) // we can add length like {3,15} and remove the +
	return match && !isReservedUsername(username)
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Function to check if a username is reserved
func isReservedUsername(username string) bool {
	reservedWords := []string{"admin", "root", "system", "superuser"}
	for _, reserved := range reservedWords {
		if strings.ToLower(username) == reserved {
			return true
		}
	}
	return false
}

func IsStrongPassword(password string) bool {
	// Ensure password is at least 6 characters long
	if len(password) < 8 {
		return false
	}

	hasLower := false
	hasDigit := false

	// Loop through the password to check for lowercase letters and digits
	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLower = true
		}
		if char >= '0' && char <= '9' {
			hasDigit = true
		}
	}

	// Password is strong if it contains at least one lowercase letter and one digit
	return hasLower && hasDigit
}

func IsExist(collumn0, collumn1, value string) (string, bool) {
	// Check if the field exists in database sqlite3 in users table
	db := database.Database
	var user, pass string
	err := db.QueryRow("SELECT "+collumn0+collumn1+" FROM users WHERE  "+collumn0+"  = ?", value).Scan(&user, &pass)
	if err != nil {
		// logger.LogWithDetails(err)
		if err == sql.ErrNoRows {
			return pass, false
		}
	}

	return pass, true
}

func IsCookieSet(r *http.Request, cookieName string) bool {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
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
