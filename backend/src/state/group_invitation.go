package state

import "forum/src/models"

func (r *State) EditGroupInvitation() *models.GroupInvitationType {
	return &r.GroupInvitation
}
