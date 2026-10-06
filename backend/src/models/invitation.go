package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

const (
	InvitationStatusPending    = "pending"
	InvitationStatusAccepted   = "accepted"
	InvitationStatusDenied     = "denied"
	InvitationStatusUnfollowed = "unfollowed"
)

type InvitationType struct {
	db.InvitationRowType
}

type InvitationRequestType struct {
	db.InvitationRowType
}

func (f *InvitationRequestType) AddInvitationInvite() error {
	return f.InvitationRowType.Insert()
}

func (f *InvitationRequestType) Accept(userId int64) error {
	return f.respondToRequest(userId, InvitationStatusAccepted)
}

func (f *InvitationRequestType) Deny(userId int64) error {
	return f.respondToRequest(userId, InvitationStatusDenied)
}

func (f *InvitationRequestType) Unfollow() error {
	return f.InvitationRowType.Unfollow()
}

func (f *InvitationRequestType) respondToRequest(userId int64, status string) error {
	var err error
	err = f.InvitationRowType.SelectInvitationById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if f.ToUserId != userId {
		err = ferror.ErrorPermissionDenied
		return ferror.ReturnErr(err)
	}
	if f.Status != InvitationStatusPending {
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
