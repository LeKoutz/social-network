package controllers

import (
	"forum/src/models"
	"forum/src/views"
	"regexp"
	"strconv"
	"strings"
)

func showCategories(data models.ResponseStruct) {
	categories, err := models.GetAllCategories()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.SetCategories(categories)
	views.Categories(data)
}

func showCategory(data models.ResponseStruct) {
	id, ok := strings.CutPrefix(data.Request.RequestURI, "/category/view/")
	if !ok || len(id) == 0 {
		(&models.Error{}).Consume(models.ErrorCategoryEmptyId).LogAndRespondError(data.Response, data.User)
		return
	}
	ok, err := regexp.MatchString(`^\d+$`, id)
	if !ok {
		(&models.Error{}).Consume(models.ErrorInvalidCategoryId).LogAndRespondError(data.Response, data.User)
		return
	}
	id_int, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	category, err := models.GetCategoryById(id_int)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.Categories = models.Categories{}
	data.Categories = append(data.Categories, category)
	posts, err := models.GetPostsByCategoryId(id_int)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	for i := range posts {
		err = posts[i].GetReactions()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		err = posts[i].GetReactionsByUserId(data.User.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	data.Posts = posts
	data.SetPosts(posts)
	views.Category(data)
}
