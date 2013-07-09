package main

import (
	"code.google.com/p/go.net/websocket"
)

type Connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan string
}

func (c *Connection) reader() {
	for {
		var message [20]byte
		n, err := c.ws.Read(message[:])
		if err != nil {
			break
		}
		h.broadcast <- string(message[:n])
	}
	c.ws.Close()
}

func (c *Connection) writer() {
	for message := range c.send {
		err := websocket.Message.Send(c.ws, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}
