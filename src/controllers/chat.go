package controllers

import (
	"forum/src/ferror"
	"forum/src/state"
)

func ShowChatHistory(data state.StateController) error {
	user1 := data.GetUser().Id
	user2 := data.GetChatMessage(0).RecipientId
	if user1 == user2 {
		return ferror.ErrorNotFound
	}
	offset := data.GetChatOffset()
	if err := data.EditChatMessages().GetChatHistory(user1, user2, offset); err != nil {
		return err
	}
	for i := range data.GetChatMessages() {
		if data.GetChatMessage(int64(i)).RecipientId == data.GetUser().Id {
			if err := data.EditChatMessage(int64(i)).MarkAsRead(); err != nil {
				return err
			}
		}
	}
	return nil
}

func ServeUnreadMessages(data state.StateController) error {
	return data.EditChatMessages().GetUnreadMessageIds(data.GetUser().Id)
}
