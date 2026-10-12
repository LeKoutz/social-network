package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type GroupInvitationsType []GroupInvitationType

func (g *GroupInvitationsType) GetAllGroupInvitations(group_id int64) error {
	var err error
	rows := &db.GroupInvitationsRowType{}
	err = rows.SelectAllGroupInvites(group_id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if g == nil {
		g = &GroupInvitationsType{}
	}
	for _, row := range *rows {
		var tmp GroupInvitationType
		tmp.GroupInvitationRowType = row
		*g = append(*g, tmp)
	}
	return nil
}

func (g *GroupInvitationsType) GetGroupInvitatesByUserId(group_id int64, user_id int64) error {
	var err error
	rows := &db.GroupInvitationsRowType{}
	err = rows.SelectGroupInvitesByUserId(group_id, user_id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if g == nil {
		g = &GroupInvitationsType{}
	}
	for _, row := range *rows {
		var tmp GroupInvitationType
		tmp.GroupInvitationRowType = row
		*g = append(*g, tmp)
	}
	return nil
}

func (g *GroupInvitationsType) GetGroupJoinRequestsByGroupId(group_id int64) error {
	var err error
	rows := &db.GroupInvitationsRowType{}
	err = rows.SelectAllGroupInvites(group_id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if g == nil {
		g = &GroupInvitationsType{}
	}
	for _, row := range *rows {
		var tmp GroupInvitationType
		tmp.GroupInvitationRowType = row
		*g = append(*g, tmp)
	}
	return nil
}

func (g *GroupInvitationsType) GetGroupJoinRequestsByUserId(group_id int64, user_id int64) error {
	var err error
	rows := &db.GroupInvitationsRowType{}
	err = rows.SelectGroupJoinRequestsByUserId(group_id, user_id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if g == nil {
		g = &GroupInvitationsType{}
	}
	for _, row := range *rows {
		var tmp GroupInvitationType
		tmp.GroupInvitationRowType = row
		*g = append(*g, tmp)
	}
	return nil
}

func (g *GroupInvitationsType) GetAllGroupsInvitesByUserId(user_id int64) error {
	var err error
	rows := &db.GroupInvitationsRowType{}
	err = rows.SelectAllGroupsInvitesByUserId(user_id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if g == nil {
		g = &GroupInvitationsType{}
	}
	for _, row := range *rows {
		var tmp GroupInvitationType
		tmp.GroupInvitationRowType = row
		*g = append(*g, tmp)
	}
	return nil
}

func (g *GroupInvitationsType) GetAllGroupsJoinRequestsByUserId(user_id int64) error {
	var err error
	rows := &db.GroupInvitationsRowType{}
	err = rows.SelectAllGroupsJoinRequestsByUserId(user_id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if g == nil {
		g = &GroupInvitationsType{}
	}
	for _, row := range *rows {
		var tmp GroupInvitationType
		tmp.GroupInvitationRowType = row
		*g = append(*g, tmp)
	}
	return nil
}
