package db

import "forum/src/ferror"

type GroupRowsType []GroupRowType

func (groups *GroupRowsType) SelectAllGroups() error {
	rows, err := db.Query(`
		SELECT groups.id, groups.timestamp, groups.title, groups.description, groups.owner_user_id, users.username
		FROM groups
		JOIN users ON groups.owner_user_id = users.id
	`)
	if err != nil {
		return ferror.ReturnErr(err)
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
			return ferror.ReturnErr(err)
		}
		*groups = append(*groups, group)
	}
	return nil
}
