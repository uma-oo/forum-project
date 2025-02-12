package auth

import (
	"log"
	"net/http"
	"time"

	"forum/internal/database"
	"forum/internal/models"
	"forum/internal/utils"
	"forum/pkg/logger"
)

func LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// lets check in first that is already have a session
	if utils.IsCookieSet(r, "token") {
		http.SetCookie(w, &http.Cookie{
			Name:   "token", // name of the cookie
			Value:  "",      // clear the cookie value
			MaxAge: -1,      // set expiration time to a time in the past
			Path:   "/",     // scope of the cookie
		})
		// db, err := database.NewDatabase()
		// log.Print("A User logged out")
		// stm, err := db.Prepare("UPDATE users SET token = ? , token_created_at = ? , expiration_date = ? where userName = ?")
		// if err != nil {
		// 	logger.LogWithDetails(err)
		// 	utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		// 	return
		// }
		// _, err = stm.Exec("", time.Now(), time.Now().Add(60*time.Minute), "")
		// if err != nil {
		// 	logger.LogWithDetails(err)
		// 	utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		// 	return
		// }

	} else {
		utils.RenderTemplate(w, "error.html", models.PageNotFound, http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
