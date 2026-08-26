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
			err = ferror.ErrorBadRequest
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
	}
	if formData["password1"] != formData["password2"] {
		err = ferror.ErrorPasswordMismatch
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.EditUser().Username = formData["username"]
	data.EditUser().Email = formData["email"]
	data.EditUser().Password = formData["password1"]
	data.EditUser().FirstName = formData["first_name"]
	data.EditUser().LastName = formData["last_name"]
	age, err := utils.StringToInt64(formData["age"])
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.EditUser().Age = age
	data.EditUser().Gender = formData["gender"]
	return nil
}

func ParseLoginForm(data state.StateHandler) error {
	var err error
	identifier := data.GetRequest().FormValue("identifier")
	if len(identifier) == 0 {
		err = ferror.ErrorEmailFieldEmpty
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	password := data.GetRequest().FormValue("password")
	if len(password) == 0 {
		err = ferror.ErrorPasswordFieldEmpty
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
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
		err = ferror.ErrorPostEmptyId
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, err
	}
Convert:
	postId, err = utils.StringToInt64(postIdStr)
	if err != nil || postId == 0 {
		err = ferror.ErrorInvalidPostId
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, err
	}
	return postId, nil
}

func ParseCommentId(data state.StateHandler) (int64, error) {
	var err error
	commentIdStr := data.GetRequest().FormValue("comment-id")
	if len(commentIdStr) == 0 {
		err = ferror.ErrorCommentEmptyId
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, err
	}
	commentId, err := utils.StringToInt64(commentIdStr)
	if err != nil {
		err = ferror.ErrorInvalidCommentId
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, err
	}
	return commentId, nil
}

func ParseCategoryId(data state.StateHandler) (int64, error) {
	var err error
	id, ok := strings.CutPrefix(data.GetRequest().RequestURI, "/api/category/view/")
	if !ok || len(id) == 0 {
		err = ferror.ErrorCategoryEmptyId
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, err
	}
	categoryId, err := utils.StringToInt64(id)
	if err != nil {
		err = ferror.ErrorInvalidCategoryId
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, err
	}
	return categoryId, nil
}

func ParseCreatePostRequest(data state.StateHandler) error {
	var err error
	var categories models.CategoriesType
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.EditPost().UserId = data.EditUser().Id
	title := data.GetRequest().FormValue("title")
	if len(title) == 0 {
		err = ferror.ErrorPostTitleEmpty
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	body := data.GetRequest().FormValue("body")
	if len(body) == 0 {
		err = ferror.ErrorPostBodyEmpty
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.EditPost().Title = title
	data.EditPost().Body = body
	err = categories.GetAll()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	var post_cat models.CategoriesType
	for _, category := range categories {
		cc := fmt.Sprintf("category-%d", category.Id)
		if data.GetRequest().FormValue(cc) == "on" {
			post_cat = append(post_cat, category)
		}
	}
	if len(post_cat) == 0 {
		err = ferror.ErrorPostHasNoCategory
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.EditPost().Categories = post_cat
	imageFile, _, err := data.GetRequest().FormFile("image")
	if err == nil {
		defer imageFile.Close()
		data.EditPost().ImagePath, err = models.SaveImage(imageFile)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
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
		err = ferror.ErrorInvalidPostId
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.EditComment().PostId = data.GetPost().Id
	return nil
}

func ParseChatId(data state.StateHandler) (id, offset int64, err error) {
	uri, ok := strings.CutPrefix(data.GetRequest().RequestURI, "/api/chat/")
	if !ok {
		err = ferror.ErrorBadRequest
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, 0, err
	}
	idStr, offsetStr, found := strings.Cut(uri, "?offset=")
	if !found {
		err = ferror.ErrorBadRequest
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, 0, err
	}
	id, err = utils.StringToInt64(idStr)
	if err != nil {
		err = ferror.ErrorInvalidChatId
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, 0, err
	}
	offset, err = utils.StringToInt64(offsetStr)
	if err != nil {
		err = ferror.ErrorBadRequest
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return 0, 0, err
	}
	return id, offset, err
}
