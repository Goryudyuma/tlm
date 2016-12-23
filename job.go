package main

import "github.com/bgpat/twtr"

func intersect(l1, l2 []int64) []int64 {
	var ret []int64
	i, j := 0, 0
	for i < len(l1) && j < len(l2) {
		if l1[i] == l2[j] {
			ret = append(ret, l1[i])
			i++
			j++
		} else if l1[i] < l2[j] {
			i++
		} else {
			j++
		}
	}
	return ret
}

func union(l1, l2 []int64) []int64 {
	var ret []int64
	i, j := 0, 0
	for i < len(l1) && j < len(l2) {
		if l1[i] == l2[j] {
			ret = append(ret, l1[i])
			i++
			j++
		} else if l1[i] < l2[j] {
			ret = append(ret, l1[i])
			i++
		} else {
			ret = append(ret, l2[i])
			j++
		}
	}
	ret = append(ret, l1[i:]...)
	ret = append(ret, l2[j:]...)

	return ret
}

func except(l1, l2 []int64) []int64 {
	var ret []int64
	i, j := 0, 0
	for i < len(l1) && j < len(l2) {
		if l1[i] == l2[j] {
			i++
			j++
		} else if l1[i] < l2[j] {
			ret = append(ret, l1[i])
			i++
		} else {
			j++
		}
	}
	ret = append(ret, l1[i:]...)
	return ret
}

func jobstask(client *twtr.Client, jobs []Job, origin map[List][]int64) (map[int64]Change, error) {
	listids := make(map[List]int64)
	ret := make(map[int64]Change)

	result := make(map[List][]int64, len(origin))
	for k, v := range origin {
		result[k] = v
	}

	for _, v := range jobs {
		l1 := result[v.List1]
		l2 := result[v.List2]
		switch v.Operator {
		case "*":
			result[v.Listresult] = intersect(l1, l2)
		case "+":
			result[v.Listresult] = union(l1, l2)
		case "-":
			result[v.Listresult] = except(l1, l2)
		}

		if v.Config.Saveflag {
			addval := Change{
				AddList: except(result[v.Listresult], origin[v.Listresult]),
				DelList: except(origin[v.Listresult], result[v.Listresult])}

			listid, ok := listids[v.Listresult]
			if !ok {
				listid = v.Listresult.ListID

				if listid == 0 {
					var mode string
					if v.Config.Publicflag {
						mode = "public"
					} else {
						mode = "private"
					}

					list, err := client.CreateList(twtr.Values{
						"name": v.Config.Name,
						"mode": mode,
					})
					if err != nil {
						return nil, err
					}

					listid = list.ID
				}
				listids[v.Listresult] = listid
			}
			ret[listid] = addval
		}
	}
	return ret, nil
}
