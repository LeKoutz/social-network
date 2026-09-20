package db

import (
	"errors"
	"forum/src/ferror"
	"forum/src/utils"

	"github.com/mattn/go-sqlite3"
)

type GroupRowType struct {
	Id            int64
	Timestamp     string
	Title         string
	Description   string
	OwnerUserId   int64
	OwnerUsername string
}

func (g *GroupRowType) InsertGroup() error {
	var err error
	query := `
		INSERT INTO groups (timestamp, title, description, owner_user_id)
		VALUES (?, ?, ?, ?)
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	res, err := stmt.Exec(
		utils.GetCurrentTimestamp(),
		g.Title,
		g.Description,
		g.OwnerUserId,
	)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				err = ferror.ErrorGroupTitleAlreadyExists
			}
		}
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	g.Id, err = res.LastInsertId()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (g *GroupRowType) SelectGroupById() error {
	var err error
	query := `
		SELECT groups.id, groups.timestamp, groups.title, groups.description, groups.owner_user_id, users.username
		FROM groups
		JOIN users ON groups.owner_user_id = users.id
		WHERE groups.id = ?
	`
	err = db.QueryRow(query, g.Id).Scan(
		&g.Id,
		&g.Timestamp,
		&g.Title,
		&g.Description,
		&g.OwnerUserId,
		&g.OwnerUsername,
	)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

// IsMember reports whether the given user belongs to the group, either as its
// owner or as a member with an accepted group invitation.
func (g *GroupRowType) IsMember(userId int64) (bool, error) {
	var err error
	var isOwner bool
	err = db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM groups WHERE id = ? AND owner_user_id = ?)`,
		g.Id,
		userId,
	).Scan(&isOwner)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return false, err
	}
	if isOwner {
		return true, nil
	}
	var invited bool
	err = db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM group_invitations
			WHERE group_id = ? AND to_user_id = ? AND status = 'accepted'
		)`,
		g.Id,
		userId,
	).Scan(&invited)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return false, err
	}
	return invited, nil
}
