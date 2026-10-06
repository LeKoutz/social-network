package state

import "forum/src/models"

func (r *State) EditGroupInvitations() *models.GroupInvitationsType {
	return &r.GroupInvitations
}
