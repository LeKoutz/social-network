package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type GroupsType []GroupType

func (g *GroupsType) GetGroups() error {
	var err error
	var groups GroupsType
	var rows db.GroupRowsType
	err = rows.SelectAllGroups()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	for _, row := range rows {
		var group GroupType
		group.GroupRowType = row
		groups = append(groups, group)
	}
	*g = groups
	return nil
}
