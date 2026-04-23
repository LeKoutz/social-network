package controllers

import (
	"fmt"
	"forum/src/models"
	"forum/src/views"
	"forum/src/utils"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func showPost(data models.ResponseStruct) {
	id, ok := strings.CutPrefix(data.Request.RequestURI, "/post/view/")
	if !ok || len(id) == 0 {
		(&models.Error{}).Consume(models.ErrorPostEmptyId).LogAndRespondError(data.Response, data.User)
		return
	}
	ok, err := regexp.MatchString(`^\d+$`, id)
	if !ok {
		(&models.Error{}).Consume(models.ErrorInvalidPostId).LogAndRespondError(data.Response, data.User)
		return
	}
	id_int, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	var post models.Post
	post.Id = id_int
	err = post.GetById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	comments, err := post.GetComments()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	for i := range comments {
		comments[i].GetReactions()
		comments[i].GetReactionsByUserId(data.User.Id)
	}
	post.Comments = comments
	categories, err := models.GetCategoriesByPostId(post.Id)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	post.Categories = categories
	err = post.GetReactions()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	err = post.GetReactionsByUserId(data.User.Id)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	// Update notifications to read for this post
	for i, notification := range data.User.Notifications {
		if notification.TargetId == post.Id && !notification.Read {
			err = data.User.MarkNotificationAsRead(notification.Id)
			if err != nil {
				(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
				return
			}
			data.User.Notifications[i].Read = true
		}
	}
	data.Posts = models.Posts{post}
	views.PostView(data)
}

func showPosts(data models.ResponseStruct) {
	posts, err := models.GetAllPosts()
	if err != nil {
		data.Error.Consume(err)
		(&models.Error{}).Consume(err).RespondError(data.Response, data.User)
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
	data.SetPosts(posts)
	views.PostsView(data)
}

func createPost(data models.ResponseStruct) {
	if data.Request.Method == http.MethodGet {
		categories, err := models.GetAllCategories()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		data.SetCategories(categories)
		views.PostCreate(data)
		return
	}
	if data.Request.Method != http.MethodPost {
		data.SetErrorConsume(models.ErrorMethodNotAllowed)
		return
	}
	if !data.User.LoggedIn {
		data.Error.Consume(models.ErrorPostPermissionDenied)
		views.UserLogin(data)
		return
	}
	post, err := parseCreatePostRequest(data)
	if err != nil {
		data.Error.Consume(err)
		views.PostCreate(data)
		return
	}
	postId, err := post.Add()
	if err != nil {
		data.Error.Consume(err)
		views.PostCreate(data)
		return
	}
	redirectURL := fmt.Sprintf("/post/view/%d", postId)
	http.Redirect(data.Response, data.Request, redirectURL, http.StatusSeeOther)
}

func parseCreatePostRequest(data models.ResponseStruct) (models.Post, error) {
	title := data.Request.FormValue("title")
	body := data.Request.FormValue("body")
	categories, err := models.GetAllCategories()
	if err != nil {
		return models.Post{}, err
	}
	var PostCategories models.Categories
	for _, category := range categories {
		cc := fmt.Sprintf("category-%d", category.Id)
		if data.Request.Form.Has(cc) && data.Request.Form.Get(cc) == "on" {
			PostCategories = append(PostCategories, category)
		}
	}
	imagePath := ""
	imageFile, _, err := data.Request.FormFile("image")
	if err == nil {
		defer imageFile.Close()
		imagePath, err = models.SaveImage(imageFile)
		if err != nil {
			return models.Post{}, err
		}
	}
	return models.Post{
		Title:      title,
		Body:       body,
		UserId:     data.User.Id,
		Categories: PostCategories,
		ImagePath: imagePath,
	}, nil
}

func handlePostCreate(data models.ResponseStruct) {
	if data.Request.Method != http.MethodPost && data.Request.Method != http.MethodGet {
		data.SetErrorConsume(models.ErrorMethodNotAllowed)
	}
	if !data.User.LoggedIn {
		(&models.Error{}).Consume(models.ErrorPostPermissionDenied).LogAndRespondError(data.Response, data.User)
		return
	}
	switch {
	case strings.Compare(data.Request.RequestURI, "/post/create") == 0:
		switch data.Request.Method {
		case http.MethodGet:
			categories, err := models.GetAllCategories()
			if err != nil {
				(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
				return
			}
			data.SetCategories(categories)
			views.PostCreate(data)
			return
		case http.MethodPost:
			err := data.Request.ParseMultipartForm(models.MaxImageSize)
			if err != nil {
				(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
				return
			}
				createPost(data)
		default:
			(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		}
	}
}

func handlePostReaction(data models.ResponseStruct) {
	if !data.User.LoggedIn {
		(&models.Error{}).Consume(models.ErrorPostPermissionDenied).LogAndRespondError(data.Response, data.User)
		return
	}
	postIdStr := data.Request.FormValue("post-id")
	if len(postIdStr) == 0 {
		(&models.Error{}).Consume(models.ErrorPostEmptyId).LogAndRespondError(data.Response, data.User)
		return
	}
	ok, err := regexp.MatchString(`^\d+$`, postIdStr)
	if !ok {
		(&models.Error{}).Consume(models.ErrorInvalidPostId).LogAndRespondError(data.Response, data.User)
		return
	}
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	post := models.Post{Id: postId}
	err = post.GetById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	notification := models.Notification{
		UserId:  post.User.Id,
		ActorId: data.User.Id,
		TargetId: post.Id,
	}
	if data.Request.FormValue("action") == "like" {
		err = DoLikePost(data.User.Id, postId)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		if data.User.Id != post.User.Id {
			notification.Type = "like"
			notification.Timestamp = utils.GetCurrentTimestamp()
			err = models.CreateNotification(&notification)
		}
	}
	if data.Request.FormValue("action") == "dislike" {
		err = DoDislikePost(data.User.Id, postId)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		if data.User.Id != post.User.Id {
			notification.Type = "dislike"
			notification.Timestamp = utils.GetCurrentTimestamp()
			err = models.CreateNotification(&notification)
		}
	}
	http.Redirect(data.Response, data.Request, "/post/view/"+postIdStr, http.StatusSeeOther)
}
