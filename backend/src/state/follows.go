package state

import "forum/src/models"

func (r *State) EditFollows() *models.InvitationsType {
	return &r.Follows
}
