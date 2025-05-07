package controllers

import (
	"hris-idn/helpers"
	"net/http"
)

func SubmitAttendance(w http.ResponseWriter, r *http.Request) {

	template := "views/static/attendance/attendance-submit.html"
	helpers.RenderTemplate(w, template, nil)
}

func ListAttendance(w http.ResponseWriter, r *http.Request) {

	template := "views/static/attendance/attendance-list.html"
	helpers.RenderTemplate(w, template, nil)
}
