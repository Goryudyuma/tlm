package list

import (
	"strconv"

	"github.com/Goryudyuma/tlm/lib/Tag"
	"github.com/Goryudyuma/tlm/lib/User"
)

type ListID int64

func (l *ListID) New(j int64) error {
	*l = ListID(j)
	return nil
}

type List struct {
	OwnerID user.UserID
	ListID  ListID
	Tag     tag.Tag
}

func (l *List) New(j JsonList) error {
	ownerid, err := strconv.ParseInt(j.OwnerID, 10, 64)
	if err != nil {
		ownerid = 0
	}
	if err := (*l).OwnerID.New(ownerid); err != nil {
		return err
	}
	listid, err := strconv.ParseInt(j.ListID, 10, 64)
	if err != nil {
		return err
	}
	if err := (*l).ListID.New(listid); err != nil {
		return err
	}
	if err := (*l).Tag.New(j.Tag); err != nil {
		return err
	}
	return nil
}

type JsonList struct {
	OwnerID string `json:"ownerid"`
	ListID  string `json:"listid"`
	Tag     string `json:"tag"`
}
