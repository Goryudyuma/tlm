package follower

import (
	"strconv"

	"github.com/Goryudyuma/tlm/lib/User"
	"github.com/bgpat/twtr"
)

func (f Follower) GetFollowerIDs(client *twtr.Client) (user.UserIDs, error) {
	userID := f.UserID
	resp, _, err := client.GetFollowerIDs(&twtr.Params{"user_id": strconv.FormatInt(int64(userID), 10)})
	if err != nil {
		return nil, err
	}
	var ret user.UserIDs
	for _, v := range resp.IDs {
		ret = append(ret, user.UserID(v))
	}
	for resp.Cursor.NextCursor != 0 {
		resp, _, err = client.GetFollowerIDs((&twtr.Params{"user_id": strconv.FormatInt(int64(userID), 10)}).AddNextCursor(resp.Cursor))
		if err != nil {
			return nil, err
		}
		for _, v := range resp.IDs {
			ret = append(ret, user.UserID(v))
		}
	}
	return ret, nil
}
