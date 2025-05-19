package controllers

import (
	"hris-idn/config"
	"hris-idn/helpers"
	"net/http"
)

// Home untuk karyawan
func Home(w http.ResponseWriter, r *http.Request) {

	template := "views/static/home/home.html"
	session, _ := config.Store.Get(r, config.SESSION_ID)
	sessionName := session.Values["name"].(string)

	data := make(map[string]interface{})
	data["name"] = sessionName

	helpers.RenderTemplate(w, template, data)
}

// Home untuk admin
func HomeAdmin(w http.ResponseWriter, r *http.Request) {

	template := "views/static/home/home-admin.html"

	session, _ := config.Store.Get(r, config.SESSION_ID)
	sessionName := session.Values["name"].(string)

	data := make(map[string]interface{})
	data["name"] = sessionName

	helpers.RenderTemplate(w, template, data)
}
