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
