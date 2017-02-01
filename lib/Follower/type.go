package follower

import (
	"strconv"

	"github.com/Goryudyuma/tlm/lib/List"
	"github.com/Goryudyuma/tlm/lib/User"
)

type Followers []Follower

func (f *Followers) New(j []JsonFollower) {
	for _, v := range j {
		var one Follower
		one.New(v)
		*f = append(*f, one)
	}
}

type Follower struct {
	List   list.List
	UserID user.UserID
}

func (f *Follower) New(j JsonFollower) {
	(*f).List.New(j.List)
	userid, _ := strconv.ParseInt(j.UserID, 10, 64)
	(*f).UserID.New(userid)
}

type JsonFollower struct {
	List   list.JsonList `json:"list"`
	UserID string        `json:"userid"`
}
