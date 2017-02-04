package change

import (
	"github.com/Goryudyuma/tlm/lib/List"
	"github.com/Goryudyuma/tlm/lib/User"
)

type Changes map[list.ListID]Change

type Change struct {
	AddList user.UserIDs
	DelList user.UserIDs
}

func (c *Change) New(j JsonChange) error {
	if err := (*c).AddList.New(j.AddList); err != nil {
		return err
	}
	if err := (*c).DelList.New(j.DelList); err != nil {
		return err
	}
	return nil
}

type JsonChange struct {
	AddList []int64
	DelList []int64
}
