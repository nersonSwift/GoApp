package sqlServise

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type SQLConnection struct {
	Link string

	DB *sql.DB
}

func Connect(link string) *SQLConnection {
	sqlConnection := SQLConnection{
		Link: link,
	}
	sqlConnection.connect()
	return &sqlConnection
}

func (sqlC *SQLConnection) connect() {
	sqlC.setConnection()
}

func (sqlC *SQLConnection) reconnect() {
	sqlC.DB = nil
	sqlC.connect()
}

func (sqlC *SQLConnection) setConnection() {
	db, err := sql.Open("mysql", sqlC.Link)
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlC.DB = db
}

func (sqlC *SQLConnection) Exec(query string, args ...interface{}) error {
	if sqlC.DB == nil {
		sqlC.reconnect()
		if sqlC.DB == nil {
			return errors.New("No DB")
		}
	}
	err := sqlC.DB.Ping()
	if err != nil {
		sqlC.DB = nil
		return err
	}

	result, err := sqlC.DB.Exec(query, args)
	fmt.Println(result)
	return err
}
