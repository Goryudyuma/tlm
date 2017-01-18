package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type database struct {
	Client *sql.DB
}
