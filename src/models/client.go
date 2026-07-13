package models

import (
	"encoding/json"
	"forum/src/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
    Hub    *Hub
    Conn   *websocket.Conn
    Send   chan []byte
    UserId int64
	Username string
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	type incomingMessage struct {
		RecipientId int64  `json:"recipientId"`
		Body        string `json:"body"`
	}
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			(&Error{}).Consume(err).LogError()
			break
		}
		var p incomingMessage
		if err := json.Unmarshal(message, &p); err != nil {
			(&Error{}).Consume(err).LogError()
        	continue
    	}
		timestampString := utils.GetCurrentTimestamp()
		msg := ChatMessage{
			SenderId:    c.UserId,
			RecipientId: p.RecipientId,
			Body:        p.Body,
			TimestampString: timestampString,
			SenderUsername: c.Username,
		}
		msg.Id, err = msg.Add()
		if err != nil {
			(&Error{}).Consume(err).LogError()
			continue
		}
		timestampTime, err := utils.ConvertStringToTime(msg.TimestampString)
		if err != nil {
			(&Error{}).Consume(err).LogError()
			continue
		}
		msg.TimestampString = utils.ConvertTimeToString(timestampTime)
		c.Hub.Broadcast <- message
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			(&Error{}).Consume(err).LogError()
			break
		}
	}
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}
	