package utils

import (
	"bytes"
	"fmt"
	"net/http"

	"forum/internal"
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
