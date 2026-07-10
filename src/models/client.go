package models

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
    Hub    *Hub
    Conn   *websocket.Conn
    Send   chan []byte
    UserId int64
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
			break
		}
		var p incomingMessage
		if err := json.Unmarshal(message, &p); err != nil {
        	continue
    	}
		msg := ChatMessage{
			SenderId:    c.UserId,
			RecipientId: p.RecipientId,
			Body:        p.Body,
		}
		if err := msg.Add(); err != nil {
			continue
		}
		c.Hub.Broadcast <- message
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}
	