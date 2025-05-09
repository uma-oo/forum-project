// internal/utils/sweet.go
package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"forum/internal/database"
	"forum/pkg/logger"
)

func IsValidUsername(username string) bool {
	// Example: Allow only alphanumeric characters and underscores
	match, _ := regexp.MatchString("^[a-zA-Z0-9_]{3,50}$", username) // we can add length like {3,15} and remove the +
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
	if len(password) <= 8 || len(password) >= 100 {
		return false
	}

	hasLower := false
	hasUpper := false
	hasDigit := false

	// Loop through the password to check for lowercase letters and digits
	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLower = true
		}
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		}
		if char >= '0' && char <= '9' {
			hasDigit = true
		}
	}

	// Password is strong if it contains at least one lowercase letter and one digit
	return hasLower && hasDigit && hasUpper
}

func IsExist(table, collumn0, collumn1, value string) (string, bool) {
	db, err := database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
	}
	// Check if the field exists in database sqlite3 in users table
	var user, pass string
	err = db.QueryRow("SELECT "+collumn0+collumn1+" FROM "+table+" WHERE  "+collumn0+"  = ?", value).Scan(&user, &pass)
	if err != nil {
		// logger.LogWithDetails(err)
		if err == sql.ErrNoRows {
			return pass, false
		}
	}
	return pass, true
}

func IsIdExist(table, column, value string) bool {
	fmt.Printf("value: %v\n", value)
	db, err := database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
	}
	defer db.Close()
	// Check if the field exists in database sqlite3 in users table
	exists := false
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE %s = ?)", table, column)
	err = db.QueryRow(query, value).Scan(&exists)
	if err != nil {
		logger.LogWithDetails(err)
	}
	fmt.Printf("exists: %v\n", exists)
	return exists
}

func IsCookieSet(r *http.Request, cookieName string) bool {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}
	if cookie.Value == "" {
		return false
	}
	db, err := database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
	}

	// lets extract the token value from the cookie and compare it with the one we have in databse
	var tokenExist bool
	// lets extract the token from users table
	// be care full with  no token
	tokenErr := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE token = $1)", cookie.Value).Scan(&tokenExist)
	if tokenErr != nil || !tokenExist {
		return false
	}

	return true
}
