// internal/middleware/middleware.go
package middlewares

import (
	"fmt"
	"log"
	"net/http"
)

func Auth_Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// lets log the request
		log.Printf("%v recieved a request >>   Method: %s, URL: %s, Client IP: %s", r.Header.Get("Date"), r.Method, r.URL.Path, r.RemoteAddr)
		log.Print("------------------------")
		// now lets check for authorisation

		cookie, err := r.Cookie("token")
		fmt.Printf("cookie: %v\n", cookie)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		fmt.Printf("cookie: %v\n", cookie)

		// now call the next hamdel in the chain
		next.ServeHTTP(w, r)
	})
}
