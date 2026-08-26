package state

import "forum/src/models"

func (r *State) SetChatMessages() {
	var message models.ChatMessageType
	r.EditUser().ChatMessages = models.ChatMessagesType{message}
}

func (r *State) EditChatMessages() *models.ChatMessagesType {
	if r.GetUser().ChatMessages == nil {
		r.SetChatMessages()
	}
	return &r.EditUser().ChatMessages
}

func (r *State) GetChatMessages() models.ChatMessagesType {
	if r.GetUser().ChatMessages == nil {
		r.EditUser().ChatMessages = models.ChatMessagesType{}
	}
	return *r.EditChatMessages()
}

func (r *State) EditChatMessage(index int64) *models.ChatMessageType {
	msgs := r.EditChatMessages()
	if index < 0 || int(index) >= len(*msgs) {
		return nil
	}
	return &(*msgs)[index]
}

func (r *State) GetChatMessage(index int64) models.ChatMessageType {
	return *r.EditChatMessage(index)
}