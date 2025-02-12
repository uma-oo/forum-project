package auth

import (
	"log"
	"net/http"
	"reflect"

	"forum/internal/database"
	"forum/internal/models"
	"forum/internal/utils"
	"forum/pkg/logger"

	"github.com/google/uuid"
)

var (
	FormErrors = models.FormErrors{}
	FormsData  = models.FormsData{}
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RenderTemplate(w, "error.html", models.MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	db, err := database.NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	UserName := r.FormValue("userName")
	Token := uuid.New().String()
	stm, err := db.Prepare("UPDATE users SET token = ? where userName = ?")
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	_, err = stm.Exec(Token, UserName)
	if err != nil {
		logger.LogWithDetails(err)
		utils.RenderTemplate(w, "error.html", models.InternalServerError, http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:   "token",
		Value:  Token,
		MaxAge: 3600,
		Path:   "/",
	}
	http.SetCookie(w, cookie)
	log.Println(UserName, "logged in")
	http.Redirect(w, r, "/", http.StatusFound)
}

func IsValidFormValues(FormErrors models.FormErrors) (*models.Data, bool) {
	values := reflect.ValueOf(FormErrors)
	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).String() != "" {
			data := &models.Data{
				FormsData: FormsData,
			}
			data.FormsData.FormErrors = FormErrors
			return data, false
		}
	}
	return &models.Data{}, true
}
