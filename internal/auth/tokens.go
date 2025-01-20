package auth

import (
	"fmt"
	"net/http"
)

func IsCookieSet(r *http.Request, cookieName string) bool {
	cookie, err := r.Cookie(cookieName)
	if err != nil {

		if err == http.ErrNoCookie {
			return false
		}

		fmt.Println("Error retrieving cookie:", err)
		return false
	}

	return cookie.Value != ""
}
