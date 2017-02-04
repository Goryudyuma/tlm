package query

import (
	"github.com/Goryudyuma/tlm/lib/Job"
	"github.com/Goryudyuma/tlm/lib/Preparation"
)

type Query struct {
	preparation preparation.Preparation
	jobs        job.Jobs
	regularflag bool
}

func (query *Query) New(j JsonQuery) error {
	if err := (*query).preparation.New(j.Preparation); err != nil {
		return err
	}
	if err := (*query).jobs.New(j.Jobs); err != nil {
		return err
	}
	(*query).regularflag = j.Regularflag

	return nil
}

type JsonQuery struct {
	Preparation preparation.JsonPreparation `json:"preparation"`
	Jobs        []job.JsonJob               `json:"jobs"`
	Regularflag bool                        `json:"regularflag"`
}
