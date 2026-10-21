package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type GroupInvitationType struct {
	db.GroupInvitationRowType
}

func (gi *GroupInvitationType) SendToUserId(user_id int64) error {
	return ferror.ReturnErr(ferror.ErrorNotImplemented)
}

func (gi *GroupInvitationType) Withdraw() error {
	return gi.GroupInvitationRowType.Delete()
}

func (gi *GroupInvitationType) JoinRequestAccept() error {
	return gi.GroupInvitationRowType.UpdateToAccepted()
}
