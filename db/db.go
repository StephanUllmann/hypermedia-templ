package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init()  {
	var err error
	DB, err = sql.Open("sqlite3", "./db/contacts.sqlite3")
	if err != nil {
		panic(err)
	}
}
