package list

import (
	"log"
	"sort"
	"strconv"
	"sync"

	"github.com/Goryudyuma/tlm/lib/User"

	"github.com/bgpat/twtr"
)

func (listid ListID) Getlist(client *twtr.Client) (user.UserIDs, error) {

	data, err := client.GetListMembers(&twtr.Values{
		"list_id": strconv.FormatInt(int64(listid), 10),
		"count":   "5000",
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}
	ret := make(user.UserIDs, 0, len(data.Users))

	for _, v := range data.Users {
		ret = append(ret, user.UserID(v.ID.ID))
	}

	return ret, nil
}

func (listarg List) GetListMembers(client *twtr.Client,
	chanerr chan error, ret *map[List]user.UserIDs, mutex *sync.Mutex) {
	list := user.UserIDs{}

	if listarg.ListID != 0 {
		var err error
		list, err = listarg.ListID.Getlist(client)
		if err != nil {
			chanerr <- err
			return
		}
	}

	sort.Sort(list)

	mutex.Lock()
	(*ret)[listarg] = list
	mutex.Unlock()
	chanerr <- nil
}
