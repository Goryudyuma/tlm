package adlib

import (
	"strconv"

	"github.com/Goryudyuma/tlm/lib/List"
	"github.com/Goryudyuma/tlm/lib/User"
)

type Adlibs []Adlib

func (a *Adlibs) New(j []JsonAdlib) error {
	for _, v := range j {
		var one Adlib
		if err := one.New(v); err != nil {
			return err
		}
		*a = append(*a, one)
	}
	return nil
}

type Adlib struct {
	List    list.List
	UserIDs user.UserIDs
}

func (a *Adlib) New(j JsonAdlib) error {
	if err := (*a).List.New(j.List); err != nil {
		return err
	}
	ret := make([]int64, len(j.UserIDs))
	for _, v := range j.UserIDs {
		one, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		ret = append(ret, one)
	}
	if err := (*a).UserIDs.New(ret); err != nil {
		return err
	}
	return nil
}

type JsonAdlib struct {
	List    list.JsonList `json:"list"`
	UserIDs []string      `json:"userids"`
}
