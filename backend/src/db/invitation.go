package db

import (
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
}

func (i *InvitationRowType) Insert() error {
	var err error
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
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	res, err := stmt.Exec(utils.GetCurrentTimestamp(), i.FromUserId, i.ToUserId)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				err = ferror.ErrorInvitationAlreadyExists
			}
		}
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	i.Id, err = res.LastInsertId()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (i *InvitationRowType) Unfollow() error {
    _, err := db.Exec(
        "UPDATE invitations SET status = 'unfollowed' WHERE from_user_id = ? AND to_user_id = ? AND status = 'accepted'",
        i.FromUserId, i.ToUserId,
    )
    if err != nil {
        if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
        return err
    }
    return nil
}