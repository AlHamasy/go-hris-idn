package controllers

import (
	"fmt"
	"hris-idn/config"
	"hris-idn/entities"
	"hris-idn/helpers"
	"hris-idn/libraries"
	"hris-idn/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var attendanceModel = models.NewAttendanceModel()
var attendanceValidation = libraries.NewValidation()

func SubmitAttendance(w http.ResponseWriter, r *http.Request) {

	template := "views/static/attendance/attendance-submit.html"

	// ambil session
	session, _ := config.Store.Get(r, config.SESSION_ID)
	sessionNIK := session.Values["nik"].(string)

	data := make(map[string]interface{})

	// ambil status terakhir
	updateAttendanceStatus(sessionNIK, data)

	// tampilkan list kantor
	office, _ := models.NewOfficeModel().FindAllOffice()
	data["office"] = office

	// tampilkan list shift
	shift, _ := models.NewShiftModel().FindAllShift()
	data["shift"] = shift

	// Hitung 5 bulan terakhir
	currentDate := time.Now()
	var months []string
	for i := 0; i < 5; i++ {
		previousMonth := currentDate.AddDate(0, -i, 0)
		months = append(months, previousMonth.Format("January 2006"))
	}
	data["months"] = months

	// Get selected month from query parameter or use current month
	selectedMonth := r.URL.Query().Get("month")
	if selectedMonth == "" {
		selectedMonth = currentDate.Format("January 2006")
	}

	getAttendanceList(sessionNIK, data, selectedMonth)

	if r.Method == http.MethodGet {
		helpers.RenderTemplate(w, template, data)
		return
	}

	r.ParseForm()

	switch data["status"] {
	case helpers.NOT_CHEKCED_IN:
		actionCheckIn(w, r, sessionNIK, data, template, selectedMonth)
	case helpers.CHECKED_IN:
		actionCheckOut(w, r, sessionNIK, data, template, selectedMonth)
	}
}

// Function untuk update status attendance
func updateAttendanceStatus(sessionNIK string, data map[string]interface{}) {
	lastAtt := attendanceModel.GetLastAttendance(sessionNIK)
	data["status"] = lastAtt
}

func getAttendanceList(sessionNIK string, data map[string]interface{}, selectedMonth string) {
	attendanceList, err := attendanceModel.GetAttendanceList(sessionNIK, selectedMonth)
	if err != nil {
		log.Println("Error getting attendance list:", err)
		data["errorList"] = "Failed to get attendance list"
	}

	data["attendances"] = attendanceList
	data["selectedMonth"] = selectedMonth
}

func actionCheckIn(w http.ResponseWriter, r *http.Request, sessionNIK string, data map[string]interface{}, template string, selectedMonth string) {

	layoutTime := "2006-01-02 15:04:05"

	officeIDStr := r.FormValue("office_id")
	shiftIDStr := r.FormValue("shift_id")
	latLongStr := r.Form.Get("latlong")
	checkIn := entities.CheckIn{
		NIK:         sessionNIK,
		Photo:       r.Form.Get("attendance_photo"),
		LatLongStr:  latLongStr,
		OfficeIDStr: officeIDStr,
		ShiftIDStr:  shiftIDStr,
		Notes:       r.Form.Get("notes"),
	}

	errorMessages := attendanceValidation.Struct(checkIn)

	if errorMessages != nil {
		data["validation"] = errorMessages
		data["attendance"] = checkIn
		helpers.RenderTemplate(w, template, data)
		return
	}

	officeID, _ := strconv.ParseInt(officeIDStr, 10, 64)
	shiftID, _ := strconv.ParseInt(shiftIDStr, 10, 64)

	parts := strings.Split(latLongStr, ",")
	latitude, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	longitude, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)

	findOffice, _ := models.NewOfficeModel().FindOfficeByID(officeID)
	findShift, _ := models.NewShiftModel().FindShiftByID(shiftID)

	distance := helpers.CalculateDistance(latitude, longitude, findOffice.Latitude, findOffice.Longitude)
	if distance > float64(findOffice.Radius) {
		data["error"] = "Anda berada di luar radius kantor, tidak bisa check-in"
		helpers.RenderTemplate(w, template, data)
		return
	}

	//now := helpers.GetCurrentTimeWIB()
	now := time.Now()
	dateToday := now.Format("2006-01-02")
	shiftStartFull := fmt.Sprintf("%s %s", dateToday, findShift.StartTime)
	shiftEndFull := fmt.Sprintf("%s %s", dateToday, findShift.EndTime)

	// Gunakan ParseInLocation dengan timezone Asia/Jakarta
	loc, _ := time.LoadLocation("Asia/Jakarta")
	shiftStartTime, _ := time.ParseInLocation(layoutTime, shiftStartFull, loc)
	shiftEndTime, _ := time.ParseInLocation(layoutTime, shiftEndFull, loc)

	if now.After(shiftEndTime) {
		data["error"] = "Shift yang anda pilih sudah selesai"
		helpers.RenderTemplate(w, template, data)
		return
	}

	isLate := now.After(shiftStartTime)

	attendance := entities.CheckIn{
		NIK:       sessionNIK,
		OfficeID:  officeID,
		ShiftID:   shiftID,
		Time:      now,
		Latitude:  latitude,
		Longitude: longitude,
		IsLate:    isLate,
		Notes:     r.Form.Get("notes"),
		Photo:     r.Form.Get("attendance_photo"),
	}

	errCheckIn := attendanceModel.CheckIn(attendance)
	if errCheckIn != nil {
		data["error"] = "Error " + errCheckIn.Error()
	} else {
		if isLate {
			data["isLate"] = "Berhasil check in, namun anda terlambat"
		} else {
			data["success"] = "Berhasil check in"
		}
		updateAttendanceStatus(sessionNIK, data)
		getAttendanceList(sessionNIK, data, selectedMonth)
	}

	helpers.RenderTemplate(w, template, data)
}

func actionCheckOut(w http.ResponseWriter, r *http.Request, sessionNIK string, data map[string]interface{}, template string, selectedMonth string) {

	layoutTime := "2006-01-02 15:04:05"

	latLongStr := r.Form.Get("latlong")
	checkIn := entities.CheckOut{
		NIK:        sessionNIK,
		Photo:      r.Form.Get("attendance_photo"),
		LatLongStr: latLongStr,
		Notes:      r.Form.Get("notes"),
	}

	errorMessages := attendanceValidation.Struct(checkIn)

	if errorMessages != nil {
		data["validation"] = errorMessages
		data["attendance"] = checkIn
		helpers.RenderTemplate(w, template, data)
		return
	}

	parts := strings.Split(latLongStr, ",")
	latitude, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	longitude, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)

	officeID, shiftID, _ := models.NewAttendanceModel().GetLatestOfficeAndShift(sessionNIK)

	findOffice, _ := models.NewOfficeModel().FindOfficeByID(officeID)
	distance := helpers.CalculateDistance(latitude, longitude, findOffice.Latitude, findOffice.Longitude)
	if distance > float64(findOffice.Radius) {
		data["error"] = "Anda berada di luar radius kantor, tidak bisa check-out"
		helpers.RenderTemplate(w, template, data)
		return
	}

	findShift, _ := models.NewShiftModel().FindShiftByID(shiftID)

	//now := helpers.GetCurrentTimeWIB()
	now := time.Now()
	dateToday := now.Format("2006-01-02")
	shiftEndFull := fmt.Sprintf("%s %s", dateToday, findShift.EndTime)

	loc, _ := time.LoadLocation("Asia/Jakarta")
	shiftEndTime, _ := time.ParseInLocation(layoutTime, shiftEndFull, loc)

	isEarly := now.Before(shiftEndTime)

	attendance := entities.CheckOut{
		Time:      now,
		Latitude:  latitude,
		Longitude: longitude,
		IsEarly:   isEarly,
		Notes:     r.Form.Get("notes"),
		Photo:     r.Form.Get("attendance_photo"),
	}

	errCheckIn := attendanceModel.CheckOut(sessionNIK, attendance)
	if errCheckIn != nil {
		data["error"] = "Error " + errCheckIn.Error()
	} else {
		if isEarly {
			data["isEarly"] = "Berhasil check out, namun anda pulang lebih awal"
		} else {
			data["success"] = "Berhasil check out"
		}
		updateAttendanceStatus(sessionNIK, data)
		getAttendanceList(sessionNIK, data, selectedMonth)
	}

	helpers.RenderTemplate(w, template, data)
}

func ListAttendance(w http.ResponseWriter, r *http.Request) {

	template := "views/static/attendance/attendance-list.html"

	data := make(map[string]interface{})

	// Hitung 5 bulan terakhir
	currentDate := time.Now()
	var months []string
	for i := 0; i < 5; i++ {
		previousMonth := currentDate.AddDate(0, -i, 0)
		months = append(months, previousMonth.Format("January 2006"))
	}
	data["months"] = months

	// Get selected month from query parameter or use current month
	selectedMonth := r.URL.Query().Get("month")
	if selectedMonth == "" {
		selectedMonth = currentDate.Format("January 2006")
	}

	// tampilkan list kehadiran
	attendanceList, err := attendanceModel.GetAttendanceList("", selectedMonth)
	if err != nil {
		log.Println("Error getting attendance list:", err)
		data["errorList"] = "Failed to get attendance list"
	}

	data["attendances"] = attendanceList
	data["selectedMonth"] = selectedMonth

	if r.Method == http.MethodGet {
		helpers.RenderTemplate(w, template, data)
		return
	}
}
