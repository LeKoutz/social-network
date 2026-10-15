package parsers

import (
	"errors"
	"fmt"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
	"net/http"
	"strings"
)

func ParseRegistrationForm(data state.StateHandler) error {
	var err error
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	formData := map[string]string{
		"username":      data.GetRequest().FormValue("username"),
		"email":         data.GetRequest().FormValue("email"),
		"first_name":    data.GetRequest().FormValue("first_name"),
		"last_name":     data.GetRequest().FormValue("last_name"),
		"password1":     data.GetRequest().FormValue("password1"),
		"password2":     data.GetRequest().FormValue("password2"),
		"date_of_birth": data.GetRequest().FormValue("date_of_birth"),
	}
	for _, value := range formData {
		if len(value) == 0 {
			err = ferror.ErrorBadRequest
			return ferror.ReturnErr(err)
		}
	}
	if formData["password1"] != formData["password2"] {
		err = ferror.ErrorPasswordMismatch
		return ferror.ReturnErr(err)
	}
	data.EditUser().Username = formData["username"]
	data.EditUser().Email = formData["email"]
	data.EditUser().Password = formData["password1"]
	data.EditUser().FirstName = formData["first_name"]
	data.EditUser().LastName = formData["last_name"]
	data.EditUser().DateOfBirth = formData["date_of_birth"]
	data.EditUser().Nickname = data.GetRequest().FormValue("nickname")
	data.EditUser().About = data.GetRequest().FormValue("about_me")
	data.EditUser().AvatarURL = data.GetRequest().FormValue("avatar")
	data.EditUser().PrivateProfile = data.GetRequest().FormValue("private_profile") == "on"
	return nil
}

func ParseLoginForm(data state.StateHandler) error {
	var err error
	identifier := data.GetRequest().FormValue("identifier")
	if len(identifier) == 0 {
		err = ferror.ErrorEmailFieldEmpty
		return ferror.ReturnErr(err)
	}
	password := data.GetRequest().FormValue("password")
	if len(password) == 0 {
		err = ferror.ErrorPasswordFieldEmpty
		return ferror.ReturnErr(err)
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
		return 0, ferror.ReturnErr(err)
	}
Convert:
	postId, err = utils.StringToInt64(postIdStr)
	if err != nil || postId == 0 {
		err = ferror.ErrorInvalidPostId
		return 0, ferror.ReturnErr(err)
	}
	return postId, nil
}

func ParseCommentId(data state.StateHandler) (int64, error) {
	var err error
	commentIdStr := data.GetRequest().FormValue("comment-id")
	if len(commentIdStr) == 0 {
		err = ferror.ErrorCommentEmptyId
		return 0, ferror.ReturnErr(err)
	}
	commentId, err := utils.StringToInt64(commentIdStr)
	if err != nil {
		err = ferror.ErrorInvalidCommentId
		return 0, ferror.ReturnErr(err)
	}
	return commentId, nil
}

func ParseCategoryId(data state.StateHandler) (int64, error) {
	var err error
	id, ok := strings.CutPrefix(data.GetRequest().RequestURI, "/api/category/view/")
	if !ok || len(id) == 0 {
		err = ferror.ErrorCategoryEmptyId
		return 0, ferror.ReturnErr(err)
	}
	categoryId, err := utils.StringToInt64(id)
	if err != nil {
		err = ferror.ErrorInvalidCategoryId
		return 0, ferror.ReturnErr(err)
	}
	return categoryId, nil
}

func ParseGroupId(data state.StateHandler) (int64, error) {
	var err error
	id, ok := strings.CutPrefix(data.GetRequest().RequestURI, "/api/group/view/")
	if !ok || len(id) == 0 {
		err = ferror.ErrorGroupEmptyId
		return 0, ferror.ReturnErr(err)
	}
	groupId, err := utils.StringToInt64(id)
	if err != nil {
		err = ferror.ErrorInvalidGroupId
		return 0, ferror.ReturnErr(err)
	}
	return groupId, nil
}

func ParseCreatePostRequest(data state.StateHandler) error {
	var err error
	var categories models.CategoriesType
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.EditPost().UserId = data.EditUser().Id
	title := data.GetRequest().FormValue("title")
	if len(title) == 0 {
		err = ferror.ErrorPostTitleEmpty
		return ferror.ReturnErr(err)
	}
	body := data.GetRequest().FormValue("body")
	if len(body) == 0 {
		err = ferror.ErrorPostBodyEmpty
		return ferror.ReturnErr(err)
	}
	data.EditPost().Title = title
	data.EditPost().Body = body
	groupIdStr := data.GetRequest().FormValue("group-id")
	if len(groupIdStr) != 0 {
		groupId, err := utils.StringToInt64(groupIdStr)
		if err != nil {
			err = ferror.ErrorInvalidGroupId
			return ferror.ReturnErr(err)
		}
		data.EditPost().GroupId = groupId
	}
	err = categories.GetAll()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	var post_cat models.CategoriesType
	for _, category := range categories {
		cc := fmt.Sprintf("category-%d", category.Id)
		if data.GetRequest().FormValue(cc) == "on" {
			post_cat = append(post_cat, category)
		}
	}
	if len(post_cat) == 0 && data.GetPost().GroupId == 0 {
		err = ferror.ErrorPostHasNoCategory
		return ferror.ReturnErr(err)
	}
	data.EditPost().Categories = post_cat
	imageFile, _, err := data.GetRequest().FormFile("image")
	if err == nil {
		defer imageFile.Close()
		data.EditPost().ImagePath, err = models.SaveImage(imageFile)
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	return nil
}

func ParseCreateCommentRequest(data state.StateHandler) error {
	var err error
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.EditComment().UserId = data.GetUser().Id
	data.EditComment().Body = data.GetRequest().FormValue("comment")
	data.EditPost().Id, err = ParsePostId(data)
	if err != nil {
		err = ferror.ErrorInvalidPostId
		return ferror.ReturnErr(err)
	}
	data.EditComment().PostId = data.GetPost().Id
	imageFile, _, err := data.GetRequest().FormFile("image")
	if err == nil {
		defer imageFile.Close()
		data.EditComment().ImagePath, err = models.SaveImage(imageFile)
		if err != nil {
			return ferror.ReturnErr(err)
		}
	}
	return nil
}

func ParseChatId(data state.StateHandler) (id, offset int64, err error) {
	uri, ok := strings.CutPrefix(data.GetRequest().RequestURI, "/api/chat/")
	if !ok {
		err = ferror.ErrorBadRequest
		return 0, 0, ferror.ReturnErr(err)
	}
	idStr, offsetStr, found := strings.Cut(uri, "?offset=")
	if !found {
		err = ferror.ErrorBadRequest
		return 0, 0, ferror.ReturnErr(err)
	}
	id, err = utils.StringToInt64(idStr)
	if err != nil {
		err = ferror.ErrorInvalidChatId
		return 0, 0, ferror.ReturnErr(err)
	}
	offset, err = utils.StringToInt64(offsetStr)
	if err != nil {
		err = ferror.ErrorBadRequest
		return 0, 0, ferror.ReturnErr(err)
	}
	return id, offset, nil
}

func ParseUserFollowRequest(data state.StateHandler) error {
	var err error
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	userIdStr := data.GetRequest().FormValue("user-id")
	if len(userIdStr) == 0 {
		err = ferror.ErrorBadRequest
		return ferror.ReturnErr(err)
	}
	userId, err := utils.StringToInt64(userIdStr)
	if err != nil {
		err = ferror.ErrorInvalidUserId
		return ferror.ReturnErr(err)
	}
	data.EditFollowRequest().ToUserId = userId
	return nil
}

func ParseAcceptFollowRequest(data state.StateHandler) error {
	return parseFollowRequestId(data, "/api/follow/accept/", ferror.ErrorInvalidInvitationId)
}

func ParseDenyFollowRequest(data state.StateHandler) error {
	return parseFollowRequestId(data, "/api/follow/deny/", ferror.ErrorInvalidInvitationId)
}

func ParseUnfollowUserRequest(data state.StateHandler) error {
	var err error
	var userId int64
	id, ok := strings.CutPrefix(data.GetRequest().RequestURI, "/api/unfollow/")
	if !ok {
		err = ferror.ErrorBadRequest
		return ferror.ReturnErr(err)
	}
	userId, err = utils.StringToInt64(id)
	if err != nil {
		err = ferror.ErrorInvalidUserId
		return ferror.ReturnErr(err)
	}
	data.EditFollowRequest().ToUserId = userId
	return nil
}

func parseFollowRequestId(data state.StateHandler, prefix string, invalidErr error) error {
	var err error
	var invitationId int64
	id, ok := strings.CutPrefix(data.GetRequest().RequestURI, prefix)
	if !ok {
		err = ferror.ErrorBadRequest
		return ferror.ReturnErr(err)
	}
	invitationId, err = utils.StringToInt64(id)
	if err != nil {
		err = invalidErr
		return ferror.ReturnErr(err)
	}
	data.EditFollowRequest().Id = invitationId
	return nil
}

func ParseCreateGroupRequest(data state.StateHandler) error {
	var err error
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.EditGroup().OwnerUserId = data.GetUser().Id
	data.EditGroup().Title = data.GetRequest().FormValue("title")
	data.EditGroup().Description = data.GetRequest().FormValue("description")
	return nil
}

func ParseGroupInvitationRequest(data state.StateHandler) error {
	var err error
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	id_str := data.GetRequest().FormValue("id")
	data.EditGroupInvitation().Id, err = utils.StringToInt64(id_str)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}

func ParseInviteToGroupRequest(data state.StateHandler) error {
	var err error
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.EditGroupInvitation().Timestamp = utils.GetCurrentTimestamp()
	data.EditGroupInvitation().Status = "pending"
	group_id_str := data.GetRequest().FormValue("group-id")
	data.EditGroupInvitation().GroupId, err = utils.StringToInt64(group_id_str)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	to_user_id_str := data.GetRequest().FormValue("to-user-id")
	data.EditGroupInvitation().ToUserId, err = utils.StringToInt64(to_user_id_str)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.EditGroupInvitation().InvitedBy = data.GetUser().Id
	return nil
}

func ParseJoinGroupRequest(data state.StateHandler) error {
	var err error
	err = data.GetRequest().ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.EditGroupInvitation().Timestamp = utils.GetCurrentTimestamp()
	data.EditGroupInvitation().Status = "pending"
	group_id_str := data.GetRequest().FormValue("group-id")
	data.EditGroupInvitation().GroupId, err = utils.StringToInt64(group_id_str)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	data.EditGroupInvitation().ToUserId = data.GetUser().Id
	return nil
}

func ParseProfileId(data state.StateHandler) (int64, error) {
	var err error
	id, ok := strings.CutPrefix(data.GetRequest().RequestURI, "/api/profile/view/")
	if !ok || len(id) == 0 {
		err = ferror.ErrorProfileEmptyId
		return 0, ferror.ReturnErr(err)
	}
	profileId, err := utils.StringToInt64(id)
	if err != nil {
		err = ferror.ErrorInvalidProfileId
		return 0, ferror.ReturnErr(err)
	}
	return profileId, nil
}

func ParseGroupInvitations(data state.StateHandler) error {
	var str string
	var ok bool
	var i64 int64
	var err error
	str, ok = strings.CutPrefix(data.GetRequest().RequestURI, "/api/group/joins/")
	if ok && len(str) != 0 {
		goto Convert
	}
	str, ok = strings.CutPrefix(data.GetRequest().RequestURI, "/api/group/invites/")
	if !ok || len(str) == 0 {
		data.EditGroup().Id = 0
		return nil
	}
Convert:
	i64, err = utils.StringToInt64(str)
	if err != nil || i64 == 0 {
		err = ferror.ErrorInvalidGroupId
		return ferror.ReturnErr(err)
	}
	data.EditGroup().Id = i64
	return nil
}

func ParseGroupChatRequest(data state.StateHandler) (groupId, offset int64, err error) {
	id, ok := strings.CutPrefix(data.GetRequest().URL.Path, "/api/group/chat/")
	if !ok || len(id) == 0 {
		err = ferror.ErrorGroupEmptyId
		return 0, 0, ferror.ReturnErr(err)
	}
	groupId, err = utils.StringToInt64(id)
	if err != nil {
		err = ferror.ErrorInvalidGroupId
		return 0, 0, ferror.ReturnErr(err)
	}
	offsetStr := data.GetRequest().URL.Query().Get("offset")
	if offsetStr == "" {
		return groupId, 0, nil
	}
	offset, err = utils.StringToInt64(offsetStr)
	if err != nil {
		err = ferror.ErrorBadRequest
		return 0, 0, ferror.ReturnErr(err)
	}
	return groupId, offset, nil
}

func ParseProfileEditRequest(data state.StateHandler) error {
	req := data.GetRequest()

	if err := req.ParseMultipartForm(models.MaxImageSize); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			return ferror.ReturnErr(ferror.ErrorImageTooBig)
		}
		return ferror.ReturnErr(ferror.ErrorBadRequest)
	}

	profile := data.EditProfile()
	profile.FirstName = strings.TrimSpace(req.PostFormValue("first_name"))
	profile.LastName = strings.TrimSpace(req.PostFormValue("last_name"))
	profile.Nickname = strings.TrimSpace(req.PostFormValue("nickname"))
	profile.About = strings.TrimSpace(req.PostFormValue("about"))
	profile.DateOfBirth = strings.TrimSpace(req.PostFormValue("date_of_birth"))
	profile.PrivateProfile = req.PostFormValue("private_profile") == "on"
	profile.RemoveAvatar = req.PostFormValue("remove_avatar") == "on"

	avatar, _, err := req.FormFile("avatar")
	if errors.Is(err, http.ErrMissingFile) {
		return nil
	}
	if err != nil {
		return ferror.ReturnErr(err)
	}

	profile.AvatarFile = avatar
	return nil
}
