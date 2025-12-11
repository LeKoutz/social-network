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
			(&Error{}).Consume(err).LogError()
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

func (c *Categories) IsEmpty() bool {
	if c == nil {
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
	categories, err := getAllCategories()
	if err != nil {
		(&Error{}).Consume(err).LogError()
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
		(&Error{}).Consume(ErrorCategoryEmptyId).LogAndRespondError(res, user)
		return
	}
	id_int, err := strconv.Atoi(id)
	if err != nil {
		(&Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	category, err := getCategoryById(id_int)
	data.Categories = Categories{}
	data.Categories = append(data.Categories, category)
	posts, err := getPostsByCategoryId(id_int)
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
	respondView(res, "category_view", data)
}

func (p *Post) getReactions() error {
	var err error
	(*p).Likes, err = getLikesCountByPostId((*p).Id)
	if err != nil {
		return err
	}
	(*p).Dislikes, err = getDislikesCountByPostId((*p).Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) getReactionsByUserId(user_id int) error {
	var err error
	(*p).Liked, err = hasUserAlreadyLikedPost(user_id, (*p).Id)
	if err != nil {
		return err
	}
	(*p).Disliked, err = hasUserAlreadyDislikedPost(user_id, (*p).Id)
	if err != nil {
		return err
	}
	return nil
}
