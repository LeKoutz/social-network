package db

import (
	"database/sql"
	"errors"
	"forum/src/ferror"
)

type GroupInvitationRowType struct {
	Id         int64
	Timestamp  string
	Status     string
	GroupId    int64
	ToUserId   int64
	InvitedBy  int64
	Username   string
	GroupTitle string
}

func (g *GroupInvitationRowType) Insert() error {
	var err error
	value := func(i int64) any {
		if i == 0 {
			return nil
		} else {
			return i
		}
	}(g.InvitedBy)

	_, err = db.Exec(`
	INSERT INTO group_invitations
	(group_id, invited_by, timestamp, to_user_id, status)
	VALUES (?, ?, ?, ?, ?)
	`, g.GroupId, value, g.Timestamp, g.ToUserId, g.Status)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (g *GroupInvitationRowType) Update() error {
	return ferror.ReturnErr(ferror.ErrorNotImplemented)
}

func (g *GroupInvitationRowType) UpdateToPending() error {
	_, err := db.Exec(`
	UPDATE group_invitations
	SET status = 'pending'
	WHERE id = ?
	`, g.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (g *GroupInvitationRowType) UpdateToAccepted() error {
	_, err := db.Exec(`
	UPDATE group_invitations
	SET status = 'accepted'
	WHERE id = ?
	`, g.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (g *GroupInvitationRowType) UpdateToDenied() error {
	_, err := db.Exec(`
	UPDATE group_invitations
	SET status = 'denied'
	WHERE id = ?
	`, g.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (g *GroupInvitationRowType) UpdateToRemoved() error {
	_, err := db.Exec(`
	UPDATE group_invitations
	SET status = 'removed'
	WHERE id = ?
	`, g.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (g *GroupInvitationRowType) UpdateToLeft() error {
	_, err := db.Exec(`
	UPDATE group_invitations
	SET status = 'left'
	WHERE id = ?
	`, g.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (g *GroupInvitationRowType) InsertOrUpdate() error {
	var err error
	var tmp GroupInvitationsRowType
	err = tmp.SelectGroupInvitesByUserId(g.GroupId, g.ToUserId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for i := range tmp {
		if tmp[i].InvitedBy == g.InvitedBy &&
			tmp[i].Status != "pending" &&
			tmp[i].Status != "accepted" {
			return g.UpdateToPending()
		}
	}
	if len(tmp) != 0 {
		return ferror.ReturnErr(ferror.ErrorGroupInvitationAlreadyExists)
	}
	return g.Insert()
}

func (gi *GroupInvitationRowType) SelectGroupInvitationById() error {
	rows, err := db.Query(`
	SELECT
	g.id,
	g.timestamp,
	g.status,
	g.group_id,
	g.invited_by,
	g.to_user_id,
	u.username,
	gr.title
	FROM group_invitations g
	JOIN users u ON g.to_user_id = u.id
	JOIN groups gr ON g.group_id = gr.id
	WHERE id = ?
	`, gi.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ferror.ReturnErr(ferror.ErrorGroupNoInvites)
		}
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&gi.Id, &gi.Timestamp, &gi.Status, &gi.GroupId, &gi.InvitedBy, &gi.ToUserId, &gi.Username, &gi.GroupTitle)
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	return nil
}

func (gi *GroupInvitationRowType) Delete() error {
	_, err := db.Exec(`
		DELETE FROM group_invitations
		WHERE id = ?
		`, gi.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
