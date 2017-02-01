package query

import (
	"github.com/Goryudyuma/tlm/lib/User"
	"github.com/bgpat/twtr"
)

func (q Query) Querytask(clients map[user.UserID]*twtr.Client) error {
	preparearr, err := q.preparation.Prepare(clients)

	listarr, err := q.jobs.Getalllist(clients, &preparearr)
	if err != nil {
		return err
	}

	commitlist, err := q.jobs.Task(clients, listarr)
	if err != nil {
		return err
	}

	client := clients[user.UserID(0)]
	err = commitlist.Commit(client)
	if err != nil {
		return err
	}
	return nil
}
