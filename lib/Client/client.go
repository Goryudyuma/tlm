package client

import (
	"math/rand"

	"github.com/Goryudyuma/tlm/lib/User"
	"github.com/bgpat/twtr"
)

func (uc UsersClients) Add(userid user.UserID, client *twtr.Client) {
	uc.clients[userid] = client
	uc.userids = append(uc.userids, userid)
}

func (uc UsersClients) Choice(userid user.UserID) *twtr.Client {
	if v, ok := uc.clients[userid]; ok {
		return v
	}

	return uc.clients[uc.userids[rand.Intn(len(uc.userids))]]
}

func (uc UsersClients) Len() int {
	return len(uc.userids)
}
