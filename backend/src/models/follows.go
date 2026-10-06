package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type FollowRequestsType []FollowRequestType

func (f *FollowRequestsType) GetFollowers(user_id int64) error {
	x, err := db.SelectFollowersByUserId(user_id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if f == nil {
		f = &FollowRequestsType{}
	}
	for i := range x {
		*f = append(*f, FollowRequestType{x[i]})
	}
	return nil
}
