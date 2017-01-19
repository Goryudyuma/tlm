package database

import (
	"database/sql"
	"fmt"

	"github.com/Goryudyuma/tlm/lib/User"
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

type CheckLoginType struct {
	UserID user.UserID
	Exit   chan string
	Err    chan error
}

func (db database) CheckLogin(u <-chan CheckLoginType, exit <-chan bool) {
	stmt, err := db.Client.Prepare("SELECT token FROM user WHERE id = ? AND owner = 1 AND lastlogin > NOW()")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	for {
		select {
		case CheckLogin := <-u:
			{
				var ret string
				if err := stmt.QueryRow(CheckLogin.UserID).Scan(&ret); err != nil {
					CheckLogin.Err <- err
				}
				CheckLogin.Exit <- ret
			}
		case <-exit:
			{
				break
			}
		}
	}
}
