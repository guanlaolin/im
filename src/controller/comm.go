package controller

import (
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

type Client struct {
	UUID uuid.UUID
	Conn *websocket.Conn
}

type ClientManager struct {
	Clients []*Client
	UUIDs   map[uuid.UUID]bool

	Login   chan uuid.UUID
	Upgrade chan *Client
}

var Manager = ClientManager{}

func NewClient(uuid uuid.UUID, conn *websocket.Conn) *Client {
	return &Client{uuid, conn}
}

func (manager *ClientManager) do() {
	for {
		select {
		case uuid := <-manager.Login:
			manager.UUIDs[uuid] = true
		case client := <-manager.Upgrade:
			manager.Clients = append(manager.Clients, client)
			delete(manager.UUIDs, client.UUID)
		}
	}
}
