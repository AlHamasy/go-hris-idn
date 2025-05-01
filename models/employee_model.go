package models

import (
	"database/sql"
	"hris-idn/config"
	"hris-idn/entities"
	"hris-idn/helpers"
	"log"
	"time"
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
		var photo sql.NullString
		err := rows.Scan(
			&employee.UUID,
			&employee.NIK,
			&employee.Name,
			&employee.Email,
			&employee.Phone,
			&employee.Gender,
			&employee.IsAdmin,
			&photo,
			&employee.BirthDate,
		)
		if err != nil {
			return []entities.Employee{}, err
		}

		if photo.Valid {
			employee.Photo = photo.String
		} else {
			if employee.Gender == "M" {
				employee.Photo = helpers.MALE_BASE64
			} else if employee.Gender == "F" {
				employee.Photo = helpers.FEMALE_BASE64
			}
		}

		employees = append(employees, employee)
	}

	return employees, nil
}

func (model EmployeeModel) FindEmployeeByUUID(uuid string) (entities.Employee, error) {
	var employee entities.Employee

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
		&employee.BirthDate,
	)

	if err != nil {
		return employee, err
	}

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
