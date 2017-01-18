package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (db database) New(username, password, dbname string) error {
	dbclient, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", username, password, dbname))
	if err != nil {
		return err
	}
	db.Client = dbclient
	return nil
}
