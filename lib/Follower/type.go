package follower

import (
	"strconv"

	"github.com/Goryudyuma/tlm/lib/List"
	"github.com/Goryudyuma/tlm/lib/User"
)

type Followers []Follower

func (f *Followers) New(j []JsonFollower) error {
	for _, v := range j {
		var one Follower
		if err := one.New(v); err != nil {
			return err
		}
		*f = append(*f, one)
	}
	return nil
}

type Follower struct {
	List   list.List
	UserID user.UserID
}

func (f *Follower) New(j JsonFollower) error {
	if err := (*f).List.New(j.List); err != nil {
		return err
	}
	userid, err := strconv.ParseInt(j.UserID, 10, 64)
	if err != nil {
		return err
	}
	if err := (*f).UserID.New(userid); err != nil {
		return err
	}
	return nil
}

type JsonFollower struct {
	List   list.JsonList `json:"list"`
	UserID string        `json:"userid"`
}
