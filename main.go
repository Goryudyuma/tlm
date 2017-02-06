// tlm project tlm.go
package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"

	"github.com/Goryudyuma/tlm/lib/Database"
	q "github.com/Goryudyuma/tlm/lib/Query"
	u "github.com/Goryudyuma/tlm/lib/User"

	"github.com/bgpat/twtr"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

var (
	config     Config
	clientmain *twtr.Client
	dbclients  database.DBClients
)

func loadconfig() []byte {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	return data
}

func loadyaml() Config {
	key := Config{}

	err := yaml.Unmarshal(loadconfig(), &key)
	if err != nil {
		panic(err)
	}
	return key
}

func checklogin(c *gin.Context) bool {
	session := sessions.Default(c)

	userid := session.Get("UserID")
	token := session.Get("Token")

	if userid == nil || token == nil {
		return false
	}

	exit := make(chan bool)
	err := make(chan error)
	dbclients.CheckLoginInput <- database.CheckLoginType{
		UserID: userid.(int64),
		Token:  token.(string),
		Exit:   exit,
		Err:    err,
	}

	select {
	case v := <-exit:
		{
			return v
		}
	case <-err:
		{
			return false
		}
	}
	return false
}

func getroot(c *gin.Context) {
	if !checklogin(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Redirect(303, "/login")
	} else {
		c.Redirect(303, "/static/main.html")
		//c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}

func login(c *gin.Context) {
	//if checklogin(c) {
	//c.Redirect(303, "/")
	//}
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	url, err := clientmain.RequestTokenURL(config.URL + "/callback")
	if err != nil {
		c.HTML(500, err.Error(), nil)
	}
	c.Redirect(303, url)
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Redirect(303, "/")
}

func callback(c *gin.Context) {
	session := sessions.Default(c)
	err := clientmain.GetAccessToken(c.Query("oauth_verifier"))

	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}

	userid := session.Get("UserID")
	token := session.Get("Token")

	myid, err := getmyid(clientmain.AccessToken.Token, clientmain.AccessToken.Secret)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}

	if userid == nil || token == nil || !checklogin(c) {
		exit := make(chan database.UserToken)
		reterr := make(chan error)

		dbclients.CreateUserInput <- database.CreateUserType{
			UserID:            myid,
			AccessToken:       clientmain.AccessToken.Token,
			AccessTokenSecret: clientmain.AccessToken.Secret,
			Exit:              exit,
			Err:               reterr,
		}
		select {
		case v := <-exit:
			{
				session.Set("UserID", v.UserID)
				session.Set("Token", v.Token)
				session.Save()
			}
		case err = <-reterr:
			{
				c.JSON(500, gin.H{"status": "error", "data": err.Error()})
				return
			}
		}
	} else {
		exit := make(chan bool)
		reterr := make(chan error)

		dbclients.AddChildUserInput <- database.AddChildUserType{
			ParentID:          userid.(int64),
			UserID:            myid,
			AccessToken:       clientmain.AccessToken.Token,
			AccessTokenSecret: clientmain.AccessToken.Secret,
			Exit:              exit,
			Err:               reterr,
		}
		select {
		case <-exit:
			{

			}
		case err = <-reterr:
			{
				c.JSON(500, gin.H{"status": "error", "data": err.Error()})
				return
			}
		}

	}

	c.Redirect(303, "/")
}

func getmyid(Token, Secret string) (u.UserID, error) {
	client, err := createclient(Token, Secret)
	if err != nil {
		return u.UserID(0), err
	}
	user, err := client.VerifyCredentials(nil)
	if err != nil {
		return u.UserID(0), err
	}
	return u.UserID(user.ID.ID), nil
}

func createclient(accesstoken, accesstokensecret string) (*twtr.Client, error) {
	consumer := oauth.Credentials{Token: config.ConsumerKey, Secret: config.ConsumerSecret}
	token := oauth.Credentials{Token: accesstoken, Secret: accesstokensecret}

	return twtr.NewClient(&consumer, &token), nil
}

func createclientparse(c *gin.Context) (map[u.UserID]*twtr.Client, error) {
	if !checklogin(c) {
		return nil, errors.New("Not login")
	}

	session := sessions.Default(c)

	userid := session.Get("UserID").(int64)
	token := session.Get("Token").(string)
	exit := make(chan map[u.UserID]database.TokenSecretType)
	reterr := make(chan error)

	dbclients.GetTokenSecretInput <- database.GetTokenSecretType{
		Id:    userid,
		Token: token,
		Exit:  exit,
		Err:   reterr,
	}

	ret := make(map[u.UserID]*twtr.Client)
	select {
	case v := <-exit:
		{
			for userid, one := range v {
				var err error
				ret[userid], err = createclient(one.AccessToken, one.AccessTokenSecret)
				if err != nil {
					return nil, err
				}
			}
		}
	case v := <-reterr:
		{
			return nil, v
		}
	}

	return ret, nil
}

func query(c *gin.Context) {

	querystring := c.PostForm("query")
	if querystring == "" {
		c.JSON(500, gin.H{"status": "error", "data": "Query parameters are missing."})
		return
	}

	var jsonquery q.JsonQuery
	err := json.Unmarshal([]byte(querystring), &jsonquery)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	client, err := createclientparse(c)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}

	var queryone q.Query
	err = queryone.New(jsonquery)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	err = queryone.Querytask(client)

	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}

	if jsonquery.Regularflag {
		myuser, err := client[0].VerifyCredentials(nil)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "data": err.Error()})
			return
		}
		exit := make(chan bool)
		reterr := make(chan error)
		dbclients.RegisterQueryInput <- database.RegisterQueryType{
			UserID: myuser.ID.ID,
			Query:  querystring,
			Exit:   exit,
			Err:    reterr,
		}
		select {
		case <-exit:
			{
			}
		case err := <-reterr:
			{
				c.JSON(500, gin.H{"status": "error", "data": err.Error()})
				return
			}
		}
	}

	c.JSON(200, gin.H{"status": "ok", "data": ""})
}

func searchuser(c *gin.Context) {
	clients, err := createclientparse(c)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	client := clients[u.UserID(0)]

	username := c.PostForm("username")

	if username == "" {
		c.JSON(500, gin.H{"status": "error", "data": "Query parameters are missing."})
		return
	}

	users, err := client.SearchUsers(&twtr.Values{
		"q":     username,
		"count": "100",
	})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	var ret [][2]string
	for _, v := range users {
		ret = append(ret, [2]string{
			v.ScreenName,
			v.IDStr,
		})
	}
	c.JSON(200, gin.H{"status": "ok", "data": ret})
}

func userlist(c *gin.Context) {
	clients, err := createclientparse(c)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	userid := c.PostForm("userid")

	if userid == "" {
		c.JSON(500, gin.H{"status": "error", "data": "Query parameters are missing."})
		return
	}

	useridint64, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": "Query parameters are missing."})
		return
	}
	client, ok := clients[u.UserID(useridint64)]
	if !ok {
		client = clients[u.UserID(0)]
	}

	lists, err := client.GetLists(&twtr.Values{
		"user_id": userid,
	})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	var ret [][2]string
	for _, v := range lists {
		ret = append(ret, [2]string{
			v.Name,
			v.ID.IDStr,
		})
	}
	c.JSON(200, gin.H{"status": "ok", "data": ret})
}

func getusers(c *gin.Context) {
	clients, err := createclientparse(c)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	userids := c.PostForm("userids")

	if userids == "" {
		c.JSON(500, gin.H{"status": "error", "data": "Query parameters are missing."})
		return
	}

	client := clients[u.UserID(0)]

	users, err := client.GetUsers(&twtr.Values{
		"user_id": userids,
	})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}

	var ret [][2]string
	for _, v := range users {
		ret = append(ret, [2]string{
			v.ScreenName,
			v.IDStr,
		})
	}
	c.JSON(200, gin.H{"status": "ok", "data": ret})
}

func test(c *gin.Context) {
	db, err := sql.Open("mysql", "test:@/test")
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	rows, err := db.Query("SELECT * FROM test")
	defer rows.Close()

	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}

	var ret [][]int

	for rows.Next() {
		var A, B, C int
		if err := rows.Scan(&A, &B, &C); err != nil {
			c.JSON(500, gin.H{"status": "error", "data": err.Error()})
			return
		}
		ret = append(ret, []int{A, B, C})
	}

	c.JSON(200, gin.H{"status": "ok", "data": ret})
}

func main() {
	config = loadyaml()
	consumer := oauth.Credentials{Token: config.ConsumerKey, Secret: config.ConsumerSecret}
	clientmain = twtr.NewClient(&consumer, nil)

	var database database.Database
	var err error

	dbclients, err = database.NewDBClients(config.DBUser, config.DBPass, config.DBName)

	if err != nil {
		panic(err)
	}

	_ = clientmain
	r := gin.Default()
	r.LoadHTMLGlob("content/index.html")

	store := sessions.NewCookieStore([]byte(config.SeedString))
	//store.Options(sessions.Options{Secure: true})
	r.Use(sessions.Sessions("tlcsession", store))

	r.GET("/", getroot)
	r.GET("/login", login)
	r.GET("/logout", logout)
	r.GET("/callback", callback)
	r.Static("/static", "./content")

	r.GET("/test", test)

	rapi := r.Group("/api")
	{
		rapi.POST("/query", query)
		rapi.POST("/userlist", userlist)
		rapi.POST("/searchuser", searchuser)
		rapi.POST("/getusers", getusers)
	}

	r.Run(":" + config.Port)
}
