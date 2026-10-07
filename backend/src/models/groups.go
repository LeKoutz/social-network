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

func (g *GroupsType) GetGroupsWithUnreadMessages(userId int64) error {
	groupRows, messageRows, err := db.SelectGroupsWithUnreadMessages(userId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	var groups GroupsType
	for i, row := range groupRows {
		if len(groups) == 0 || groups[len(groups)-1].Id != row.Id {
			var group GroupType
			group.GroupRowType = row
			group.Member = true
			groups = append(groups, group)
		}
		if messageRows[i].Id == 0 {
			continue
		}
		current := &groups[len(groups)-1]
		var message GroupMessageType
		message.GroupMessageRowType = messageRows[i]
		current.ChatMessages = append(current.ChatMessages, message)
	}
	*g = groups
	return nil
}
