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

var shiftModel = models.NewShiftModel()
var shiftValidation = libraries.NewValidation()

func Shift(w http.ResponseWriter, r *http.Request) {

	template := "views/static/shift/shift.html"

	if r.Method == http.MethodGet {

		var data = make(map[string]interface{})
		shifts, err := shiftModel.FindAllShift()

		if err != nil {
			data["error"] = "Terdapat kesahalan saat menampilkan data shift " + err.Error()
			log.Println("error :", err.Error())
		} else {
			data["shift"] = shifts
		}

		helpers.RenderTemplate(w, template, data)
		return
	}
}

func AddShift(w http.ResponseWriter, r *http.Request) {

	template := "views/static/shift/add-shift.html"

	if r.Method == http.MethodGet {
		helpers.RenderTemplate(w, template, nil)
		return
	}

	r.ParseForm()

	var data = make(map[string]interface{})

	shift := entities.Shift{
		Name:      r.Form.Get("name"),
		StartTime: r.Form.Get("start_time"),
		EndTime:   r.Form.Get("end_time"),
	}

	errorMessages := shiftValidation.Struct(shift)

	if errorMessages != nil {
		data["validation"] = errorMessages
		data["shift"] = shift
		helpers.RenderTemplate(w, template, data)
		return
	}

	err := shiftModel.AddShift(shift)

	if err != nil {
		data["error"] = "Gagal menambahkan shift: " + err.Error()
	} else {
		data["success"] = "Berhasil menambahkan shift"
	}

	helpers.RenderTemplate(w, template, data)

}

func EditShift(w http.ResponseWriter, r *http.Request) {

	template := "views/static/shift/edit-shift.html"
	id := r.URL.Query().Get("id")
	int64Id, _ := strconv.ParseInt(id, 10, 64)

	var data = make(map[string]interface{})

	if r.Method == http.MethodGet {

		// Ambil data shift berdasarkan ID
		shift, err := shiftModel.FindShiftByID(int64Id)
		if err != nil || id == "" {
			http.Error(w, "ID tidak ditemukan", http.StatusBadRequest)
			return
		}

		// Kirim data ke template
		data["shift"] = shift
		helpers.RenderTemplate(w, template, data)
		return
	}

	r.ParseForm()

	shift := entities.Shift{
		Id:        int64Id,
		Name:      r.Form.Get("name"),
		StartTime: r.Form.Get("start_time"),
		EndTime:   r.Form.Get("end_time"),
	}

	errorMessages := shiftValidation.Struct(shift)

	if errorMessages != nil {
		data["validation"] = errorMessages
		data["shift"] = shift
		helpers.RenderTemplate(w, template, data)
		return
	}

	err := shiftModel.EditShift(shift)

	if err != nil {
		data["error"] = "Edit data gagal: " + err.Error()
	} else {
		data["success"] = "Edit data berhasil"
	}

	helpers.RenderTemplate(w, template, data)

}

func DeleteShift(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	int64Id, _ := strconv.ParseInt(id, 10, 64)

	if id == "" {
		http.Error(w, "ID tidak ditemukan", http.StatusBadRequest)
		return
	}

	err := shiftModel.SoftDeleteShift(int64Id)
	if err != nil {
		http.Error(w, "Gagal menghapus data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/shift", http.StatusSeeOther)
}
