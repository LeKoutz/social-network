package controllers

import (
	"forum/src/ferror"
	"forum/src/state"
)

func ShowChatHistory(data state.StateController) error {
	user1 := data.GetUser().Id
	user2 := data.GetChatMessage(0).RecipientId
	if user1 == user2 {
		err := ferror.ErrorNotFound
		data.SetErrorConsume(err)
		return err
	}
	offset := data.GetChatOffset()
	if err := data.EditChatMessages().GetChatHistory(user1, user2, offset); err != nil {
		data.SetErrorConsume(err)
		return err
	}
	for i := range data.GetChatMessages() {
		if data.GetChatMessage(int64(i)).RecipientId == data.GetUser().Id {
			data.EditChatMessage(int64(i)).MarkAsRead()
		}
	}
	return nil
}

func ServeUnreadMessages(data state.StateController) error {
	if err := data.EditChatMessages().GetUnreadMessageIds(data.GetUser().Id); err != nil {
		data.SetErrorConsume(err)
		return err
	}
	return nil
}
