package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type InvitationsType struct {
	db.InvitationsRowType `json:"-"`

	IncomingInvites []InvitationType `json:"Followers,omitzero"`
	OutgoingInvites []InvitationType `json:"Following,omitzero"`
}

func (f *InvitationsType) GetForUser(userId int64) error {
	var err error
	f.IncomingInvites = make([]InvitationType, 0)
	f.OutgoingInvites = make([]InvitationType, 0)
	var invites db.InvitationsRowType
	err = invites.SelectAllReferencingUserId(userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range invites {
		switch {
		case userId == row.ToUserId:
			f.IncomingInvites = append(f.IncomingInvites, InvitationType{InvitationRowType: row})
		case userId == row.FromUserId:
			f.OutgoingInvites = append(f.OutgoingInvites, InvitationType{InvitationRowType: row})
		}
	}
	invites = nil
	return nil
}
