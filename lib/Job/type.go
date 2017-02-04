package job

import (
	"github.com/Goryudyuma/tlm/lib/List"
	"github.com/Goryudyuma/tlm/lib/ResultListConfig"
)

type Jobs []Job

func (j *Jobs) New(json []JsonJob) error {
	for _, v := range json {
		var one Job
		if err := one.New(v); err != nil {
			return err
		}
		*j = append(*j, one)
	}
	return nil
}

type Job struct {
	Operator    string
	ListOne     list.List
	ListAnother list.List
	ListResult  list.List
	Config      resultlistconfig.ResultListConfig
}

func (j *Job) New(json JsonJob) error {
	(*j).Operator = json.Operator
	if err := (*j).ListOne.New(json.ListOne); err != nil {
		return err
	}
	if err := (*j).ListAnother.New(json.ListAnother); err != nil {
		return err
	}
	if err := (*j).ListResult.New(json.ListResult); err != nil {
		return err
	}
	if err := (*j).Config.New(json.Config); err != nil {
		return err
	}
	return nil
}

type JsonJob struct {
	Operator    string                                `json:"operator"`
	ListOne     list.JsonList                         `json:"listone"`
	ListAnother list.JsonList                         `json:"listanother"`
	ListResult  list.JsonList                         `json:"listresult"`
	Config      resultlistconfig.JsonResultListConfig `json:"config"`
}
