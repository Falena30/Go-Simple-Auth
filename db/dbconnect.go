package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//Connect digunakan untuk menghubungkan ke DB
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_db")
	if err != nil {
		return nil, err
	}
	return db, nil
}
