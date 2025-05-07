package controllers

import (
	"hris-idn/helpers"
	"net/http"
)

func SubmitLeave(w http.ResponseWriter, r *http.Request) {

	template := "views/static/leave/leave-submit.html"
	helpers.RenderTemplate(w, template, nil)
}

func ListLeave(w http.ResponseWriter, r *http.Request) {

	template := "views/static/leave/leave-list.html"
	helpers.RenderTemplate(w, template, nil)
}
