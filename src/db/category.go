package db

import (
	"errors"
	"forum/src/ferror"
	"forum/src/utils"

	"github.com/mattn/go-sqlite3"
)

type CategoryRowType struct {
	Id          int64
	Name        string
	Description string
}

func (cr *CategoryRowType) SelectCategoryById() error {
	var err error
	query := `SELECT name, description FROM categories WHERE id = ?`
	err = db.QueryRow(query, cr.Id).Scan(&cr.Name, &cr.Description)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}

func (cr *CategoryRowType) InsertCategory() error {
	var err error
	query := "INSERT INTO categories (name, description) VALUES (?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	res, err := stmt.Exec(cr.Name, cr.Description)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return ferror.ErrorCategoryAlreadyExists
			}
		}
		return errors.Join(utils.GetFunctionName(), err)
	}
	cr.Id, err = res.LastInsertId()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return nil
}
