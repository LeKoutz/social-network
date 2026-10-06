package db

import (
	"forum/src/ferror"
)

type GroupInvitationRowType struct {
	Id        int64
	Timestamp string
	Status    string
	GroupId   int64
	ToUserId  int64
	InvitedBy int64
	Username  string
}

func (g *GroupInvitationRowType) Insert() error {
	return ferror.ReturnErr(ferror.ErrorNotImplemented)
}

func (g *GroupInvitationRowType) Update() error {
	return ferror.ReturnErr(ferror.ErrorNotImplemented)
}

func (g *GroupInvitationRowType) InsertOrUpdate() error {
	rows, err := db.Query(`
	SELECT id, to_user_id
	FROM group_invitations
	WHERE group_id = ? AND to_user_id = ?
	`, g.GroupId, g.ToUserId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	var tmp GroupInvitationsRowType
	for rows.Next() {
		var t GroupInvitationRowType
		err = rows.Scan(&t.Id, &t.InvitedBy)
		tmp = append(tmp, t)
	}
	for i := range tmp {
		if tmp[i].InvitedBy == g.InvitedBy &&
			tmp[i].Status != "pending" &&
			tmp[i].Status != "accepted" {
			_, err = db.Exec(`
			UPDATE group_invitations
			SET status = 'pending'
			WHERE id = ?
			`, tmp[i].Id)
			if err != nil {
				return ferror.ReturnErr(err)
			}
			return nil
		}
	}
	if len(tmp) != 0 {
		return ferror.ReturnErr(ferror.ErrorGroupInvitationAlreadyExists)
	}
	_, err = db.Exec(`
	INSERT INTO group_invitations
	(group_id, invited_by, timestamp, to_user_id, status)
	VALUES (?, ?, ?, ?, ?)
	`, g.GroupId, g.InvitedBy, g.Timestamp, g.ToUserId, g.Status)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
