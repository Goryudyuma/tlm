package database

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"strconv"

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
	token  string
	Exit   chan bool
	Err    chan error
}

func (db database) CheckLogin(u <-chan CheckLoginType, exit <-chan bool) {
	stmt, err := db.Client.Prepare("SELECT token FROM account WHERE id = ? AND parent = 0 AND expiration > NOW()")
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
					continue
				}
				CheckLogin.Exit <- (ret == CheckLogin.token)
			}
		case <-exit:
			{
				break
			}
		}
	}
}

type UserToken struct {
	UserID user.UserID
	Token  string
}

type CreateUserType struct {
	AccessToken       string
	AccessTokenSecret string
	Exit              chan UserToken
	Err               chan error
}

func createToken() string {
	const base = 36
	const length = 250
	size := big.NewInt(base)
	n := make([]byte, length)
	for i, _ := range n {
		c, _ := rand.Int(rand.Reader, size)
		n[i] = strconv.FormatInt(c.Int64(), base)[0]
	}
	return string(n)
}

func (db database) CreateUser(u <-chan CreateUserType, exit <-chan bool) {
	token := createToken()
	stmt, err := db.Client.Prepare(`
		INSERT INTO account 
			(parent, token, accesstoken, accesstokensecret, expiration)
		VALUES
			(0, ?, ?, ?, NOW() + INTERVAL 1 DAY)`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	for {
		select {
		case CreateUser := <-u:
			{
				res, err := stmt.Exec(token, CreateUser.AccessToken, CreateUser.AccessTokenSecret)
				if err != nil {
					CreateUser.Err <- err
					continue
				}
				id, err := res.LastInsertId()
				if err != nil {
					CreateUser.Err <- err
					continue
				}
				CreateUser.Exit <- UserToken{UserID: user.UserID(id), Token: token}
			}
		case <-exit:
			{
				break
			}
		}
	}
}
