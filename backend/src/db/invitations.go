package db

import (
	"database/sql"
	"forum/src/ferror"
)

type InvitationsRowType []InvitationRowType

func (irt *InvitationsRowType) SelectAllReferencingUserId(user_id int64) error {
	if irt == nil {
		irt = &InvitationsRowType{}
	}
	var err error
	query := `
	SELECT i.id, i.status, i.timestamp, i.from_user_id, i.to_user_id, u.username
	FROM invitations i
	JOIN users u ON u.id = CASE WHEN i.from_user_id = ? THEN i.to_user_id ELSE i.from_user_id END
	WHERE from_user_id = ? OR to_user_id = ?
	`
	rows, err := db.Query(query, user_id, user_id, user_id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var inv InvitationRowType
		err = rows.Scan(&inv.Id, &inv.Status, &inv.Timestamp, &inv.FromUserId, &inv.ToUserId, &inv.Username)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		*irt = append(*irt, inv)
	}
	return nil
}

func (irt *InvitationsRowType) SelectInvitationsFromUserId(user_id int64) error {
	if irt == nil {
		irt = &InvitationsRowType{}
	}
	var err error
	query := `
	SELECT i.id, i.status, i.timestamp, i.to_user_id, u.username
	FROM invitations i
	JOIN users u ON u.id = i.to_user_id
	WHERE from_user_id = ?
	`
	rows, err := db.Query(query, user_id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var inv InvitationRowType
		err = rows.Scan(&inv.Id, &inv.Status, &inv.Timestamp, &inv.ToUserId, &inv.Username)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		*irt = append(*irt, inv)
	}
	return nil
}

func (inv *InvitationsRowType) SelectPendingInvitesByUserId(userId int64) error {
	rows, err := db.Query(`
	SELECT
		i.id,
		i.timestamp,
		i.status,
		i.from_user_id,
		i.to_user_id,
		u.username
	FROM invitations i
	JOIN users u ON u.id = CASE WHEN i.from_user_id = ? THEN i.to_user_id ELSE i.from_user_id END
	WHERE i.status = 'pending' AND i.from_user_id = ?
	`, userId, userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	*inv, err = scanInvitations(rows)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (inv *InvitationsRowType) SelectPendingInvitesToUserId(userId int64) error {
	rows, err := db.Query(`
	SELECT
		i.id,
		i.timestamp,
		i.status,
		i.from_user_id,
		i.to_user_id,
		u.username
	FROM invitations i
	JOIN users u ON u.id = CASE WHEN i.from_user_id = ? THEN i.to_user_id ELSE i.from_user_id END
	WHERE i.status = 'pending' AND i.to_user_id = ?
	`, userId, userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	*inv, err = scanInvitations(rows)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (inv *InvitationsRowType) SelectAcceptedInvitationsToUserId(userId int64) error {
	rows, err := db.Query(`
	SELECT
		i.id,
		i.timestamp,
		i.status,
		i.from_user_id,
		i.to_user_id,
		u.username
	FROM invitations i
	JOIN users u ON u.id = CASE WHEN i.from_user_id = ? THEN i.to_user_id ELSE i.from_user_id END
	WHERE i.status = 'accepted' AND i.to_user_id = ?
	`, userId, userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	*inv, err = scanInvitations(rows)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (inv *InvitationsRowType) SelectAcceptedInvitationsFromUserId(userId int64) error {
	rows, err := db.Query(`
	SELECT
		i.id,
		i.timestamp,
		i.status,
		i.from_user_id,
		i.to_user_id,
		u.username
	FROM invitations i
	JOIN users u ON u.id = CASE WHEN i.from_user_id = ? THEN i.to_user_id ELSE i.from_user_id END
	WHERE i.status = 'accepted' AND i.from_user_id = ?
	`, userId, userId, userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	*inv, err = scanInvitations(rows)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (irt *InvitationsRowType) SelectInvitationsToUserId(user_id int64) error {
	if irt == nil {
		irt = &InvitationsRowType{}
	}
	var err error
	query := `
	SELECT i.id, i.status, i.timestamp, i.from_user_id, u.username, to_user_id
	FROM invitations i
	JOIN users u ON u.id = i.from_user_id
	WHERE to_user_id = ?
	`
	rows, err := db.Query(query, user_id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	defer rows.Close()
	for rows.Next() {
		var inv InvitationRowType
		err = rows.Scan(&inv.Id, &inv.Status, &inv.Timestamp, &inv.FromUserId, &inv.Username, &inv.ToUserId)
		if err != nil {
			return ferror.ReturnErr(err)
		}
		*irt = append(*irt, inv)
	}
	return nil
}

func scanInvitations(rows *sql.Rows) (InvitationsRowType, error) {
	var invitations []InvitationRowType
	for rows.Next() {
		var invitation InvitationRowType
		err := rows.Scan(
			&invitation.Id,
			&invitation.Timestamp,
			&invitation.Status,
			&invitation.FromUserId,
			&invitation.ToUserId,
			&invitation.Username,
		)
		if err != nil {
			return InvitationsRowType{}, ferror.ReturnErr(err)
		}
		invitations = append(invitations, invitation)
	}
	return invitations, nil
}

func SelectFollowersByUserId(userId int64) ([]InvitationRowType, error) {
	rows, err := db.Query(`
	SELECT
		i.id,
		i.timestamp,
		i.status,
		i.from_user_id,
		i.to_user_id,
		u.username
	FROM invitations i
	JOIN users u ON u.id = i.from_user_id
	WHERE i.to_user_id = ? AND i.status = 'accepted'
	`, userId)
	if err != nil {
		return []InvitationRowType{}, ferror.ReturnErr(err)
	}
	defer rows.Close()
	return scanInvitations(rows)
}

func SelectFollowingByUserId(userId int64) ([]InvitationRowType, error) {
	rows, err := db.Query(`
	SELECT
		i.id,
		i.timestamp,
		i.status,
		i.from_user_id,
		i.to_user_id,
		u.username
	FROM invitations i
	JOIN users u ON u.id = i.to_user_id
	WHERE i.from_user_id = ? AND i.status = 'accepted'
	`, userId)
	if err != nil {
		return []InvitationRowType{}, ferror.ReturnErr(err)
	}
	defer rows.Close()
	return scanInvitations(rows)
}
