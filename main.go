package main

import (
	"fmt"
	"hris-idn/controllers"
	"net/http"
)

func main() {
	// baca file untuk resource yg ada di static (css, js, img)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))

	// home
	http.HandleFunc("/home", controllers.Home)
	http.HandleFunc("/home-admin", controllers.HomeAdmin)

	// news
	http.HandleFunc("/news", controllers.News)

	// attendance
	http.HandleFunc("/attendance-submit", controllers.SubmitAttendance)
	http.HandleFunc("/attendance-list", controllers.ListAttendance)

	// leave
	http.HandleFunc("/leave-submit", controllers.SubmitLeave)
	http.HandleFunc("/leave-list", controllers.ListLeave)

	// payroll
	http.HandleFunc("/payroll-self", controllers.MyPayroll)
	http.HandleFunc("/payroll-list", controllers.ListPayroll)

	// employee
	http.HandleFunc("/employee", controllers.Employee)
	http.HandleFunc("/employee/add-employee", controllers.AddEmployee)
	http.HandleFunc("/employee/edit-employee", controllers.EditEmployee)
	http.HandleFunc("/employee/delete-employee", controllers.DeleteEmployee)

	// office
	http.HandleFunc("/office", controllers.Office)
	http.HandleFunc("/office/add-office", controllers.AddOffice)
	http.HandleFunc("/office/edit-office", controllers.EditOffice)
	http.HandleFunc("/office/delete-office", controllers.DeleteOffice)

	//shift
	http.HandleFunc("/shift", controllers.Shift)
	http.HandleFunc("/shift/add-shift", controllers.AddShift)
	http.HandleFunc("/shift/edit-shift", controllers.EditShift)
	http.HandleFunc("/shift/delete-shift", controllers.DeleteShift)

	fmt.Println("Server running in : http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
