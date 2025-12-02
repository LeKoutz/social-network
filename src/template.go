package forum

import (
	"html/template"
	"log"
	"net/http"
)

func respondView(res http.ResponseWriter, view string, data ResponseStruct) {
	var templatesDir string = "templates"
	var tmpl *template.Template
	tmpl, err := template.ParseGlob(templatesDir + "/*.html")
	if err != nil {
		log.Printf("Error: %s", err.Error())
		respondError(http.StatusInternalServerError, res, "Not implemented")
		return
	}
	err = tmpl.ExecuteTemplate(res, view, data)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		respondError(http.StatusInternalServerError, res, "Not implemented")
		return
	}
}
