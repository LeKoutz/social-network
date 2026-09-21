package models

import "forum/src/db"

type FollowRequestType struct {
	db.InvitationRowType
}

func (f *FollowRequestType) AddFollowInvite() error {
	return f.InvitationRowType.Insert()
}

func (f *FollowRequestType) Unfollow() error {
    return f.InvitationRowType.Unfollow()
}