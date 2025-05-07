package controllers

import (
	"hris-idn/helpers"
	"net/http"
)

func MyPayroll(w http.ResponseWriter, r *http.Request) {

	template := "views/static/payroll/payroll-self.html"
	helpers.RenderTemplate(w, template, nil)
}

func ListPayroll(w http.ResponseWriter, r *http.Request) {

	template := "views/static/payroll/payroll-list.html"
	helpers.RenderTemplate(w, template, nil)
}
