package internal

import (
	"bytes"
	"log"
	"text/template"
)

type Pages struct {
	All_Templates *template.Template
	Buf           bytes.Buffer
}

var Pagess Pages

func ParseTemplates() {
	var err error
	Pagess.All_Templates, err = template.ParseGlob("./web/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	Pagess.All_Templates, err = Pagess.All_Templates.ParseGlob("./web/components/*.html")
	if err != nil {
		log.Fatal(err)
	}
}
