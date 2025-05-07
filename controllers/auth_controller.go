package controllers

import (
	"hris-idn/config"
	"hris-idn/entities"
	"hris-idn/helpers"
	"hris-idn/libraries"
	"hris-idn/models"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var authModel = models.NewAuthModel()
var authValidation = libraries.NewValidation()

func Login(w http.ResponseWriter, r *http.Request) {

	template := "views/static/login/login.html"
	if r.Method == http.MethodGet {
		helpers.RenderTemplate(w, template, nil)
		return
	}

	r.ParseForm()

	data := make(map[string]interface{})

	authInput := entities.Auth{
		NIK:      r.Form.Get("nik"),
		Password: r.Form.Get("password"),
	}

	if validationErrors := authValidation.Struct(authInput); validationErrors != nil {
		data["validation"] = validationErrors
		helpers.RenderTemplate(w, template, data)
		return
	}

	// Ambil data employee
	employee, err := authModel.FindByNIK(authInput.NIK)
	if err != nil {
		data["error"] = "NIK tidak ditemukan"
		log.Println("Login gagal untuk NIK:", authInput.NIK, "Error:", err)
		helpers.RenderTemplate(w, template, data)
		return
	}

	// Bandingkan password
	if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(authInput.Password)); err != nil {
		data["error"] = "NIK atau Password salah"
		helpers.RenderTemplate(w, template, data)
		return
	}

	session, _ := config.Store.Get(r, config.SESSION_ID)
	session.Values["loggedIn"] = true
	session.Values["employee"] = map[string]interface{}{
		"name":    employee.Name,
		"nik":     employee.NIK,
		"email":   employee.Email,
		"phone":   employee.Phone,
		"isAdmin": employee.IsAdmin,
		"photo":   employee.Photo,
	}
	session.Save(r, w)

	if employee.IsAdmin {
		http.Redirect(w, r, "/home-admin", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

}
