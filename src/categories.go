package forum

import (
	"net/http"
	"strconv"
)

type Category struct {
	Id   int
	Name string
}

type Categories []Category

func ReturnMockCategories() Categories {
	return Categories{
		{
			Id:   1,
			Name: "various",
		},
		{
			Id:   2,
			Name: "general",
		},
	}
}

func InitCategories() {
	for _, category := range ReturnMockCategories() {
		err := addCategory(category)
		if err != nil {
			var e Error
			e = e.Consume(err)
			e.LogError()
		}
	}
}

func (c *Category) IsEmpty() bool {
	empty := Category{}
	if c == nil || *c == empty {
		return true
	}
	return false
}

func (c *Category) validateCategory() error {
	if len((*c).Name) == 0 {
		return ErrorCategoryNameEmpty
	}
	if len((*c).Name) >= 128 {
		return ErrorCategoryNameTooLong
	}
	return nil
}

func (c *Category) DoesCategoryExist() bool {
	categories, err := getAllCategories()
	if err != nil {
		var e Error
		e = e.Consume(err)
		e.LogError()
		return false
	}
	for _, category := range categories {
		if category.Name == (*c).Name {
			return true
		}
	}
	return false
}

func showCategories(res http.ResponseWriter, _ *http.Request, user User) {
	categories, err := getAllCategories()
	if err != nil {
		var e Error
		e = e.Consume(err)
		e.LogError()
		e.RespondError(res)
		return
	}
	data := ResponseStruct{
		WebsiteName: "Forum",
		User:        user,
		Categories:  categories,
	}
	respondView(res, "categories_view", data)
}

func showCategory(res http.ResponseWriter, req *http.Request, user User) {
	data := ReturnMockResponse()
	data.User = user
	id := req.URL.Query().Get("id")
	if len(id) == 0 {
		var e Error
		e = e.Consume(ErrorCategoryEmptyId)
		e.LogError()
		e.RespondError(res)
		return
	}
	id_int, err := strconv.Atoi(id)
	if err != nil {
		var e Error
		e = e.Consume(err)
		e.LogError()
		e.RespondError(res)
		return
	}
	category, err := getCategoryById(id_int)
	data.Categories = Categories{}
	data.Categories = append(data.Categories, category)
	posts, err := getPostsByCategoryId(id_int)
	data.Posts = posts
	respondView(res, "category_view", data)
}
