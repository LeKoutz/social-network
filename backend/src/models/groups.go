package models

import (
	"errors"
	"forum/src/db"
	"forum/src/utils"
)

type GroupsType []GroupType

func (g *GroupsType) GetGroups() error {
	var err error
	var groups GroupsType
	var rows db.GroupRowsType
	err = rows.SelectAllGroups()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for _, row := range rows {
		var group GroupType
		group.GroupRowType = row
		groups = append(groups, group)
	}
	*g = groups
	return nil
}
