package controllers

import (
	"hris-idn/helpers"
	"net/http"
)

func News(w http.ResponseWriter, r *http.Request) {

	template := "views/static/news/news.html"
	helpers.RenderTemplate(w, template, nil)
}
