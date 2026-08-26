package controllers

import (
	"errors"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(data state.StateController) error {
	conn, err := upgrader.Upgrade(*data.EditResponse(), data.GetRequest(), nil)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
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
