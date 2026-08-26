package state

import "forum/src/models"

func (r *State) InitCategory() {
	var category models.CategoryType
	r.Categories = models.CategoriesType{category}
}

func (r *State) EditCategory() *models.CategoryType {
	if r.Categories == nil {
		r.InitCategory()
	}
	return &r.Categories[0]
}

func (r *State) GetCategory() models.CategoryType {
	return *r.EditCategory()
}

func (r *State) SetCategory(category models.CategoryType) *State {
	r.Categories = models.CategoriesType{category}
	return r
}
