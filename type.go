package main

type Config struct {
	ConsumerKey    string `yaml:"ConsumerKey"`
	ConsumerSecret string `yaml:"ConsumerSecret"`
	SeedString     string `yaml:"SeedString"`
	CallbackURL    string `yaml:"CallbackURL"`
}

type Change struct {
	AddList []int64
	DelList []int64
}

type Adlib struct {
	List    List    `json:"list"`
	UserIDs []int64 `json:"userids"`
}

type Follower struct {
	List   List  `json:"list"`
	UserID int64 `json:"userid"`
}

type Preparation struct {
	Adlib    []Adlib    `json:"adlib"`
	Follower []Follower `json:"follower"`
}

type List struct {
	OwnerID int64  `json:"ownerid"`
	ListID  int64  `json:"listid"`
	Tag     string `json:"tag"`
}

type ResultListConfig struct {
	Name       string `json:"name"`
	Publicflag bool   `json:"publicflag"`
	Saveflag   bool   `json:"saveflag"`
}

type Job struct {
	Operator   string           `json:"operator"`
	List1      List             `json:"list1"`
	List2      List             `json:"list2"`
	Listresult List             `json:"listresult"`
	Config     ResultListConfig `json:"config"`
}

type Query struct {
	Preparation Preparation `json:"preparation"`
	Jobs        []Job       `json:"jobs"`
	Regularflag bool        `json:"regularflag"`
}
