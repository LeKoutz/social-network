package controllers

import (
	"errors"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
)

func ShowChatHistory(data state.StateController) error {
	var err error
	user1 := data.GetUser().Id
	user2 := data.GetChatMessage(0).RecipientId
	if user1 == user2 {
		err = ferror.ErrorNotFound
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	var recipient models.UserType
	recipient.Id = user2
	if err = recipient.GetById(); err != nil {
		err = ferror.ErrorNotFound
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	offset := data.GetChatOffset()
	if err = data.EditChatMessages().GetChatHistory(user1, user2, offset); err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	for i := range data.GetChatMessages() {
		if data.GetChatMessage(int64(i)).RecipientId == data.GetUser().Id {
			if err = data.EditChatMessage(int64(i)).MarkAsRead(); err != nil {
				if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
				return err
			}
		}
	}
	return nil
}

func ServeUnreadMessages(data state.StateController) error {
	return data.EditChatMessages().GetUnreadMessageIds(data.GetUser().Id)
}
