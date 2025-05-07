package controllers

import (
	"hris-idn/entities"
	"hris-idn/helpers"
	"hris-idn/libraries"
	"hris-idn/models"
	"log"
	"net/http"
	"strconv"
)

var officeModel = models.NewOfficeModel()
var officeValidation = libraries.NewValidation()

func Office(w http.ResponseWriter, r *http.Request) {

	template := "views/static/office/office.html"
	if r.Method == http.MethodGet {

		var data = make(map[string]interface{})
		offices, err := officeModel.FindAllOffice()

		if err != nil {
			data["error"] = "Terdapat kesahalan saat menampilkan data kantor " + err.Error()
			log.Println("error :", err.Error())
		} else {
			data["office"] = offices
		}

		helpers.RenderTemplate(w, template, data)
		return
	}
}

func AddOffice(w http.ResponseWriter, r *http.Request) {

	template := "views/static/office/add-office.html"
	if r.Method == http.MethodGet {
		helpers.RenderTemplate(w, template, nil)
		return
	}

	r.ParseForm()

	var data = make(map[string]interface{})

	// Parsing radius
	radiusStr := r.Form.Get("radius")
	radius, _ := strconv.ParseInt(radiusStr, 10, 64)

	// Parsing latitude
	latitudeStr := r.Form.Get("latitude")
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)

	// Parsing longitude
	longitudeStr := r.Form.Get("longitude")
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)

	office := entities.Office{
		Name:      r.Form.Get("name"),
		Address:   r.Form.Get("address"),
		Latitude:  latitude,
		Longitude: longitude,
		Radius:    radius,
	}

	errorMessages := officeValidation.Struct(office)

	if errorMessages != nil {
		data["validation"] = errorMessages
		data["office"] = office
		helpers.RenderTemplate(w, template, data)
		return
	}

	err := officeModel.AddOffice(office)

	if err != nil {
		data["error"] = "Gagal menambahkan kantor: " + err.Error()
	} else {
		data["success"] = "Berhasil menambahkan kantor"
	}

	helpers.RenderTemplate(w, template, data)
}

func EditOffice(w http.ResponseWriter, r *http.Request) {

	template := "views/static/office/edit-office.html"
	id := r.URL.Query().Get("id")
	int64Id, _ := strconv.ParseInt(id, 10, 64)

	var data = make(map[string]interface{})

	if r.Method == http.MethodGet {

		// Ambil data office berdasarkan ID
		office, err := officeModel.FindOfficeByID(int64Id)
		if err != nil || id == "" {
			http.Error(w, "ID tidak ditemukan", http.StatusBadRequest)
			return
		}

		// Kirim data ke template
		data["office"] = office
		helpers.RenderTemplate(w, template, data)
		return
	}

	r.ParseForm()

	// Parsing radius
	radiusStr := r.Form.Get("radius")
	radius, _ := strconv.ParseInt(radiusStr, 10, 64)

	// Parsing latitude
	latitudeStr := r.Form.Get("latitude")
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)

	// Parsing longitude
	longitudeStr := r.Form.Get("longitude")
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)

	office := entities.Office{
		Id:        int64Id,
		Name:      r.Form.Get("name"),
		Address:   r.Form.Get("address"),
		Latitude:  latitude,
		Longitude: longitude,
		Radius:    radius,
	}

	errorMessages := officeValidation.Struct(office)

	if errorMessages != nil {
		data["validation"] = errorMessages
		data["office"] = office
		helpers.RenderTemplate(w, template, data)
		return
	}

	err := officeModel.EditOffice(office)

	if err != nil {
		data["error"] = "Edit data gagal: " + err.Error()
	} else {
		data["success"] = "Edit data berhasil"
	}

	helpers.RenderTemplate(w, template, data)

}

func DeleteOffice(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	int64Id, _ := strconv.ParseInt(id, 10, 64)

	if id == "" {
		http.Error(w, "ID tidak ditemukan", http.StatusBadRequest)
		return
	}

	err := officeModel.SoftDeleteOffice(int64Id)
	if err != nil {
		http.Error(w, "Gagal menghapus data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/office", http.StatusSeeOther)
}
