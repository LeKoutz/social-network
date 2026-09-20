package db

import (
	"errors"
	"forum/src/utils"
)

type GroupRowsType []GroupRowType

func (groups *GroupRowsType) SelectAllGroups() error {
	rows, err := db.Query(`
		SELECT groups.id, groups.timestamp, groups.title, groups.description, groups.owner_user_id, users.username
		FROM groups
		JOIN users ON groups.owner_user_id = users.id
	`)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var group GroupRowType
		err = rows.Scan(
			&group.Id,
			&group.Timestamp,
			&group.Title,
			&group.Description,
			&group.OwnerUserId,
			&group.OwnerUsername,
		)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		*groups = append(*groups, group)
	}
	return nil
}
