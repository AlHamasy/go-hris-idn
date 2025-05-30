package models

import (
	"database/sql"
	"hris-idn/config"
	"hris-idn/entities"
	"log"
	"time"

	"github.com/goodsign/monday"
)

type EmployeeModel struct {
	db *sql.DB
}

func NewEmployeeModel() *EmployeeModel {
	conn, err := config.DBConn()
	if err != nil {
		log.Println("Failed connect to database:", err)
	}
	return &EmployeeModel{
		db: conn,
	}
}

func (model EmployeeModel) AddEmployee(user entities.Employee) error {

	_, err := model.db.Exec(
		"INSERT INTO employee (uuid, name, email, phone, address, nik, is_admin, password, gender, birth_date) VALUES (?,?,?,?,?,?,?,?,?,?)",
		user.UUID, user.Name, user.Email, user.Phone, user.Address, user.NIK, user.IsAdmin, user.Password, user.Gender, user.BirthDate)

	return err
}

func (model EmployeeModel) FindAllEmployee() ([]entities.Employee, error) {

	rows, err := model.db.Query(`
		SELECT uuid, nik, name, email, phone, gender, is_admin, photo, birth_date 
		FROM employee 
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return []entities.Employee{}, err
	}
	defer rows.Close()

	var employees []entities.Employee

	for rows.Next() {
		var employee entities.Employee
		var birthDateTime time.Time
		err := rows.Scan(
			&employee.UUID,
			&employee.NIK,
			&employee.Name,
			&employee.Email,
			&employee.Phone,
			&employee.Gender,
			&employee.IsAdmin,
			&employee.Photo,
			&birthDateTime,
		)
		if err != nil {
			return []entities.Employee{}, err
		}

		employee.BirthDate = monday.Format(birthDateTime, "02 January 2006", monday.LocaleIdID)

		employees = append(employees, employee)
	}

	return employees, nil
}

func (model EmployeeModel) FindEmployeeByUUID(uuid string) (entities.Employee, error) {
	var employee entities.Employee
	var birthDate time.Time

	query := "SELECT uuid, nik, name, email, phone, gender, is_admin, address, birth_date FROM employee WHERE uuid = ?"

	err := model.db.QueryRow(query, uuid).Scan(
		&employee.UUID,
		&employee.NIK,
		&employee.Name,
		&employee.Email,
		&employee.Phone,
		&employee.Gender,
		&employee.IsAdmin,
		&employee.Address,
		&birthDate,
	)

	if err != nil {
		return employee, err
	}

	employee.BirthDate = birthDate.Format("2006-01-02")

	return employee, nil
}

func (model EmployeeModel) EditEmployee(employee entities.EditEmployee) error {

	query := `
		UPDATE employee 
		SET name = ?, email = ?, phone = ?, gender = ?, is_admin = ?, address = ?, nik = ?, updated_at = ?, birth_date = ?
		WHERE uuid = ?
	`
	_, err := model.db.Exec(
		query,
		employee.Name,
		employee.Email,
		employee.Phone,
		employee.Gender,
		employee.IsAdmin,
		employee.Address,
		employee.NIK,
		time.Now(),
		employee.BirthDate,
		employee.UUID,
	)

	return err
}

func (model EmployeeModel) SoftDeleteEmployee(uuid string) error {
	query := `
		UPDATE employee 
		SET deleted_at = ? 
		WHERE uuid = ?
	`
	_, err := model.db.Exec(query, time.Now(), uuid)
	return err
}
