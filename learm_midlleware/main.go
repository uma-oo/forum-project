package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Logger struct {
	LogFile *os.File
}

// lets craete our servmux
type Mine_ServMux struct {
	Routers map[string]http.Handler
}

// lets craete our handlefunc
func (mux *Mine_ServMux) Mine_Handlfunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if mux.Routers == nil {
		mux.Routers = make(map[string]http.Handler)
	}
	mux.Routers[pattern] = http.HandlerFunc(handler)
}

// lets create our serveHTTP
func (mux *Mine_ServMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// lets range over our routs
	if handler, ok := mux.Routers[r.URL.Path]; ok {
		fmt.Println(r.URL.Path)
		handler.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}

func init() {
	var loger Logger
	var err error
	loger.LogFile, err = os.OpenFile("server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(loger.LogFile)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from server\n")
	w.Write([]byte("this is a public site\n"))
}

func secure_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from server\n")
	w.Write([]byte("this privete site\n"))
	w.Write([]byte("you are succefully loged\n"))
}

func LoginMidlleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// lets log the request
		log.Printf("%v recieved a request >>   Method: %s, URL: %s, Client IP: %s", r.Header.Get("Date"), r.Method, r.URL.Path, r.RemoteAddr)
		log.Print("------------------------")
		// now lets check for authorisation
		authorisationHrader := r.Header.Get("Authorization")
		if authorisationHrader == "" || !strings.HasPrefix(authorisationHrader, "aghyor ") || !validToken(authorisationHrader[7:]) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// now call the next hamdel in the chain
		next.ServeHTTP(w, r)
	})
}

func validToken(token string) bool {
	return token == "boda"
}

func main() {
	var loger Logger
	defer loger.LogFile.Close()
	// Create a new HTTP multiplexer (router)
	mux := &Mine_ServMux{}
	// mux := http.DefaultServeMux
	mux.Mine_Handlfunc("/", handler)
	// mux.HandleFunc("/", handler)
	// mux.Handle("/secure", LoginMidlleware(http.HandlerFunc(secure_handler)))
	fmt.Println("server is running on port 8080 ... http://localhost:8080")
	// Start the server with the logging middleware
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
