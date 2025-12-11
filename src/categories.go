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

func (c *Category) IsEmpty() bool {
	return c == nil || *c == Category{}
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
	id_int, err := strconv.Atoi(id)
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
