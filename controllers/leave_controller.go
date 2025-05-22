package controllers

import (
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

var leaveModel = models.NewLeaveModel()
var leaveValidation = libraries.NewValidation()

func SubmitLeave(w http.ResponseWriter, r *http.Request) {

	template := "views/static/leave/leave-submit.html"

	session, _ := config.Store.Get(r, config.SESSION_ID)
	sessionNIK := session.Values["nik"].(string)

	data := make(map[string]interface{})

	leaveType, _ := leaveModel.FindAllLeaveType()
	data["leaveType"] = leaveType

	getLeaveList(data, sessionNIK)

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

	isValid, errMsg := helpers.IsLeaveDateValid(leaveDate)
	if !isValid {
		data["error"] = errMsg
		data["leave"] = submitLeave
		helpers.RenderTemplate(w, template, data)
		return
	}

	errSubmit := leaveModel.InsertLeave(submitLeave)
	if (errSubmit) != nil {
		data["error"] = "Error " + errSubmit.Error()
	} else {
		data["success"] = "Pengajuan cuti berhasil, silahkan tunggu persetujuan dari Admin"
		getLeaveList(data, sessionNIK)
	}

	helpers.RenderTemplate(w, template, data)
}

func getLeaveList(data map[string]interface{}, sessionNIK string) {
	leaveList, err := leaveModel.GetLeaveList(sessionNIK, "")
	if err != nil {
		log.Println("Error getting leave list:", err)
		data["errorList"] = "Failed to get leave list"
	}
	data["leaves"] = leaveList
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
	leaveList, err := leaveModel.GetLeaveList("", selectedMonth)
	if err != nil {
		log.Println("Error getting leave list:", err)
		data["errorList"] = "Failed to get leave list"
	}

	data["leaves"] = leaveList
	data["selectedMonth"] = selectedMonth

	if r.Method == http.MethodGet {
		helpers.RenderTemplate(w, template, data)
		return
	}
}

func ApprovalLeave(w http.ResponseWriter, r *http.Request) {

	template := "views/static/leave/leave-approval.html"

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var data = make(map[string]interface{})

	getLeave(id, w, data)

	if r.Method == http.MethodGet {
		helpers.RenderTemplate(w, template, data)
		return
	}

	r.ParseForm()

	status := r.Form.Get("status")
	reason := r.Form.Get("reason_status")

	idInt64, _ := strconv.ParseInt(idStr, 10, 64)
	statusInt64, _ := strconv.ParseInt(status, 10, 64)
	approval := entities.ApprovalLeave{
		Id:           idInt64,
		Status:       statusInt64,
		ReasonStatus: reason,
		UpdatedAt:    time.Now(),
	}

	errorValidation := leaveValidation.Struct(approval)

	if errorValidation != nil {
		data["validation"] = errorValidation
		data["approval"] = approval
		helpers.RenderTemplate(w, template, data)
		return
	}

	errApprove := leaveModel.UpdateLeaveStatus(approval)
	if errApprove != nil {
		data["error"] = "Gagal memproses cuti " + errApprove.Error()
	} else {
		data["success"] = "Cuti berhasil diproses"
		getLeave(id, w, data)
	}

	helpers.RenderTemplate(w, template, data)
}

func getLeave(id int64, w http.ResponseWriter, data map[string]interface{}) {
	leave, err := leaveModel.GetLeaveById(id)
	if err != nil || leave == nil {
		http.Error(w, "Data cuti tidak ditemukan", http.StatusBadRequest)
		return
	}
	data["leave"] = leave
}
