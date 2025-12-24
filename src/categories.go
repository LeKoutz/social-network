package forum

import (
	"net/http"
	"regexp"
	"strconv"
)

type Category struct {
	Id   int64
	Name string
	Description string
}

type Categories []Category

func (c *Category) IsEmpty() bool {
	return c == nil || *c == Category{}
}

func (c *Categories) IsEmpty() bool {
	if c == nil {
		return true
	}
	if c == (&Categories{}) {
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
		(&Error{}).Consume(err).LogError()
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
	data := ResponseStruct{}
	data.Init().SetUser(user)
	categories, err := getAllCategories()
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	data.SetCategories(categories)
	data.SetView("categories_view").WriteResponse(res)
}

func showCategory(res http.ResponseWriter, req *http.Request, user User) {
	data := ResponseStruct{}
	data.Init().SetUser(user)
	id := req.URL.Query().Get("id")
	if len(id) == 0 {
		(&Error{}).Consume(ErrorCategoryEmptyId).LogAndRespondError(res, user)
		return
	}
	ok, err := regexp.MatchString(`^\d+$`, id)
	if !ok {
		(&Error{}).Consume(ErrorInvalidCategoryId).LogAndRespondError(res, user)
		return
	}
	id_int, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	category, err := getCategoryById(id_int)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	data.Categories = Categories{}
	data.Categories = append(data.Categories, category)
	posts, err := getPostsByCategoryId(id_int)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	for i := range posts {
		err = posts[i].getReactions()
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
		err = posts[i].getReactionsByUserId(user.Id)
		if err != nil {
			(&Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	data.Posts = posts
	data.SetPosts(posts).SetView("category_view").WriteResponse(res)
}
