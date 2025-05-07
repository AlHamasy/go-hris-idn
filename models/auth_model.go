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

func (model AuthModel) FindByNIK(nik string) (entities.LoginEmployee, error) {
	var employee entities.LoginEmployee

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
		&employee.Photo,
		&employee.Password,
		&employee.IsAdmin,
	)

	return employee, err
}
