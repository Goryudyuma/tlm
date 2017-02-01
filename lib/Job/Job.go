package job

import (
	"sync"

	"github.com/Goryudyuma/tlm/lib/Change"
	"github.com/Goryudyuma/tlm/lib/List"
	"github.com/Goryudyuma/tlm/lib/User"

	"github.com/bgpat/twtr"
)

func (v Job) dojob(clients map[user.UserID]*twtr.Client, result, origin *map[list.List]user.UserIDs,
	ret *map[list.ListID]change.Change, listids *map[list.List]list.ListID) error {
	l1 := (*result)[v.ListOne]
	l2 := (*result)[v.ListAnother]
	switch v.Operator {
	case "*":
		(*result)[v.ListResult] = l1.Intersect(l2)
	case "+":
		(*result)[v.ListResult] = l1.Union(l2)
	case "-":
		(*result)[v.ListResult] = l1.Except(l2)
	}

	if v.Config.Saveflag {
		addval := change.Change{
			AddList: (*result)[v.ListResult].Except((*origin)[v.ListResult]),
			DelList: (*origin)[v.ListResult].Except((*result)[v.ListResult])}

		listid, ok := (*listids)[v.ListResult]
		if !ok {
			listid = v.ListResult.ListID

			if listid == 0 {
				var mode string
				if v.Config.Publicflag {
					mode = "public"
				} else {
					mode = "private"
				}

				client := clients[user.UserID(0)]
				createlist, err := client.CreateList(&twtr.Values{
					"name": v.Config.Name,
					"mode": mode,
				})
				if err != nil {
					return err
				}

				listid = list.ListID(createlist.ID.ID)
			}
			(*listids)[v.ListResult] = listid
		}
		(*ret)[listid] = addval
	}
	return nil
}

func (jobs Jobs) Task(clients map[user.UserID]*twtr.Client, origin map[list.List]user.UserIDs) (
	change.Changes, error) {
	listids := make(map[list.List]list.ListID)
	ret := make(map[list.ListID]change.Change)

	result := make(map[list.List]user.UserIDs, len(origin))
	for k, v := range origin {
		result[k] = v
	}

	for _, v := range jobs {
		v.dojob(clients, &result, &origin, &ret, &listids)
	}
	return ret, nil
}

func (job Job) GetListMember(clients map[user.UserID]*twtr.Client, ret *map[list.List]user.UserIDs, chanerr chan error, mutex *sync.Mutex) {
	if _, ok := (*ret)[job.ListOne]; !ok {
		client, ok := clients[job.ListOne.OwnerID]
		if !ok {
			client = clients[user.UserID(0)]
		}
		go job.ListOne.GetListMembers(client, chanerr, ret, mutex)
	} else {
		chanerr <- nil
	}
	if _, ok := (*ret)[job.ListAnother]; !ok {
		client, ok := clients[job.ListAnother.OwnerID]
		if !ok {
			client = clients[user.UserID(0)]
		}
		go job.ListAnother.GetListMembers(client, chanerr, ret, mutex)
	} else {
		chanerr <- nil
	}
	if _, ok := (*ret)[job.ListResult]; !ok {
		client, ok := clients[job.ListResult.OwnerID]
		if !ok {
			client = clients[user.UserID(0)]
		}
		go job.ListResult.GetListMembers(client, chanerr, ret, mutex)
	} else {
		chanerr <- nil
	}
}

func (j Jobs) Getalllist(clients map[user.UserID]*twtr.Client, ret *map[list.List]user.UserIDs) (map[list.List]user.UserIDs, error) {
	var mutex sync.Mutex
	chanerr := make(chan error, len(j)*3+1)
	defer close(chanerr)
	for _, v := range j {
		v.GetListMember(clients, ret, chanerr, &mutex)
	}

	var err error
	for i := 0; i < len(j)*3; i++ {
		select {
		case v := <-chanerr:
			{
				if v != nil {
					err = v
				}
			}
		}
	}
	return *ret, err
}
