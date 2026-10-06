package controllers

import (
	"forum/src/ferror"
	"forum/src/state"
)

func ShowGroupInvites(data state.StateController) error {
	err := data.EditGroup().GetById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	b, err := data.EditGroup().IsMember(data.GetUser().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if !b {
		return ferror.ReturnErr(ferror.ErrorPermissionDenied)
	}
	if data.EditGroup().OwnerUserId == data.GetUser().Id {
		return data.EditGroupInvitations().GetAllGroupInvitations(data.GetGroup().Id)
	}
	return data.EditGroupInvitations().GetGroupInvitatesByUserId(data.GetGroup().Id, data.GetUser().Id)
}

func ShowGroupJoins(data state.StateController) error {
	err := data.EditGroup().GetById()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	b, err := data.EditGroup().IsMember(data.GetUser().Id)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if !b {
		return ferror.ReturnErr(ferror.ErrorPermissionDenied)
	}
	if data.EditGroup().OwnerUserId == data.GetUser().Id {
		return data.EditGroupInvitations().GetAllGroupInvitations(data.GetGroup().Id)
	}
	return data.EditGroupInvitations().GetGroupInvitatesByUserId(data.GetGroup().Id, data.GetUser().Id)
}
