// tlm project tlm.go
package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	//	"sync"

	"github.com/bgpat/twtr"
	"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

var (
	config     Config
	clientmain *twtr.Client
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
	OauthToken := session.Get("OauthToken")
	OauthTokenSecret := session.Get("OauthTokenSecret")
	if OauthToken == nil || OauthTokenSecret == nil {
		return false
	}
	return true
}

func getroot(c *gin.Context) {
	if !checklogin(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Redirect(301, "/login")
	} else {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}

func login(c *gin.Context) {
	if checklogin(c) {
		c.Redirect(301, "/")
	}
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	url, err := clientmain.RequestTokenURL(config.CallbackURL)
	if err != nil {
		c.HTML(500, err.Error(), nil)
	}
	c.Redirect(301, url)
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Redirect(301, "/")
}

func callback(c *gin.Context) {
	session := sessions.Default(c)
	err := clientmain.GetAccessToken(c.Query("oauth_verifier"))

	if err != nil {

		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}

	//spew.Dump(clientmain.GetLists(url.Values{}))
	spew.Dump(clientmain)
	session.Set("OauthToken", clientmain.AccessToken.Token)
	session.Set("OauthTokenSecret", clientmain.AccessToken.Secret)
	session.Save()

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Redirect(301, "/")
}

func createclient(c *gin.Context) (*twtr.Client, error) {
	if !checklogin(c) {
		return nil, errors.New("Not login")
	}

	session := sessions.Default(c)

	var OauthToken, OauthTokenSecret string
	if session.Get("OauthToken") != nil {
		OauthToken = session.Get("OauthToken").(string)
	} else {
		return nil, errors.New("Not login")
	}
	if session.Get("OauthTokenSecret") != nil {
		OauthTokenSecret = session.Get("OauthTokenSecret").(string)
	} else {
		return nil, errors.New("Not login")
	}

	consumer := oauth.Credentials{Token: config.ConsumerKey, Secret: config.ConsumerSecret}
	token := oauth.Credentials{Token: OauthToken, Secret: OauthTokenSecret}

	return twtr.NewClient(&consumer, &token), nil
}

func query(c *gin.Context) {

	querystring := c.PostForm("query")
	var queryone Query
	err := json.Unmarshal([]byte(querystring), &queryone)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	client, err := createclient(c)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	err = querytask(queryone, client)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	c.JSON(200, queryone)
}

func userlist(c *gin.Context) {
	client, err := createclient(c)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	userid := c.PostForm("userid")
	lists, err := client.GetLists(twtr.Values{
		"userid": userid,
	})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": err.Error()})
		return
	}
	c.JSON(200, lists)
}

func main() {
	config = loadyaml()
	consumer := oauth.Credentials{Token: config.ConsumerKey, Secret: config.ConsumerSecret}
	clientmain = twtr.NewClient(&consumer, nil)

	_ = clientmain
	r := gin.Default()
	r.LoadHTMLGlob("content/*")

	store := sessions.NewCookieStore([]byte(config.SeedString))
	//store.Options(sessions.Options{Secure: true})
	r.Use(sessions.Sessions("tlcsession", store))

	r.GET("/", getroot)
	r.GET("/login", login)
	r.GET("/logout", logout)
	r.GET("/callback", callback)

	rapi := r.Group("/api")
	{
		rapi.POST("/query", query)
		rapi.POST("/userlist", userlist)
	}

	r.Run()
}
