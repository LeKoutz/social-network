package controllers

import (
	"errors"
	"forum/src/ferror"
	"forum/src/state"
	"forum/src/utils"
)

func CommentCreate(data state.StateController) error {
	var err error
	data.EditComment().UserId = data.GetUser().Id
	data.EditComment().Body = data.GetRequest().FormValue("comment")
	data.EditComment().PostId, err = ParsePostId(data)
	if err != nil {
		return ferror.ErrorInvalidPostId
	}
	data.EditPost().Id = data.GetComment().PostId
	err = data.EditPost().SelectPostById()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	err = data.EditComment().Add()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	data.SetRedirect(GetRedirectLinkToCommentOfPost(data))
	err = data.EditComment().CreateCommentNotification(data.GetPost())
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	// TODO Look at this
	// data.SetResponse().Header().Set("Content-Type", "application/json")
	// json.NewEncoder(data.Response).Encode(map[string]int64{
	// 	"postId":    post.Id,
	// 	"commentId": comment.Id,
	// })
	return nil
}

func CommentReaction(data state.StateController) error {
	var err error
	data.EditComment().Id, err = ParseCommentId(data)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	data.EditComment().PostId, err = ParsePostId(data)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	data.EditPost().Id = data.EditComment().PostId
	err = data.EditComment().SelectCommentById()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	switch data.GetRequest().FormValue("action") {
	case "like":
		err = data.EditUser().LikeComment(data.EditComment().Id)
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
		err = data.EditComment().CreateReactionNotification(data.GetUser().Id, "commentLike")
		if err != nil {
			(&ferror.Error{}).Consume(err).LogError()
		}
	case "dislike":
		err = data.EditUser().DislikeComment(data.EditComment().Id)
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
		err = data.EditComment().CreateReactionNotification(data.GetUser().Id, "commentDislike")
		if err != nil {
			(&ferror.Error{}).Consume(err).LogError()
		}
	default:
		return ferror.ErrorUnknownAction
	}
	data.SetRedirect(GetRedirectLinkToCommentOfPost(data))
	return nil
}

func CommentDelete(data state.StateController) error {
	var err error
	data.EditComment().Id, err = ParseCommentId(data)
	if err != nil {
		return ferror.ErrorInvalidCommentId
	}
	data.EditPost().Id, err = ParsePostId(data)
	if err != nil {
		return ferror.ErrorInvalidCommentId
	}
	err = data.EditComment().SelectCommentById()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	if data.EditComment().UserId != data.GetUser().Id {
		return ferror.ErrorCommentPermissionDenied
	}
	err = data.EditComment().DeleteCommentById()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	data.SetRedirect(GetRedirectLinkToPost(data))
	return nil
}

func commonCommentEditPrework(data state.StateController) error {
	var err error
	err = ValidateFormCommentEdit(data)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	return VerifyCommentOwnership(data)
}

func CommentEdit(data state.StateController) error {
	var err error
	err = commonCommentEditPrework(data)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	if data.GetRequest().FormValue("save-comment") == "1" {
		err = UpdateCommentFromForm(data)
		if err != nil {
			return errors.Join(utils.GetFunctionName(), err)
		}
		return nil
	}
	return ShowEditComment(data)
}

func ValidateFormCommentEdit(data state.StateController) error {
	var err error
	data.EditComment().Id, err = ParseCommentId(data)
	if err != nil {
		return ferror.ErrorInvalidCommentId
	}
	data.EditPost().Id, err = ParsePostId(data)
	if err != nil {
		return ferror.ErrorInvalidPostId
	}
	data.EditComment().PostId = data.GetPost().Id
	return nil
}

func VerifyCommentOwnership(data state.StateController) error {
	var err error
	err = data.EditComment().SelectCommentById()
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	// Check your priviledge
	if data.GetComment().UserId != data.GetUser().Id {
		return ferror.ErrorCommentPermissionDenied
	}
	return nil
}

func ShowEditComment(data state.StateController) error {
	var err error
	err = getPostDataById(data)
	if err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}
	data.SetEditCommentId(data.GetComment().Id)
	return nil
}

func UpdateCommentFromForm(data state.StateController) error {
	data.EditComment().Body = data.GetRequest().FormValue("comment")
	return data.EditComment().Update()
}
