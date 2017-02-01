package list

import (
	"strconv"

	"github.com/Goryudyuma/tlm/lib/Tag"
	"github.com/Goryudyuma/tlm/lib/User"
)

type ListID int64

func (l *ListID) New(j int64) {
	*l = ListID(j)
}

type List struct {
	OwnerID user.UserID
	ListID  ListID
	Tag     tag.Tag
}

func (l *List) New(j JsonList) {
	ownerid, _ := strconv.ParseInt(j.OwnerID, 10, 64)
	(*l).OwnerID.New(ownerid)
	listid, _ := strconv.ParseInt(j.ListID, 10, 64)
	(*l).ListID.New(listid)
	(*l).Tag.New(j.Tag)
}

type JsonList struct {
	OwnerID string `json:"ownerid"`
	ListID  string `json:"listid"`
	Tag     string `json:"tag"`
}
