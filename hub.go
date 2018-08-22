// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Envelope contains routing information for a message.
type Envelope struct {
	// id of sending peer
	srcGUID string
	// id of target peer to receive message
	dstGUID string
	// message data
	message []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Envelope

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	noClients chan uint32

	key uint32

	serverGUID string
}

func newHub(key uint32, noClients chan uint32) *Hub {
	return &Hub{
		broadcast:  make(chan Envelope),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		noClients:  noClients,
		key:        key,
	}
}

func (h *Hub) count() int {
	return len(h.clients)
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case envelope := <-h.broadcast:
			for client := range h.clients {
                if(envelope.srcGUID == client.GUID) {
                    continue
                }
                if envelope.dstGUID == "ALL" || envelope.dstGUID == client.GUID {
					select {
					case client.send <- envelope:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
		if len(h.clients) == 0 {
			h.noClients <- h.key
			break
		}
	}
}
