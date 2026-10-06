package db

import (
	"database/sql"
	"errors"
	"forum/src/ferror"
	"forum/src/utils"

	"github.com/mattn/go-sqlite3"
)

type InvitationRowType struct {
	Id         int64
	Timestamp  string
	Status     string
	FromUserId int64
	ToUserId   int64
	Username   string
}

func (i *InvitationRowType) Insert() error {
	var err error
	existing, err := SelectInvitationByUsers(i.FromUserId, i.ToUserId)
	if err == nil {
		i.Id = existing.Id
		switch existing.Status {
		case "accepted", "pending":
			err = ferror.ErrorInvitationAlreadyExists
			return ferror.ReturnErr(err)
		case "denied", "unfollowed":
			return i.reactivateInvitation()
		}
	} else if !errors.Is(err, ferror.ErrorInvitationNotFound) {
		return ferror.ReturnErr(err)
	}
	query := `
	INSERT INTO invitations (status, timestamp, from_user_id, to_user_id)
	VALUES (
		COALESCE(
			CASE
			WHEN (SELECT private_profile FROM user_profiles WHERE user_id = ?) = 0
			THEN 'accepted'
			END,
			'pending'
		),
		?, ?, ?
	);
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	res, err := stmt.Exec(i.ToUserId, utils.GetCurrentTimestamp(), i.FromUserId, i.ToUserId)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				err = ferror.ErrorInvitationAlreadyExists
			}
		}
		return ferror.ReturnErr(err)
	}
	i.Id, err = res.LastInsertId()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (i *InvitationRowType) Unfollow() error {
	_, err := db.Exec(
		"UPDATE invitations SET status = 'unfollowed' WHERE from_user_id = ? AND to_user_id = ? AND status = 'accepted'",
		i.FromUserId, i.ToUserId,
	)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func (i *InvitationRowType) AreFollowingEachOther(fromUserId, toUserId int64) (bool, error) {
	var following bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM invitations
			WHERE status = 'accepted'
			AND (
				(from_user_id = ? AND to_user_id = ?)
				OR (from_user_id = ? AND to_user_id = ?)
			)
		)
	`, fromUserId, toUserId, toUserId, fromUserId).Scan(&following)
	if err != nil {
		return false, ferror.ReturnErr(err)
	}
	return following, nil
}

func (i *InvitationRowType) SelectInvitationById() error {
	err := db.QueryRow(`
	SELECT id, timestamp, status, from_user_id, to_user_id
	FROM invitations
	WHERE id = ?
	`, i.Id).Scan(&i.Id, &i.Timestamp, &i.Status, &i.FromUserId, &i.ToUserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ferror.ErrorInvitationNotFound
		}
		return ferror.ReturnErr(err)
	}
	return nil
}

func SelectInvitationByUsers(fromUserId, toUserId int64) (InvitationRowType, error) {
	var invitation InvitationRowType
	err := db.QueryRow(`
	SELECT id, timestamp, status, from_user_id, to_user_id
	FROM invitations
	WHERE from_user_id = ? AND to_user_id = ?
	`, fromUserId, toUserId).Scan(
		&invitation.Id,
		&invitation.Timestamp,
		&invitation.Status,
		&invitation.FromUserId,
		&invitation.ToUserId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ferror.ErrorInvitationNotFound
		}
		return InvitationRowType{}, ferror.ReturnErr(err)
	}
	return invitation, nil
}

func (i *InvitationRowType) reactivateInvitation() error {
	privateProfile := true
	err := db.QueryRow(`SELECT private_profile FROM user_profiles WHERE user_id = ?`, i.ToUserId).Scan(&privateProfile)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return ferror.ReturnErr(err)
	}
	if !privateProfile {
		i.Status = "accepted"
	} else {
		i.Status = "pending"
	}
	i.Timestamp = utils.GetCurrentTimestamp()
	return i.UpdateInvitationStatus()
}

func (i *InvitationRowType) UpdateInvitationStatus() error {
	_, err := db.Exec("UPDATE invitations SET status = ? WHERE id = ?", i.Status, i.Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
