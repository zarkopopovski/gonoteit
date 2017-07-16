package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type DBConnection struct {
	db *sql.DB
}

func CreateDBConnection() *DBConnection {
	dbConnection := &DBConnection{}
	dbConnection.OpenConnectionSession()

	return dbConnection
}

func (dbConnection *DBConnection) OpenConnectionSession() (err error) {
	db, err := sql.Open("sqlite3", "./gonoteit.db")
	if err != nil {
		panic(err)
	}

	fmt.Println("SQLite Connection is Active")
	dbConnection.db = db

	dbConnection.setupInitialDatabase()

	return
}

func (dbConnection *DBConnection) setupInitialDatabase() (err error) {
	statement, _ := dbConnection.db.Prepare("CREATE TABLE IF NOT EXISTS users (id VARCHAR PRIMARY KEY, username VARCHAR, email VARCHAR, password VARCHAR)")
	statement.Exec()

	statement, _ = dbConnection.db.Prepare("CREATE TABLE IF NOT EXISTS notes (id VARCHAR PRIMARY KEY, user_id VARCHAR, title VARCHAR, body VARCHAR, date_created VARCHAR)")
	statement.Exec()

	return
}
