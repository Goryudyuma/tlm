package adlib

import (
	"strconv"

	"github.com/Goryudyuma/tlm/lib/List"
	"github.com/Goryudyuma/tlm/lib/User"
)

type Adlibs []Adlib

func (a *Adlibs) New(j []JsonAdlib) {
	for _, v := range j {
		var one Adlib
		one.New(v)
		*a = append(*a, one)
	}
}

type Adlib struct {
	List    list.List
	UserIDs user.UserIDs
}

func (a *Adlib) New(j JsonAdlib) {
	(*a).List.New(j.List)
	ret := make([]int64, len(j.UserIDs))
	for _, v := range j.UserIDs {
		one, _ := strconv.ParseInt(v, 10, 64)
		ret = append(ret, one)
	}
	(*a).UserIDs.New(ret)
}

type JsonAdlib struct {
	List    list.JsonList `json:"list"`
	UserIDs []string      `json:"userids"`
}
