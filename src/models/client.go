package models

import (
	"errors"
	"encoding/json"
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"

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
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			(&ferror.Error{}).Consume(err).LogError()
			break
		}
		var incoming WsMessage
		if err := json.Unmarshal(message, &incoming); err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			(&ferror.Error{}).Consume(err).LogError()
			continue
		}
		switch incoming.Type {
		case "chat-message":
			var p struct {
				RecipientId int64  `json:"recipientId"`
				Body        string `json:"body"`
			}
			if err := json.Unmarshal(incoming.Payload, &p); err != nil {
				if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
				(&ferror.Error{}).Consume(err).LogError()
				continue
			}
			timestampString := utils.GetCurrentTimestamp()
			msg := ChatMessageType{
				ChatMessageRowType: db.ChatMessageRowType{
					SenderId:        c.UserId,
					RecipientId:     p.RecipientId,
					Body:            p.Body,
					TimestampString: timestampString,
					SenderUsername:  c.Username,
				},
			}
			msg.Id, err = msg.Add()
			if err != nil {
				if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
				(&ferror.Error{}).Consume(err).LogError()
				continue
			}
			timestampTime, err := utils.ConvertStringToTime(msg.TimestampString)
			if err != nil {
				if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
				(&ferror.Error{}).Consume(err).LogError()
				continue
			}
			msg.TimestampString = utils.ConvertTimeToString(timestampTime)
			c.Hub.Transmit <- msg
		case "message-read":
			message := ChatMessageType{}
			if err := json.Unmarshal(incoming.Payload, &message); err != nil {
				if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
				(&ferror.Error{}).Consume(err).LogError()
				continue
			}
			if err := message.MarkAsRead(); err != nil {
				if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
				(&ferror.Error{}).Consume(err).LogError()
			}
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			(&ferror.Error{}).Consume(err).LogError()
			break
		}
	}
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}
