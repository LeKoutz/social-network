package state

import "forum/src/models"

func (r *State) InitGroup() {
	var group models.GroupType
	r.Groups = models.GroupsType{group}
}

func (r *State) EditGroup() *models.GroupType {
	if r.Groups == nil {
		r.InitGroup()
	}
	return &r.Groups[0]
}

func (r *State) GetGroup() models.GroupType {
	return *r.EditGroup()
}

func (r *State) SetGroup(group models.GroupType) *State {
	r.Groups = models.GroupsType{group}
	return r
}