package state

import "forum/src/models"

func (r *State) InitProfile() {
	var profile models.UserProfileType
	r.Profiles = models.UserProfilesType{profile}
}

func (r *State) EditProfile() *models.UserProfileType {
	if r.Profiles == nil {
		r.InitProfile()
	}
	return &r.Profiles[0]
}

func (r *State) GetProfile() models.UserProfileType {
	return *r.EditProfile()
}

func (r *State) SetProfile(profile models.UserProfileType) *State {
	r.Profiles = models.UserProfilesType{profile}
	return r
}