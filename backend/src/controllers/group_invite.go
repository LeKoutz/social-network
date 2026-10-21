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

func AcceptGroupInvitation(data state.StateController) error {
	return data.EditGroupInvitation().UpdateToAccepted()
}

func DeclineGroupInvitation(data state.StateController) error {
	return data.EditGroupInvitation().UpdateToDenied()
}

func AcceptJoinRequest(data state.StateController) error {
	err := data.EditGroupInvitation().SelectGroupInvitationById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	err = data.EditGroup().SelectGroupById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if data.EditGroup().OwnerUserId == data.GetUser().Id {
		return data.EditGroupInvitation().JoinRequestAccept()
	}
	return ferror.ReturnErr(ferror.ErrorPermissionDenied)
}

func DeclineJoinRequest(data state.StateController) error {
	return ferror.ReturnErr(ferror.ErrorNotImplemented)
}

func WithdrawJoinRequest(data state.StateController) error {
	err := data.EditGroupInvitation().Withdraw()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func RequestJoinGroup(data state.StateController) error {
	data.EditGroup().Id = data.EditGroupInvitation().GroupId
	err := data.EditGroup().GetById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.EditGroup().Member, err = data.EditGroup().IsMember(data.GetUser().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if data.EditGroup().OwnerUserId == data.GetUser().Id {
		data.EditGroup().Owner = true
	}
	if data.EditGroup().Member || data.EditGroup().Owner {
		return ferror.ReturnErr(ferror.ErrorGroupAlreadyMember)
	}
	data.EditGroupInvitation().ToUserId = data.GetUser().Id
	data.EditGroupInvitation().InvitedBy = 0
	data.EditGroupInvitation().Status = "pending"
	return data.EditGroupInvitation().InsertOrUpdate()
}
