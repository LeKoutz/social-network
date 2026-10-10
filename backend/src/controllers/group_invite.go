package controllers

import (
	"forum/src/ferror"
	"forum/src/state"
)

func InviteUserToGroup(data state.StateController) error {
	data.EditGroup().Id = data.EditGroupInvitation().GroupId
	err := data.EditGroup().GetById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	b, err := data.EditGroup().IsMember(data.GetUser().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if !b && data.EditGroup().OwnerUserId != data.GetUser().Id {
		return ferror.ReturnErr(ferror.ErrorPermissionDenied)
	}
	b, err = data.EditGroup().IsMember(data.EditGroupInvitation().ToUserId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if b {
		return ferror.ReturnErr(ferror.ErrorGroupAlreadyMember)
	}
	data.EditGroupInvitation().InvitedBy = data.GetUser().Id
	data.EditGroupInvitation().Status = "pending"
	return data.EditGroupInvitation().InsertOrUpdate()
}

func RequestJoinGroup(data state.StateController) error {
	data.EditGroup().Id = data.EditGroupInvitation().GroupId
	err := data.EditGroup().GetById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if data.EditGroup().Member || data.EditGroup().Owner {
		return ferror.ReturnErr(ferror.ErrorGroupAlreadyMember)
	}
	data.EditGroupInvitation().ToUserId = data.GetUser().Id
	data.EditGroupInvitation().InvitedBy = 0
	data.EditGroupInvitation().Status = "pending"
	return data.EditGroupInvitation().InsertOrUpdate()
}
