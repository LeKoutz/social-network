package controllers

import (
	"errors"
	"forum/src/ferror"
	"forum/src/state"
	"forum/src/utils"
)

func CommentCreate(data state.StateController) error {
	var err error
	err = data.EditComment().Add()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = data.EditPost().GetById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	err = data.EditComment().CreateCommentNotification(data.GetPost())
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func CommentReaction(data state.StateController) error {
	var err error
	err = data.EditComment().GetById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.EditPost().Id = data.GetComment().PostId
	switch data.GetRequest().FormValue("action") {
	case "like":
		err = data.EditUser().LikeComment(data.GetComment().Id)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		err = data.EditComment().CreateReactionNotification(data.GetUser().Id, "commentLike")
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			(&ferror.Error{}).Consume(err).LogError()
		}
	case "dislike":
		err = data.EditUser().DislikeComment(data.GetComment().Id)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		err = data.EditComment().CreateReactionNotification(data.GetUser().Id, "commentDislike")
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			(&ferror.Error{}).Consume(err).LogError()
		}
	default:
		err = ferror.ErrorUnknownAction
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func CommentDelete(data state.StateController) error {
	var err error
	err = VerifyCommentOwnership(data)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.EditPost().Id = data.GetComment().PostId
	err = data.EditComment().Delete()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func CommentEdit(data state.StateController) error {
	var err error
	err = VerifyCommentOwnership(data)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.EditPost().Id = data.GetComment().PostId
	if data.GetRequest().FormValue("save-comment") == "1" {
		err = UpdateCommentFromForm(data)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return err
		}
		return nil
	}
	return ShowEditComment(data)
}

func VerifyCommentOwnership(data state.StateController) error {
	var err error
	err = data.EditComment().GetById()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	// Check your priviledge
	if data.GetComment().UserId != data.GetUser().Id {
		err = ferror.ErrorCommentPermissionDenied
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func ShowEditComment(data state.StateController) error {
	var err error
	data.EditPost().Id = data.GetComment().PostId
	err = getPostDataById(data)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	data.SetEditCommentId(data.GetComment().Id)
	return nil
}

func UpdateCommentFromForm(data state.StateController) error {
	data.EditComment().Body = data.GetRequest().FormValue("comment")
	return data.EditComment().Update()
}
