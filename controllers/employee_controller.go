package controllers

import (
	"hris-idn/entities"
	"hris-idn/helpers"
	"hris-idn/libraries"
	"hris-idn/models"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var employeeModel = models.NewEmployeeModel()
var employeeValidation = libraries.NewValidation()

func Employee(w http.ResponseWriter, r *http.Request) {

	template := "views/static/employee/employee.html"
	if r.Method == http.MethodGet {

		var data = make(map[string]interface{})
		employees, err := employeeModel.FindAllEmployee()

		if err != nil {
			data["error"] = "Terdapat kesahalan saat menampilkan data karyawan " + err.Error()
			log.Println("error :", err.Error())
		} else {
			data["employee"] = employees
		}

		helpers.RenderTemplate(w, template, data)
		return
	}
}

func AddEmployee(w http.ResponseWriter, r *http.Request) {

	template := "views/static/employee/add-employee.html"
	if r.Method == http.MethodGet {
		helpers.RenderTemplate(w, template, nil)
		return
	}

	r.ParseForm()

	isAdmin := r.Form.Get("is-admin") != ""
	employee := entities.Employee{
		UUID:      uuid.New().String(),
		Name:      r.Form.Get("name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		Address:   r.Form.Get("address"),
		NIK:       r.Form.Get("nik"),
		Gender:    r.Form.Get("gender"),
		BirthDate: r.Form.Get("birth_date"),
		//Photo:     r.Form.Get("photo-base64"),
		IsAdmin: isAdmin,
	}

	errorMessages := employeeValidation.Struct(employee)

	var data = make(map[string]interface{})
	if errorMessages != nil {
		data["validation"] = errorMessages
		data["employee"] = employee
		helpers.RenderTemplate(w, template, data)
		return
	}

	// hash password
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(helpers.DEFAULT_PASSWORD), bcrypt.DefaultCost)
	employee.Password = string(hashPassword)

	// insert ke database
	err := employeeModel.AddEmployee(employee)

	if err != nil {
		data["error"] = "Registrasi gagal: " + err.Error()
	} else {
		data["success"] = "Registrasi berhasil, silahkan arahkan karyawan untuk login lalu ubah passwordnya"
	}

	helpers.RenderTemplate(w, template, data)
}

func EditEmployee(w http.ResponseWriter, r *http.Request) {

	template := "views/static/employee/edit-employee.html"
	uuid := r.URL.Query().Get("uuid")
	var data = make(map[string]interface{})

	if r.Method == http.MethodGet {

		// cek uuid
		if r.URL.Query().Get("uuid") == "" {
			http.Error(w, "UUID tidak ditemukan", http.StatusBadRequest)
			return
		}

		// Ambil data employee berdasarkan UUID
		employee, err := employeeModel.FindEmployeeByUUID(uuid)
		if err != nil {
			data["error"] = "Gagal mengambil data karyawan: " + err.Error()
			helpers.RenderTemplate(w, template, data)
			return
		}

		// Kirim data ke template
		data["employee"] = employee
		helpers.RenderTemplate(w, template, data)
		return
	}

	r.ParseForm()

	isAdmin := r.Form.Get("is-admin") != ""
	employee := entities.EditEmployee{
		UUID:      uuid,
		Name:      r.Form.Get("name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		Address:   r.Form.Get("address"),
		NIK:       r.Form.Get("nik"),
		Gender:    r.Form.Get("gender"),
		BirthDate: r.Form.Get("birth_date"),
		IsAdmin:   isAdmin,
	}

	errorMessages := employeeValidation.Struct(employee)

	if errorMessages != nil {
		data["validation"] = errorMessages
		data["employee"] = employee
		helpers.RenderTemplate(w, template, data)
		return
	}

	// update ke database
	err := employeeModel.EditEmployee(employee)

	if err != nil {
		data["error"] = "Edit data gagal: " + err.Error()
	} else {
		data["success"] = "Edit data berhasil"
	}

	helpers.RenderTemplate(w, template, data)
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")

	if uuid == "" {
		http.Error(w, "UUID tidak ditemukan", http.StatusBadRequest)
		return
	}

	err := employeeModel.SoftDeleteEmployee(uuid)
	if err != nil {
		http.Error(w, "Gagal menghapus data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/employee", http.StatusSeeOther)
}
