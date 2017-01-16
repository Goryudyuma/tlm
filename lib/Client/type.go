package client

import (
	"github.com/Goryudyuma/tlm/lib/User"
	"github.com/bgpat/twtr"
)

type UsersClients struct {
	clients map[user.UserID]*twtr.Client
	userids user.UserIDs
}
