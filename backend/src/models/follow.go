package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

const (
	FollowStatusPending    = "pending"
	FollowStatusAccepted   = "accepted"
	FollowStatusDenied     = "denied"
	FollowStatusUnfollowed = "unfollowed"
)

type FollowType struct {
	db.InvitationRowType
}

type FollowsType struct {
	Followers       []FollowType
	Following       []FollowType
	PendingIncoming []FollowType
	PendingOutgoing []FollowType
}

func (f *FollowsType) GetForUser(userId int64) error {
	var err error
	rows, err := db.SelectFollowersByUserId(userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range rows {
		f.Followers = append(f.Followers, FollowType{InvitationRowType: row})
	}

	rows, err = db.SelectFollowingByUserId(userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range rows {
		f.Following = append(f.Following, FollowType{InvitationRowType: row})
	}

	for _, row := range rows {
		follow := FollowType{InvitationRowType: row}
		if follow.FromUserId == userId {
			f.PendingOutgoing = append(f.PendingOutgoing, follow)
		} else {
			f.PendingIncoming = append(f.PendingIncoming, follow)
		}
	}
	return nil
}

type FollowRequestType struct {
	db.InvitationRowType
}

func (f *FollowRequestType) AddFollowInvite() error {
	return f.InvitationRowType.Insert()
}

func (f *FollowRequestType) Accept(userId int64) error {
	return f.respondToRequest(userId, FollowStatusAccepted)
}

func (f *FollowRequestType) Deny(userId int64) error {
	return f.respondToRequest(userId, FollowStatusDenied)
}

func (f *FollowRequestType) Unfollow() error {
	return f.InvitationRowType.Unfollow()
}

func (f *FollowRequestType) respondToRequest(userId int64, status string) error {
	var err error
	err = f.InvitationRowType.SelectInvitationById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if f.ToUserId != userId {
		err = ferror.ErrorPermissionDenied
		return ferror.ReturnErr(err)
	}
	if f.Status != FollowStatusPending {
		err = ferror.ErrorInvitationNotPending
		return ferror.ReturnErr(err)
	}
	f.Status = status
	err = f.InvitationRowType.UpdateInvitationStatus()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
