package controllers

import (
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
)

func ShowChatHistory(data state.StateController) error {
	var err error
	user1 := data.GetUser().Id
	user2 := data.GetChatMessage(0).RecipientId
	if user1 == user2 {
		err = ferror.ErrorNotFound
		return ferror.ReturnErr(err)
	}
	var recipient models.UserType
	recipient.Id = user2
	if err = recipient.GetById(); err != nil {
		err = ferror.ErrorNotFound
		return ferror.ReturnErr(err)
	}
	offset := data.GetChatOffset()
	if err = data.EditChatMessages().GetChatHistory(user1, user2, offset); err != nil {
		return ferror.ReturnErr(err)
	}
	for i := range data.GetChatMessages() {
		if data.GetChatMessage(int64(i)).RecipientId == data.GetUser().Id {
			if err = data.EditChatMessage(int64(i)).MarkAsRead(); err != nil {
				return ferror.ReturnErr(err)
			}
		}
	}
	return nil
}

func ServeUnreadMessages(data state.StateController) error {
	return data.EditChatMessages().GetUnreadMessageIds(data.GetUser().Id)
}
