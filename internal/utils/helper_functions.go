package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"forum/internal"
	"forum/internal/database"
	"forum/internal/models"
	"forum/pkg/logger"
)

func RenderTemplate(w http.ResponseWriter, tmp string, data interface{}, status int) {
	var buf bytes.Buffer
	err := internal.Templates.ExecuteTemplate(&buf, tmp, data)
	if err != nil {
		logger.LogWithDetails(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>Internal Server Error 500</h1>")
		return
	}
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func UserData(r *http.Request, cookieName string, currentPath string) (*models.User, error) {
	var user models.User
	// if !IsCookieSet(r, cookieName) {
	// 	return &models.User{
	// 		CurrentPath: currentPath,
	// 		IsLoged:     false,
	// 		UserName:    "",
	// 		UserEmail:   "",
	// 		UserId:      "",
	// 	}, nil
	// } else {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	db, err := database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}

	query := `SELECT id, userEmail , userName FROM users where token= ? `
	statement, err := db.Prepare(query)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	row := statement.QueryRow(cookie.Value)
	row.Scan(&user.UserId, &user.UserEmail, &user.UserName)

	fmt.Printf("user.UserId: %v\n", user.UserId)
	fmt.Printf("user.UserName: %v\n", user.UserName)
	fmt.Printf("user.IsLoged: %v\n", user.IsLoged)
	return &models.User{
		IsLoged:     true,
		CurrentPath: currentPath,
		UserName:    user.UserName,
		UserEmail:   user.UserEmail,
		UserId:      user.UserId,
	}, nil
}

func CheckTokenExpired(user *models.User) (bool, error) {
	var exp_date string
	db, err := database.NewDatabase()
	if err != nil {
		return true, err
	}
	query := `SELECT expiration_date FROM users WHERE id = ?`
	statement, err := db.Prepare(query)
	if err != nil {
		return true, err
	}
	err = statement.QueryRow(user.UserId).Scan(&exp_date)
	if err != nil {
		return true, err
	}
	if time.Now().Format(time.RFC3339) > exp_date {
		return true, nil
	}
	fmt.Println(exp_date)
	return false, nil
}
