package models

import (
	"database/sql"
	"hris-idn/config"
	"hris-idn/entities"
	"log"
)

type AuthModel struct {
	db *sql.DB
}

func NewAuthModel() *AuthModel {
	conn, err := config.DBConn()
	if err != nil {
		log.Println("Failed connect to database:", err)
	}
	return &AuthModel{
		db: conn,
	}
}

func (model AuthModel) FindByNIK(nik string) (entities.Employee, error) {
	var employee entities.Employee
	var photo sql.NullString

	query := `
		SELECT name, email, phone, nik, gender, birth_date, photo, password, is_admin
		FROM employee 
		WHERE nik = ? AND deleted_at IS NULL
	`
	err := model.db.QueryRow(query, nik).Scan(
		&employee.Name,
		&employee.Email,
		&employee.Phone,
		&employee.NIK,
		&employee.Gender,
		&employee.BirthDate,
		&photo,
		&employee.Password,
		&employee.IsAdmin,
	)

	if err != nil {
		return employee, err
	}

	return employee, err
}
