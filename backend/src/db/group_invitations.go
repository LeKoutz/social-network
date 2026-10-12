package db

import (
	"database/sql"
	"errors"
	"forum/src/ferror"
	"forum/src/utils"
)

type GroupInvitationsRowType []GroupInvitationRowType

func (gi *GroupInvitationsRowType) SelectAllGroupInvites(group_id int64) error {
	rows, err := db.Query(`
	SELECT 
		g.id,
		g.timestamp,
		g.status,
		g.to_user_id,
		g.invited_by,
		u.username
	FROM group_invitations g
	JOIN users u ON g.to_user_id = u.id
	WHERE group_id = ?
	`, group_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ferror.ReturnErr(ferror.ErrorGroupNoInvites)
		}
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tmp GroupInvitationRowType
		err = rows.Scan(&tmp.Id, &tmp.Timestamp, &tmp.Status, &tmp.ToUserId, &tmp.InvitedBy, &tmp.Username)
		*gi = append(*gi, tmp)
	}
	if rows.Err() != nil {
		return ferror.ReturnErr(rows.Err())
	}
	return nil
}

func (gi *GroupInvitationsRowType) SelectGroupInvitesByUserId(group_id int64, user_id int64) error {
	rows, err := db.Query(`
	SELECT
		g.id,
		g.timestamp,
		g.status,
		g.to_user_id,
		g.invited_by,
		u.username
	FROM group_invitations g
	JOIN users u ON g.to_user_id = u.id
	WHERE group_id = ? AND invited_by = ?
	`, group_id, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ferror.ReturnErr(ferror.ErrorGroupNoInvites)
		}
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tmp GroupInvitationRowType
		err = rows.Scan(&tmp.Id, &tmp.Timestamp, &tmp.Status, &tmp.ToUserId, &tmp.InvitedBy, &tmp.Username)
		*gi = append(*gi, tmp)
	}
	if rows.Err() != nil {
		return ferror.ReturnErr(rows.Err())
	}
	return nil
}

func (gi *GroupInvitationsRowType) SelectAllGroupJoinRequestsByUserId(group_id int64) error {
	rows, err := db.Query(`
	SELECT
		g.id,
		g.timestamp,
		g.to_user_id,
		g.status,
		u.username
	FROM group_invitations g
	JOIN users u ON g.to_user_id = u.id
	WHERE group_id = ?
	`, group_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ferror.ReturnErr(ferror.ErrorGroupNoJoinRequests)
		}
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tmp GroupInvitationRowType
		err = rows.Scan(&tmp.Id, &tmp.Timestamp, &tmp.ToUserId, &tmp.Status, &tmp.Username)
		*gi = append(*gi, tmp)
	}
	if rows.Err() != nil {
		return ferror.ReturnErr(rows.Err())
	}
	return nil
}

func (gi *GroupInvitationsRowType) SelectGroupJoinRequestsByUserId(group_id int64, user_id int64) error {
	rows, err := db.Query(`
	SELECT
		g.id,
		g.timestamp,
		g.to_user_id,
		g.status,
		u.username
	FROM group_invitations g
	JOIN users u ON g.to_user_id = u.id
	WHERE group_id = ? AND invited_by = ?
	`, group_id, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ferror.ReturnErr(ferror.ErrorGroupNoJoinRequests)
		}
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tmp GroupInvitationRowType
		err = rows.Scan(&tmp.Id, &tmp.Timestamp, &tmp.ToUserId, &tmp.Status, &tmp.InvitedBy, &tmp.Username)
		*gi = append(*gi, tmp)
	}
	if rows.Err() != nil {
		return ferror.ReturnErr(rows.Err())
	}
	return nil
}

func (gi *GroupInvitationsRowType) SelectAllGroupsInvitesByUserId(user_id int64) error {
	rows, err := db.Query(`
	SELECT
		g.id,
		g.timestamp,
		g.status,
		g.group_id,
		g.to_user_id,
		g.invited_by,
		u.username,
		gr.title
	FROM group_invitations g
	JOIN users u ON g.invited_by = u.id
	JOIN groups gr ON g.group_id = gr.id
	WHERE to_user_id = ?
	`, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ferror.ReturnErr(ferror.ErrorGroupNoInvites)
		}
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tmp GroupInvitationRowType
		err = rows.Scan(&tmp.Id, &tmp.Timestamp, &tmp.Status, &tmp.GroupId, &tmp.ToUserId, &tmp.InvitedBy, &tmp.Username, &tmp.GroupTitle)
		*gi = append(*gi, tmp)
	}
	if rows.Err() != nil {
		return ferror.ReturnErr(rows.Err())
	}
	return nil
}

func (gi *GroupInvitationsRowType) SelectAllGroupsJoinRequestsByUserId(user_id int64) error {
	rows, err := db.Query(`
	SELECT
		g.id,
		g.timestamp,
		g.status,
		g.group_id,
		g.to_user_id,
		u.username,
		gr.title
	FROM group_invitations g
	JOIN users u ON g.to_user_id = u.id
	JOIN groups gr ON g.group_id = gr.id
	WHERE to_user_id = ?
	`, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ferror.ReturnErr(ferror.ErrorGroupNoInvites)
		}
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tmp GroupInvitationRowType
		err = rows.Scan(&tmp.Id, &tmp.Timestamp, &tmp.Status, &tmp.GroupId, &tmp.ToUserId, &tmp.Username, &tmp.GroupTitle)
		utils.LogDebug(tmp)
		if tmp.InvitedBy == 0 {
			*gi = append(*gi, tmp)
		}
	}
	if rows.Err() != nil {
		return ferror.ReturnErr(rows.Err())
	}
	return nil
}
