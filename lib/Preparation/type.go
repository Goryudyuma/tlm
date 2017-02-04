package preparation

import (
	"github.com/Goryudyuma/tlm/lib/Adlib"
	"github.com/Goryudyuma/tlm/lib/Follower"
)

type Preparation struct {
	Adlib    adlib.Adlibs
	Follower follower.Followers
}

func (p *Preparation) New(j JsonPreparation) error {
	if err := (*p).Adlib.New(j.Adlib); err != nil {
		return err
	}
	if err := (*p).Follower.New(j.Follower); err != nil {
		return err
	}
	return nil
}

type JsonPreparation struct {
	Adlib    []adlib.JsonAdlib       `json:"adlib"`
	Follower []follower.JsonFollower `json:"follower"`
}
