package database

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"strconv"

	"github.com/Goryudyuma/tlm/lib/Query"
	"github.com/Goryudyuma/tlm/lib/User"
	_ "github.com/go-sql-driver/mysql"
)

type dbclients struct {
	CheckLoginInput    chan<- CheckLoginType
	CreateUserInput    chan<- CreateUserType
	AddChildUserInput  chan<- AddChildUserType
	LoginInput         chan<- LoginType
	RegisterQueryInput chan<- RegisterQueryType
	QueryAllInput      chan<- QueryAllType
}

func (db database) New(username, password, dbname string) (dbclients, chan<- bool, error) {
	dbclient, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", username, password, dbname))
	if err != nil {
		return dbclients{}, nil, err
	}
	db.Client = dbclient

	exit := make(chan bool)

	var ret dbclients

	CheckLoginTypeChan := make(chan CheckLoginType, 10)
	ret.CheckLoginInput = CheckLoginTypeChan
	go db.CheckLogin(CheckLoginTypeChan, exit)

	CreateUserTypeChan := make(chan CreateUserType, 10)
	ret.CreateUserInput = CreateUserTypeChan
	go db.CreateUser(CreateUserTypeChan, exit)

	AddChildUserTypeChan := make(chan AddChildUserType, 10)
	ret.AddChildUserInput = AddChildUserTypeChan
	go db.AddChildUser(AddChildUserTypeChan, exit)

	LoginTypeChan := make(chan LoginType, 10)
	ret.LoginInput = LoginTypeChan
	go db.Login(LoginTypeChan, exit)

	RegisterQueryTypeChan := make(chan RegisterQueryType, 10)
	ret.RegisterQueryInput = RegisterQueryTypeChan
	go db.RegisterQuery(RegisterQueryTypeChan, exit)

	QueryAllTypeChan := make(chan QueryAllType, 10)
	ret.QueryAllInput = QueryAllTypeChan
	go db.QueryAll(QueryAllTypeChan, exit)

	return ret, exit, nil
}

type CheckLoginType struct {
	UserID user.UserID
	token  string
	Exit   chan<- bool
	Err    chan<- error
}

func (db database) CheckLogin(u <-chan CheckLoginType, exit <-chan bool) {
	stmt, err := db.Client.Prepare("SELECT token FROM account WHERE id = ? AND parentid = 0 AND expiration > NOW()")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	for {
		select {
		case CheckLoginValue := <-u:
			{
				var ret string
				if err := stmt.QueryRow(CheckLoginValue.UserID).Scan(&ret); err != nil {
					CheckLoginValue.Err <- err
					continue
				}
				CheckLoginValue.Exit <- (ret == CheckLoginValue.token)
			}
		case <-exit:
			{
				break
			}
		}
	}
}

type UserToken struct {
	UserID int64
	Token  string
}

type CreateUserType struct {
	UserID            user.UserID
	AccessToken       string
	AccessTokenSecret string
	Exit              chan<- UserToken
	Err               chan<- error
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
			(parentid,userid, token, accesstoken, accesstokensecret, expiration)
		VALUES
			(0, ?, ?, ?, ?, NOW() + INTERVAL 1 DAY);`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	for {
		select {
		case CreateUserValue := <-u:
			{
				res, err := stmt.Exec(token, CreateUserValue.UserID, CreateUserValue.AccessToken, CreateUserValue.AccessTokenSecret)
				if err != nil {
					CreateUserValue.Err <- err
					continue
				}
				id, err := res.LastInsertId()
				if err != nil {
					CreateUserValue.Err <- err
					continue
				}
				CreateUserValue.Exit <- UserToken{UserID: id, Token: token}
			}
		case <-exit:
			{
				break
			}
		}
	}
}

type AddChildUserType struct {
	ParentID          int
	UserID            user.UserID
	AccessToken       string
	AccessTokenSecret string
	Exit              chan<- bool
	Err               chan<- error
}

func (db database) AddChildUser(u <-chan AddChildUserType, exit <-chan bool) {
	stmt, err := db.Client.Prepare(`
		INSERT INTO account (parentid, userid, accesstoken, accesstokensecret)
		VALUES (?, ?, ?, ?);
	`)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case AddChildUserValue := <-u:
			{
				_, err = stmt.Exec(AddChildUserValue.ParentID, AddChildUserValue.UserID, AddChildUserValue.AccessToken, AddChildUserValue.AccessTokenSecret)
				if err != nil {
					AddChildUserValue.Err <- err
					continue
				}
				AddChildUserValue.Exit <- true
			}
		case <-exit:
			{
				break
			}
		}
	}
}

type LoginType struct {
	AccessToken       string
	AccessTokenSecret string
	Exit              chan<- UserToken
	Err               chan<- error
}

func (db database) Login(u <-chan LoginType, exit <-chan bool) {
	stmt, err := db.Client.Prepare(`
		SELECT id FROM account WHERE accesstoken = ? AND accesstokensecret = ?;
	`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	stmtin, err := db.Client.Prepare(`
		UPDATE account SET token = ?, expiration = NOW() + INTERVAL 1 DAY WHERE id = ?;
	`)
	if err != nil {
		panic(err)
	}
	defer stmtin.Close()

	for {
		select {
		case LoginValue := <-u:
			{
				var ret int64
				if err := stmt.QueryRow(LoginValue.AccessToken, LoginValue.AccessTokenSecret).Scan(&ret); err != nil {
					LoginValue.Err <- err
					continue
				}
				token := createToken()
				if _, err := stmtin.Exec(token, ret); err != nil {
					LoginValue.Err <- err
					continue
				}
				LoginValue.Exit <- UserToken{
					UserID: ret,
					Token:  token}
			}
		}
	}
}

type RegisterQueryType struct {
	userid int64
	query  query.JsonQuery
	Exit   chan<- bool
	Err    chan<- error
}

func (db database) RegisterQuery(u <-chan RegisterQueryType, exit <-chan bool) {
	stmt, err := db.Client.Prepare(`
		INSERT INTO query (accountid, query, failcount) VALUES (?, ?, 0);
	`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for {
		select {
		case RegisterQueryValue := <-u:
			{
				if _, err := stmt.Exec(RegisterQueryValue.userid, RegisterQueryValue.query); err != nil {
					RegisterQueryValue.Err <- err
					continue
				}
				RegisterQueryValue.Exit <- true
			}
		case <-exit:
			{
				break
			}
		}
	}
}

type QueryAllType struct {
	userid int64
	Exit   chan<- []query.JsonQuery
	Err    chan<- error
}

func (db database) QueryAll(u <-chan QueryAllType, exit <-chan bool) {
	stmt, err := db.Client.Prepare(`
		SELECT query FROM query WHERE userid = ?;
	`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	for {
		select {
		case QueryAllValue := <-u:
			{
				rows, err := stmt.Query(QueryAllValue)
				if err != nil {
					QueryAllValue.Err <- err
					continue
				}
				ret := make([]query.JsonQuery, 0)
				for rows.Next() {
					var one query.JsonQuery
					err = rows.Scan(&one)
					if err != nil {
						QueryAllValue.Err <- err
						break
					}
					ret = append(ret, one)
				}
				if err == nil {
					QueryAllValue.Exit <- ret
				}
			}
		case <-exit:
			{
				break
			}
		}
	}
}
