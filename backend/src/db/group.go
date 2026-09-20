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
