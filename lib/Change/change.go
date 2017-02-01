package change

import (
	"strconv"
	"strings"
	"time"

	"github.com/bgpat/twtr"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (change Changes) Commit(client *twtr.Client) error {
	var err error
	for id, v := range change {
		i := 0
		for {
			_, err := client.GetList(&twtr.Values{
				"list_id": strconv.FormatInt(int64(id), 10),
			})
			if err != nil {
				//5回リストがあるかどうかチェックして、それでも無ければerrorとして返す。
				i++
				if i > 5 {
					return err
				}
				time.Sleep(1 * time.Second)
			} else {
				break
			}
		}
		for len(v.DelList) != 0 {
			list := make([]string, 0, 100)
			handled := v.DelList[:min(100, len(v.DelList))]
			for _, one := range handled {
				list = append(list, strconv.FormatInt(int64(one), 10))
			}
			count := 0
			for count < 10 {
				_, err := client.DeleteListMembers(&twtr.Values{
					"list_id": strconv.FormatInt(int64(id), 10),
					"user_id": strings.Join(list[:], ","),
				})
				if err == nil {
					break
				}
			}
			if err != nil {
				return err
			}
			v.DelList = v.DelList[min(100, len(v.DelList)):]
		}
		for len(v.AddList) != 0 {
			list := make([]string, 0, 100)
			handled := v.AddList[:min(100, len(v.AddList))]
			for _, one := range handled {
				list = append(list, strconv.FormatInt(int64(one), 10))
			}
			count := 0
			for count < 10 {
				_, err := client.AddListMembers(&twtr.Values{
					"list_id": strconv.FormatInt(int64(id), 10),
					"user_id": strings.Join(list[:], ","),
				})
				if err == nil {
					break
				}
			}
			if err != nil {
				return err
			}
			v.AddList = v.AddList[min(100, len(v.AddList)):]
		}
	}
	return nil
}
