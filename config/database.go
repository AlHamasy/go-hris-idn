package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DBConn() (db *sql.DB, err error) {

	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "hris-idn"

	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:8889)/"+dbName)
	return db, err
}
