package state

import "forum/src/models"

func (r *State) EditCategories() *models.CategoriesType {
	return &r.Categories
}

func (r *State) GetCategories() models.CategoriesType {
	return r.Categories
}

func (r *State) SetCategories(categories models.CategoriesType) *State {
	r.Categories = categories
	return r
}
