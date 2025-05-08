package controllers

import (
	"hris-idn/config"
	"hris-idn/entities"
	"hris-idn/helpers"
	"hris-idn/libraries"
	"hris-idn/models"
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
	// redirectTo := r.URL.Query().Get("redirect") // ⬅️ ambil nilai redirect

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
	session.Values["name"] = employee.Name
	session.Values["nik"] = employee.NIK
	session.Values["email"] = employee.Email
	session.Values["phone"] = employee.Phone
	session.Values["isAdmin"] = employee.IsAdmin
	session.Save(r, w)

	// // Jika ada redirect, pakai itu
	// if redirectTo != "" {
	// 	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	// 	return
	// }

	if employee.IsAdmin {
		http.Redirect(w, r, "/home-admin", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {

	session, _ := config.Store.Get(r, config.SESSION_ID)
	// delete session
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
