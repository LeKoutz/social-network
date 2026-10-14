package models

import (
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"
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
	f.IncomingInvites = make([]InvitationType, 0)
	f.OutgoingInvites = make([]InvitationType, 0)
	f.PendingIncoming = make([]InvitationType, 0)
	f.PendingOutgoing = make([]InvitationType, 0)
	var invites db.InvitationsRowType
	err = invites.SelectInvitationsToUserId(userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range invites {
		f.IncomingInvites = append(f.IncomingInvites, InvitationType{InvitationRowType: row})
	}
	invites = nil
	utils.LogInfo("iv to id")

	err = invites.SelectAcceptedInvitationsFromUserId(userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range invites {
		f.OutgoingInvites = append(f.OutgoingInvites, InvitationType{InvitationRowType: row})
	}
	invites = nil
	utils.LogInfo("acc to id")

	err = invites.SelectPendingInvitesByUserId(userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range invites {
		follow := InvitationType{InvitationRowType: row}
		if follow.FromUserId == userId {
			f.PendingOutgoing = append(f.PendingOutgoing, follow)
		} else {
			f.PendingIncoming = append(f.PendingIncoming, follow)
		}
	}
	invites = nil
	utils.LogDebug(invites)
	return nil
}
