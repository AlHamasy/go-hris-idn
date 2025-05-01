package main

import (
	"fmt"
	"hris-idn/controllers"
	"net/http"
)

func main() {
	// baca file untuk resource yg ada di static (css, js, img)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))

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

	fmt.Println("Server running in : http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
