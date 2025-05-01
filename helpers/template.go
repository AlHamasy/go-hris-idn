package helpers

import (
	"html/template"
	"net/http"
)

// Fungsi global untuk render template
func RenderTemplate(w http.ResponseWriter, path string, data interface{}) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
