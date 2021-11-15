package main


//type hub struct {
//	rooms map[string]map[*connection]bool
//	broadcast chan message
//	unregister chan subscription
//}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan Subscription

	// Unregister requests from clients.
	unregister chan Subscription
}

type Subscription struct {
	client *Client
	roomId string
}

type Message struct {
	data []byte
	roomId string
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan Subscription),
		unregister: make(chan Subscription),
		clients:    make(map[string]map[*Client]bool),
	}
}

func (h *Hub) run(){
	for {
		select {
		case s := <-h.register:
			clientsWithRoomId := h.clients[s.roomId]
			if clientsWithRoomId == nil {
				clientsWithRoomId = make(map[*Client]bool)
				h.clients[s.roomId] = clientsWithRoomId
			}
			h.clients[s.roomId][s.client] = true

		case s := <-h.unregister:
			clientsWithRoomId := h.clients[s.roomId]
			if clientsWithRoomId != nil {
				if _, ok := clientsWithRoomId[s.client]; ok{
					delete(clientsWithRoomId, s.client)
					close(s.client.send)
					if len(clientsWithRoomId) == 0 {
						delete(h.clients, s.roomId)
					}
				}
			}

		case m := <-h.broadcast:
			clientsWithRoomId := h.clients[m.roomId]
			for c := range clientsWithRoomId{
				select{
				case c.send <- m.data:
				default:
					close(c.send)
					delete(clientsWithRoomId, c)
					if len(clientsWithRoomId) == 0{
						delete(h.clients, m.roomId)
					}
				}
			}
		}
	}
}