package ws

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID    uint
	Role      string
	Conn      *websocket.Conn
	Hub       *Hub
	Send      chan []byte
	closeSend sync.Once
}

type Hub struct {
	mu            sync.RWMutex
	CustomerConns map[uint]*Client
	AdminConns    map[uint]*Client
	Register      chan *Client
	Unregister    chan *Client
}

var DefaultHub *Hub

func NewHub() *Hub {
	return &Hub{
		CustomerConns: make(map[uint]*Client),
		AdminConns:    make(map[uint]*Client),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Hub panic recovered: %v", r)
				}
			}()

			select {
			case client := <-h.Register:
				h.mu.Lock()
				if client.Role == "admin" || client.Role == "service" {
					if existing, ok := h.AdminConns[client.UserID]; ok {
						existing.closeSend.Do(func() { close(existing.Send) })
					}
					h.AdminConns[client.UserID] = client
				} else {
					if existing, ok := h.CustomerConns[client.UserID]; ok {
						existing.closeSend.Do(func() { close(existing.Send) })
					}
					h.CustomerConns[client.UserID] = client
				}
				h.mu.Unlock()
				log.Printf("WebSocket connected: userID=%d role=%s", client.UserID, client.Role)

			case client := <-h.Unregister:
				h.mu.Lock()
				if client.Role == "admin" || client.Role == "service" {
					if existing, ok := h.AdminConns[client.UserID]; ok && existing == client {
						delete(h.AdminConns, client.UserID)
					}
				} else {
					if existing, ok := h.CustomerConns[client.UserID]; ok && existing == client {
						delete(h.CustomerConns, client.UserID)
					}
				}
				h.mu.Unlock()
				client.closeSend.Do(func() { close(client.Send) })
				log.Printf("WebSocket disconnected: userID=%d role=%s", client.UserID, client.Role)
			}
		}()
	}
}

func (h *Hub) SendToUser(userID uint, msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.CustomerConns[userID]; ok {
		select {
		case client.Send <- data:
		default:
		}
	}
}

func (h *Hub) SendToAdmins(msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.AdminConns {
		select {
		case client.Send <- data:
		default:
		}
	}
}

func (h *Hub) SendToAll(msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.CustomerConns {
		select {
		case client.Send <- data:
		default:
		}
	}
	for _, client := range h.AdminConns {
		select {
		case client.Send <- data:
		default:
		}
	}
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump(handler func(*Client, []byte)) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error (userID=%d): %v", c.UserID, err)
			return
		}
		handler(c, message)
	}
}
