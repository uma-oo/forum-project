package internal

import (
	"log"
	"text/template"
)

var Templates *template.Template

func ParseTemplates() {
	var err error
	Templates, err = template.ParseGlob("./web/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	Templates, err = Templates.ParseGlob("./web/components/*.html")
	if err != nil {
		log.Fatal(err)
	}
}
