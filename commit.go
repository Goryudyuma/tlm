package main

import (
	"strconv"
	"strings"

	"github.com/bgpat/twtr"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func commit(client *twtr.Client, change map[int64]Change) error {
	for id, v := range change {
		i := 0
		for {
			_, err := client.GetList(twtr.Values{
				"list_id": strconv.FormatInt(id, 10),
			})
			if err != nil {
				//5回リストがあるかどうかチェックして、それでも無ければerrorとして返す。
				i++
				if i > 5 {
					return err
				}
			} else {
				break
			}
		}
		for len(v.DelList) != 0 {
			list := make([]string, 0, 100)
			handled := v.DelList[:min(100, len(v.DelList))]
			for _, one := range handled {
				list = append(list, strconv.FormatInt(one, 10))
			}
			v.DelList = v.DelList[min(100, len(v.DelList)):]
			_, err := client.DeleteListMembers(twtr.Values{
				"list_id": strconv.FormatInt(id, 10),
				"user_id": strings.Join(list[:], ","),
			})
			if err != nil {
				return err
			}
		}
		for len(v.AddList) != 0 {
			list := make([]string, 0, 100)
			handled := v.AddList[:min(100, len(v.AddList))]
			for _, one := range handled {
				list = append(list, strconv.FormatInt(one, 10))
			}
			v.AddList = v.AddList[min(100, len(v.AddList)):]
			_, err := client.AddListMembers(twtr.Values{
				"list_id": strconv.FormatInt(id, 10),
				"user_id": strings.Join(list[:], ","),
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
