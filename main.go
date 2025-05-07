package main

import (
	"fmt"
	"hris-idn/controllers"
	"net/http"
)

func main() {
	// baca file untuk resource yg ada di static (css, js, img)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))

	// login
	http.HandleFunc("/login", controllers.Login) // diakses oleh admin & non-admin

	// home
	http.HandleFunc("/home", controllers.Home)            // diakses oleh non-admin
	http.HandleFunc("/home-admin", controllers.HomeAdmin) // diakses oleh admin

	// news
	http.HandleFunc("/news", controllers.News) // diakses oleh admin

	// attendance
	http.HandleFunc("/attendance-submit", controllers.SubmitAttendance) // diakses oleh non admin
	http.HandleFunc("/attendance-list", controllers.ListAttendance)     // diakses oleh admin

	// leave
	http.HandleFunc("/leave-submit", controllers.SubmitLeave) // diakses oleh non admin
	http.HandleFunc("/leave-list", controllers.ListLeave)     // diakses oleh admin

	// payroll
	http.HandleFunc("/payroll-self", controllers.MyPayroll)   // diakses oleh non admin
	http.HandleFunc("/payroll-list", controllers.ListPayroll) // diakses oleh admin

	// employee
	http.HandleFunc("/employee", controllers.Employee)                       // diakses oleh admin
	http.HandleFunc("/employee/add-employee", controllers.AddEmployee)       // diakses oleh admin
	http.HandleFunc("/employee/edit-employee", controllers.EditEmployee)     // diakses oleh admin
	http.HandleFunc("/employee/delete-employee", controllers.DeleteEmployee) // diakses oleh admin

	// office
	http.HandleFunc("/office", controllers.Office)                     // diakses oleh admin
	http.HandleFunc("/office/add-office", controllers.AddOffice)       // diakses oleh admin
	http.HandleFunc("/office/edit-office", controllers.EditOffice)     // diakses oleh admin
	http.HandleFunc("/office/delete-office", controllers.DeleteOffice) // diakses oleh admin

	//shift
	http.HandleFunc("/shift", controllers.Shift)                    // diakses oleh admin
	http.HandleFunc("/shift/add-shift", controllers.AddShift)       // diakses oleh admin
	http.HandleFunc("/shift/edit-shift", controllers.EditShift)     // diakses oleh admin
	http.HandleFunc("/shift/delete-shift", controllers.DeleteShift) // diakses oleh admin

	fmt.Println("Server running in : http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
