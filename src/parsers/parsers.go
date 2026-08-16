package parsers

import (
	"errors"
	"fmt"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
	"strings"
)

func ParseRegistrationForm(data state.StateHandler) error {
	var err error
	formData := map[string]string{
		"username":   data.GetRequest().FormValue("username"),
		"email":      data.GetRequest().FormValue("email"),
		"first_name": data.GetRequest().FormValue("first_name"),
		"last_name":  data.GetRequest().FormValue("last_name"),
		"age":        data.GetRequest().FormValue("age"),
		"gender":     data.GetRequest().FormValue("gender"),
		"password1":  data.GetRequest().FormValue("password1"),
		"password2":  data.GetRequest().FormValue("password2"),
	}
	for _, value := range formData {
		if len(value) == 0 {
			return ferror.ErrorBadRequest
		}
	}
	if formData["password1"] != formData["password2"] {
		return ferror.ErrorPasswordMismatch
	}
	data.EditUser().Username = formData["username"]
	data.EditUser().Email = formData["email"]
	data.EditUser().Password = formData["password1"]
	data.EditUser().FirstName = formData["first_name"]
	data.EditUser().LastName = formData["last_name"]
	age, err := utils.StringToInt64(formData["age"])
	if err != nil {
		return err
	}
	data.EditUser().Age = age
	data.EditUser().Gender = formData["gender"]
	return nil
}

func ParseLoginForm(data state.StateHandler) error {
	identifier := data.GetRequest().FormValue("identifier")
	if len(identifier) == 0 {
		return ferror.ErrorEmailFieldEmpty
	}
	password := data.GetRequest().FormValue("password")
	if len(password) == 0 {
		return ferror.ErrorPasswordFieldEmpty
	}
	data.EditUser().Identifier = identifier
	data.EditUser().Password = password
	return nil
}

func ParsePostId(data state.StateHandler) (int64, error) {
	var postIdStr string
	var postId int64
	var ok bool
	var err error
	postIdStr = data.GetRequest().FormValue("post-id")
	if len(postIdStr) != 0 {
		utils.LogDebug(postIdStr)
		goto Convert
	}
	postIdStr, ok = strings.CutPrefix(data.GetRequest().RequestURI, "/api/post/view/")
	if ok && len(postIdStr) != 0 {
		utils.LogDebug(postIdStr)
		goto Convert
	}
	postIdStr, ok = strings.CutPrefix(data.GetRequest().RequestURI, "/api/post/edit/")
	if ok && len(postIdStr) != 0 {
		utils.LogDebug(postIdStr)
		goto Convert
	}
	if len(postIdStr) == 0 {
		utils.LogDebug(postIdStr)
		return 0, ferror.ErrorPostEmptyId
	}
Convert:
	postId, err = utils.StringToInt64(postIdStr)
	if err != nil || postId == 0 {
		utils.LogDebug(postIdStr)
		return 0, ferror.ErrorInvalidPostId
	}
	return postId, nil
}

func ParseCommentId(data state.StateHandler) (int64, error) {
	commentIdStr := data.GetRequest().FormValue("comment-id")
	if len(commentIdStr) == 0 {
		return 0, ferror.ErrorCommentEmptyId
	}
	commentId, err := utils.StringToInt64(commentIdStr)
	if err != nil {
		return 0, ferror.ErrorInvalidCommentId
	}
	return commentId, nil
}

func ParseCategoryId(data state.StateHandler) (int64, error) {
	id, ok := strings.CutPrefix(data.GetRequest().RequestURI, "/api/category/view/")
	if !ok || len(id) == 0 {
		return 0, ferror.ErrorCategoryEmptyId
	}
	categoryId, err := utils.StringToInt64(id)
	if err != nil {
		return 0, ferror.ErrorInvalidCategoryId
	}
	return categoryId, nil
}

func ParseCreatePostRequest(data state.StateHandler) error {
	var err error
	var categories models.CategoriesType
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	data.EditPost().UserId = data.EditUser().Id
	data.EditPost().Title = data.GetRequest().FormValue("title")
	data.EditPost().Body = data.GetRequest().FormValue("body")
	err = categories.GetAll()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	var post_cat models.CategoriesType
	for _, category := range categories {
		cc := fmt.Sprintf("category-%d", category.Id)
		if data.GetRequest().FormValue(cc) == "on" {
			post_cat = append(post_cat, category)
		}
	}
	data.EditPost().Categories = post_cat
	imageFile, _, err := data.GetRequest().FormFile("image")
	if err == nil {
		defer imageFile.Close()
		data.EditPost().ImagePath, err = models.SaveImage(imageFile)
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
	}
	return nil
}

func ParseCreateCommentRequest(data state.StateHandler) error {
	var err error
	data.EditComment().UserId = data.GetUser().Id
	data.EditComment().Body = data.GetRequest().FormValue("comment")
	data.EditPost().Id, err = ParsePostId(data)
	if err != nil {
		return ferror.ErrorInvalidPostId
	}
	data.EditComment().PostId = data.GetPost().Id
	return nil
}

func ParseChatId(data state.StateHandler) (id, offset int64, err error) {
	uri, ok := strings.CutPrefix(data.GetRequest().RequestURI, "/api/chat/")
	if !ok {
		return 0, 0, ferror.ErrorBadRequest
	}
	idStr, offsetStr, found := strings.Cut(uri, "?offset=")
	if !found {
		return 0, 0, ferror.ErrorBadRequest
	}
	id, err = utils.StringToInt64(idStr)
	if err != nil {
		return 0, 0, ferror.ErrorInvalidChatId
	}
	offset, err = utils.StringToInt64(offsetStr)
	if err != nil {
		return 0, 0, ferror.ErrorBadRequest
	}
	return id, offset, err
}
