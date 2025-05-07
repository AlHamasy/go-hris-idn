package controllers

import (
	"hris-idn/helpers"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {

	template := "views/static/home/home.html"
	helpers.RenderTemplate(w, template, nil)
}

func HomeAdmin(w http.ResponseWriter, r *http.Request) {

	template := "views/static/home/home-admin.html"
	helpers.RenderTemplate(w, template, nil)
}
