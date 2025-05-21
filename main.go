package main

import (
	"fmt"
	"hris-idn/config"
	"hris-idn/controllers"
	"net/http"
)

func main() {
	// baca file untuk resource yg ada di static (css, js, img)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))

	// login & logout
	http.HandleFunc("/login", config.GuestMiddleware(controllers.Login)) // bebas diakses (belum login)
	http.HandleFunc("/logout", controllers.Logout)                       // harus login (admin & non-admin)

	// home
	http.HandleFunc("/home", config.NonAdminOnly(controllers.Home))         // hanya non-admin
	http.HandleFunc("/home-admin", config.AdminOnly(controllers.HomeAdmin)) // hanya admin

	// news
	http.HandleFunc("/news", config.AdminOnly(controllers.News)) // hanya admin

	// attendance
	http.HandleFunc("/attendance-submit", config.NonAdminOnly(controllers.SubmitAttendance)) // hanya non-admin
	http.HandleFunc("/attendance-list", config.AdminOnly(controllers.ListAttendance))        // hanya admin

	// leave
	http.HandleFunc("/leave-submit", config.NonAdminOnly(controllers.SubmitLeave)) // hanya non-admin
	http.HandleFunc("/leave-list", config.AdminOnly(controllers.ListLeave))        // hanya admin
	http.HandleFunc("/leave/approval", config.AdminOnly(controllers.ApprovalLeave))

	// payroll
	http.HandleFunc("/payroll-self", config.NonAdminOnly(controllers.MyPayroll)) // hanya non-admin
	http.HandleFunc("/payroll-list", config.AdminOnly(controllers.ListPayroll))  // hanya admin

	// employee
	http.HandleFunc("/employee", config.AdminOnly(controllers.Employee))                       // hanya admin
	http.HandleFunc("/employee/add-employee", config.AdminOnly(controllers.AddEmployee))       // hanya admin
	http.HandleFunc("/employee/edit-employee", config.AdminOnly(controllers.EditEmployee))     // hanya admin
	http.HandleFunc("/employee/delete-employee", config.AdminOnly(controllers.DeleteEmployee)) // hanya admin

	// office
	http.HandleFunc("/office", config.AdminOnly(controllers.Office))                     // hanya admin
	http.HandleFunc("/office/add-office", config.AdminOnly(controllers.AddOffice))       // hanya admin
	http.HandleFunc("/office/edit-office", config.AdminOnly(controllers.EditOffice))     // hanya admin
	http.HandleFunc("/office/delete-office", config.AdminOnly(controllers.DeleteOffice)) // hanya admin

	// shift
	http.HandleFunc("/shift", config.AdminOnly(controllers.Shift))                    // hanya admin
	http.HandleFunc("/shift/add-shift", config.AdminOnly(controllers.AddShift))       // hanya admin
	http.HandleFunc("/shift/edit-shift", config.AdminOnly(controllers.EditShift))     // hanya admin
	http.HandleFunc("/shift/delete-shift", config.AdminOnly(controllers.DeleteShift)) // hanya admin

	fmt.Println("Server running in : http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
