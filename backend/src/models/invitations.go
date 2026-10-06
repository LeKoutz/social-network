package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type InvitationsType struct {
	db.InvitationsRowType `json:"-"`

	IncomingInvites []InvitationType `json:"Followers,omitzero"`
	OutgoingInvites []InvitationType `json:"Following,omitzero"`
	PendingIncoming []InvitationType `json:"PendingIncoming,omitzero"`
	PendingOutgoing []InvitationType `json:"PendingOutgoing,omitzero"`
}

func (f *InvitationsType) GetForUser(userId int64) error {
	var err error
	err = f.InvitationsRowType.SelectInvitationsToUserId(userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range f.InvitationsRowType {
		f.IncomingInvites = append(f.IncomingInvites, InvitationType{InvitationRowType: row})
	}
	f.InvitationsRowType = nil

	err = f.InvitationsRowType.SelectAcceptedInvitationsFromUserId(userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range f.InvitationsRowType {
		f.OutgoingInvites = append(f.OutgoingInvites, InvitationType{InvitationRowType: row})
	}
	f.InvitationsRowType = nil

	err = f.InvitationsRowType.SelectPendingInvitesByUserId(userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range f.InvitationsRowType {
		follow := InvitationType{InvitationRowType: row}
		if follow.FromUserId == userId {
			f.PendingOutgoing = append(f.PendingOutgoing, follow)
		} else {
			f.PendingIncoming = append(f.PendingIncoming, follow)
		}
	}
	f.InvitationsRowType = nil
	return nil
}
