package controllers

import (
	"forum/src/models"
	"net/http"
	"regexp"
	"strconv"
)

func showCategories(res http.ResponseWriter, _ *http.Request, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetUser(user).SetResponse(res)
	categories, err := models.GetAllCategories()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	data.SetCategories(categories)
	data.SetView("categories_view").WriteResponse()
}

func showCategory(res http.ResponseWriter, req *http.Request, user models.User) {
	data := models.ResponseStruct{}
	data.Init().SetUser(user).SetResponse(res)
	id := req.URL.Query().Get("id")
	if len(id) == 0 {
		(&models.Error{}).Consume(models.ErrorCategoryEmptyId).LogAndRespondError(res, user)
		return
	}
	ok, err := regexp.MatchString(`^\d+$`, id)
	if !ok {
		(&models.Error{}).Consume(models.ErrorInvalidCategoryId).LogAndRespondError(res, user)
		return
	}
	id_int, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	category, err := models.GetCategoryById(id_int)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	data.Categories = models.Categories{}
	data.Categories = append(data.Categories, category)
	posts, err := models.GetPostsByCategoryId(id_int)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(res, user)
		return
	}
	for i := range posts {
		err = posts[i].GetReactions()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
		err = posts[i].GetReactionsByUserId(user.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(res, user)
			return
		}
	}
	data.Posts = posts
	data.SetPosts(posts).SetView("category_view").WriteResponse()
}
