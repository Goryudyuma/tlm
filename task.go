package main

import (
	"log"
	"sort"
	"strconv"
	"sync"

	"github.com/bgpat/twtr"
	"github.com/cznic/sortutil"
)

type taglistnode struct {
	list []int64
	ok   bool
}

func getlist(client *twtr.Client, id int64) ([]int64, error) {

	data, err := client.GetListMembers(twtr.Values{
		"list_id": strconv.FormatInt(id, 10),
		"count":   "5000",
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}
	ret := make([]int64, len(data.Users))

	for _, v := range data.Users {
		ret = append(ret, int64(v.ID))
	}

	return ret, nil
}

func getlistmembers(listarg List, client *twtr.Client,
	chanerr chan error, ret *map[List][]int64, mutex *sync.Mutex) {
	list := []int64{}

	if listarg.ListID != 0 {
		var err error
		list, err = getlist(client, listarg.ListID)
		if err != nil {
			chanerr <- err
			return
		}
	}

	/*
		//Golang1.8から以下の構文が入るかもしれない？

		sort.Slice(list, func(i, j int64) bool {
			return i < j
		})
	*/
	/*
		入るまでのつなぎ。
	*/

	var listint64 sortutil.Int64Slice
	listint64 = list
	sort.Sort(listint64)
	list = listint64

	/*
		ここまで
	*/

	mutex.Lock()
	(*ret)[listarg] = list
	mutex.Unlock()
	chanerr <- nil
}

func getalllist(jobs []Job, client *twtr.Client) (map[List][]int64, error) {
	var mutex sync.Mutex
	ret := make(map[List][]int64)
	chanerr := make(chan error, len(jobs)*3+1)
	defer close(chanerr)
	for _, v := range jobs {
		if _, ok := ret[v.List1]; !ok {
			go getlistmembers(v.List1, client, chanerr, &ret, &mutex)
		} else {
			chanerr <- nil
		}
		if _, ok := ret[v.List2]; !ok {
			go getlistmembers(v.List2, client, chanerr, &ret, &mutex)
		} else {
			chanerr <- nil
		}
		if _, ok := ret[v.Listresult]; !ok {
			go getlistmembers(v.Listresult, client, chanerr, &ret, &mutex)
		} else {
			chanerr <- nil
		}
	}

	var err error
	for i := 0; i < len(jobs)*3; i++ {
		select {
		case v := <-chanerr:
			{
				if v != nil {
					err = v
				}
			}
		}
	}
	return ret, err
}

func prepare(p Preparation, client *twtr.Client) (map[List][]int64, error) {
	ret := make(map[List][]int64)
	for _, v := range p.Adlib {
		ret[v.List] = v.UserIDs
	}

	//ここでフォロワー一覧のIDを取ってくる

	return ret, nil
}

func querytask(queryparam Query, client *twtr.Client) error {

	//今の状態
	preparearr, err := prepare(queryparam.Preparation, client)

	listarr, err := getalllist(queryparam.Jobs, client)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for k, v := range listarr {
		preparearr[k] = v
	}

	commitlist, err := jobstask(client, queryparam.Jobs, preparearr)
	if err != nil {
		return err
	}

	err = commit(client, commitlist)
	if err != nil {
		return err
	}
	return nil
}
