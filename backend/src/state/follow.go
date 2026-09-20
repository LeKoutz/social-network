package state

import "forum/src/models"

func (r *State) InitFollowRequest() {
	var followRequest models.FollowRequestType
	followRequest.FromUserId = r.User.Id
	r.SetFollowRequest(followRequest)
}

func (r *State) SetFollowRequest(followRequest models.FollowRequestType) {
	*r.EditFollowRequest() = followRequest
}

func (r *State) GetFollowRequest() models.FollowRequestType {
	return *r.EditFollowRequest()
}

func (r *State) EditFollowRequest() *models.FollowRequestType {
	if &r.FollowRequest == nil {
		r.InitFollowRequest()
	}
	return &r.FollowRequest
}
