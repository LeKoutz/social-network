package controllers

import (
	"forum/src/models"
	"forum/src/state"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(data state.StateController) error {
	conn, err := upgrader.Upgrade(*data.EditResponse(), data.GetRequest(), nil)
	if err != nil {
		return err
	}
	client := &models.Client{
		Hub:      models.MainHub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserId:   data.GetUser().Id,
		Username: data.GetUser().Username,
	}
	client.Hub.Register <- client
	go client.WritePump()
	go client.ReadPump()
	return nil
}
