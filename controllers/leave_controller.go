package controllers

import (
	"hris-idn/config"
	"hris-idn/entities"
	"hris-idn/helpers"
	"hris-idn/libraries"
	"hris-idn/models"
	"log"
	"net/http"
	"strings"
)

var leaveModel = models.NewLeaveModel()
var leaveValidation = libraries.NewValidation()

func SubmitLeave(w http.ResponseWriter, r *http.Request) {

	template := "views/static/leave/leave-submit.html"

	session, _ := config.Store.Get(r, config.SESSION_ID)
	sessionNIK := session.Values["nik"].(string)

	data := make(map[string]interface{})

	leaveType, _ := leaveModel.FindAllLeaveType()
	data["leaveType"] = leaveType

	leaveList, err := leaveModel.GetLeaveList(sessionNIK, "")
	if err != nil {
		log.Println("Error getting leave list:", err)
		data["errorList"] = "Failed to get leave list"
	}
	data["leaves"] = leaveList

	if r.Method == http.MethodGet {
		helpers.RenderTemplate(w, template, data)
		return
	}

	r.ParseForm()

	rawLeaveDates := r.Form["leave_date[]"]
	leaveDate := cleanLeaveDates(rawLeaveDates)
	submitLeave := entities.SubmitLeave{
		NIK:           sessionNIK,
		LeaveTypeId:   r.Form.Get("leave_type_id"),
		LeaveDate:     leaveDate,
		LeaveDateJoin: strings.Join(leaveDate, ","),
		Reason:        r.Form.Get("reason"),
		Attachment:    r.Form.Get("attachment_photo"),
		Status:        helpers.PENDING_LEAVE,
	}

	validationErrors := leaveValidation.Struct(submitLeave)

	if validationErrors != nil {
		data["validation"] = validationErrors
		data["leave"] = submitLeave
		helpers.RenderTemplate(w, template, data)
		return
	}

	errSubmit := leaveModel.InsertLeave(submitLeave)
	if (errSubmit) != nil {
		data["error"] = "Error " + errSubmit.Error()
	} else {
		data["success"] = "Pengajuan cuti berhasil, silahkan tunggu persetujuan dari Admin"
	}

	helpers.RenderTemplate(w, template, data)
}

func cleanLeaveDates(input []string) []string {
	var cleaned []string
	for _, date := range input {
		trimmed := strings.TrimSpace(date)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

func ListLeave(w http.ResponseWriter, r *http.Request) {

	template := "views/static/leave/leave-list.html"

	helpers.RenderTemplate(w, template, nil)
}
