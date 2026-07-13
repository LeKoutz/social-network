package models

import (
    "encoding/json"
)

type WsMessage struct {
    Type    string          `json:"type"`
    Payload json.RawMessage `json:"payload"`
}

type OnlineQuery struct {
    Response chan map[int64]bool
}

type Hub struct {
    Clients    map[int64]*Client
    Register   chan *Client
    Unregister chan *Client
    Broadcast  chan []byte
    Transmit   chan ChatMessage
    Query      chan OnlineQuery
}

func NewHub() *Hub {
    return &Hub{
        Clients:   make(map[int64]*Client),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
        Broadcast:  make(chan []byte),
        Transmit:   make(chan ChatMessage),
        Query:      make(chan OnlineQuery),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.Register:
            h.Clients[client.UserId] = client
        case client := <-h.Unregister:
            if _, ok := h.Clients[client.UserId]; ok {
                delete(h.Clients, client.UserId)
                close(client.Send)
            }
        case alert := <-h.Broadcast:
            for _, client := range h.Clients {
                select {
                case client.Send <- alert:
                default:
                    delete(h.Clients, client.UserId)
                    close(client.Send)
                }
            }
        case msg := <-h.Transmit:
            payload , err := json.Marshal(msg)
            if err != nil {
                (&Error{}).Consume(err).LogError()
                continue
            }
            msgBytes, err := json.Marshal(WsMessage{
                Type:    "chat_message",
                Payload: json.RawMessage(payload),
            })
            if err != nil {
                (&Error{}).Consume(err).LogError()
                continue
            }
            if sender, ok := h.Clients[msg.SenderId]; ok {
                sender.Send <- msgBytes
            }
            if recipient, ok := h.Clients[msg.RecipientId]; ok {
                recipient.Send <- msgBytes
            }
        case query := <-h.Query:
            onlineUsers := make(map[int64]bool)
            for _, client := range h.Clients {
                onlineUsers[client.UserId] = true
            }
            query.Response <- onlineUsers
        }
    }
}

func (h *Hub) GetOnlineUsers() map[int64]bool {
    query := OnlineQuery{Response: make(chan map[int64]bool)}
    h.Query <- query
    return <-query.Response
}