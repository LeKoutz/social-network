package controllers

import (
	"forum/src/models"
	"github.com/gorilla/websocket"
)

var Hub *models.Hub

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(data models.ResponseStruct) {
	conn, err := upgrader.Upgrade(data.Response, data.Request, nil)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	client := &models.Client{Hub: Hub, Conn: conn, Send: make(chan []byte, 256), UserId: data.User.Id}
	client.Hub.Register <- client
	go client.WritePump()
	go client.ReadPump()
}

func getUsers(data models.ResponseStruct) {
	users, err := models.GetAllUsers()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.Users = users
	data.WriteResponse() // no view necessary
}