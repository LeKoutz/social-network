package models

type OnlineQuery struct {
    Response chan map[int64]bool
}

type Hub struct {
    Clients    map[int64]*Client
    Register   chan *Client
    Unregister chan *Client
    Broadcast  chan []byte
    Query      chan OnlineQuery
}

func NewHub() *Hub {
    return &Hub{
        Clients:   make(map[int64]*Client),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
        Broadcast:  make(chan []byte),
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
        case message := <-h.Broadcast:
            for _, client := range h.Clients {
                select {
                case client.Send <- message:
                default:
                    delete(h.Clients, client.UserId)
                    close(client.Send)
                }
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