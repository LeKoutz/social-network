package state

import (
	"forum/src/db"
	"forum/src/models"
)

func (r *State) InitFollowRequest() {
	r.SetFollowRequest(models.FollowRequestType{
		InvitationRowType: db.InvitationRowType{FromUserId: r.User.Id},
	})
}

func (r *State) SetFollowRequest(followRequest models.FollowRequestType) {
	r.FollowRequest = followRequest
}

func (r *State) GetFollowRequest() models.FollowRequestType {
	return *r.EditFollowRequest()
}

func (r *State) EditFollowRequest() *models.FollowRequestType {
	if r.FollowRequest.FromUserId == 0 {
		r.InitFollowRequest()
	}
	return &r.FollowRequest
}
