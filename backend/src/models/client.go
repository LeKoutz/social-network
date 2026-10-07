package models

import (
	"encoding/json"
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserId   int64
	Username string
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			(&ferror.Error{}).Consume(ferror.ReturnErr(err)).LogError()
			break
		}
		var incoming WsMessage
		if err := json.Unmarshal(message, &incoming); err != nil {
			(&ferror.Error{}).Consume(ferror.ReturnErr(err)).LogError()
			continue
		}
		switch incoming.Type {
		case "chat-message":
			err = c.handleChatMessage(incoming.Payload)
		case "group-message":
			err = c.handleGroupMessage(incoming.Payload)
		case "message-read":
			err = c.handleMessageRead(incoming.Payload)
		default:
			err = ferror.ErrorBadRequest
		}
		if err != nil {
			(&ferror.Error{}).Consume(ferror.ReturnErr(err)).LogError()
		}
	}
}

func (c *Client) handleChatMessage(payload json.RawMessage) error {
	var p struct {
		RecipientId int64  `json:"recipientId"`
		Body        string `json:"body"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return ferror.ReturnErr(err)
	}
	if p.RecipientId == c.UserId {
		return ferror.ReturnErr(ferror.ErrorNotFound)
	}
	if len(strings.TrimSpace(p.Body)) == 0 {
		return ferror.ReturnErr(ferror.ErrorChatMessageEmpty)
	}
	inv := FollowRequestType{
		InvitationRowType: db.InvitationRowType{
			FromUserId: c.UserId,
			ToUserId:   p.RecipientId,
		},
	}
	allowed, err := inv.HasFollowRelation()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if !allowed {
		return ferror.ReturnErr(ferror.ErrorRecipientNotFollowed)
	}
	msg := ChatMessageType{
		ChatMessageRowType: db.ChatMessageRowType{
			SenderId:    c.UserId,
			RecipientId: p.RecipientId,
			Body:        p.Body,
			Timestamp:   utils.GetCurrentTimestamp(),
		},
		SenderUsername: c.Username,
	}
	msg.Id, err = msg.Add()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	c.Hub.Transmit <- msg
	return nil
}

func (c *Client) handleGroupMessage(payload json.RawMessage) error {
	var p struct {
		GroupId int64  `json:"groupId"`
		Body    string `json:"body"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return ferror.ReturnErr(err)
	}
	if len(strings.TrimSpace(p.Body)) == 0 {
		return ferror.ReturnErr(ferror.ErrorChatMessageEmpty)
	}
	group := GroupType{}
	group.Id = p.GroupId
	isMember, err := group.IsMember(c.UserId)
	if err != nil {
		return ferror.ReturnErr(err)
	}
	if !isMember {
		return ferror.ReturnErr(ferror.ErrorPermissionDenied)
	}
	msg := GroupMessageType{
		GroupMessageRowType: db.GroupMessageRowType{
			GroupId:   p.GroupId,
			SenderId:  c.UserId,
			Body:      p.Body,
			Timestamp: utils.GetCurrentTimestamp(),
		},
		SenderUsername: c.Username,
	}
	msg.Id, err = msg.Add()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	msg.MemberIds, err = group.GetMemberIds()
	if err != nil {
		return ferror.ReturnErr(err)
	}
	c.Hub.TransmitGroup <- msg
	return nil
}

func (c *Client) handleMessageRead(payload json.RawMessage) error {
	var p struct {
		Id      int64 `json:"Id"`
		GroupId int64 `json:"GroupId"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return ferror.ReturnErr(err)
	}
	if p.GroupId != 0 {
		message := GroupMessageType{
			GroupMessageRowType: db.GroupMessageRowType{
				Id:      p.Id,
				GroupId: p.GroupId,
			},
		}
		return message.MarkAsRead(c.UserId)
	}
	message := ChatMessageType{
		ChatMessageRowType: db.ChatMessageRowType{
			Id:          p.Id,
			RecipientId: c.UserId,
		},
	}
	return message.MarkAsRead()
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			(&ferror.Error{}).Consume(ferror.ReturnErr(err)).LogError()
			break
		}
	}
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}
