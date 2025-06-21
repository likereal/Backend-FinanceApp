package main

import (
	"context"
	"log"
	"github.com/coder/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	Hub  *Hub
}

func (c *Client) ReadLoop(ctx context.Context) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close(websocket.StatusNormalClosure, "read done")
	}()

	for {
		_, msg, err := c.Conn.Read(ctx)
		if err != nil {
			break
		}

		// Encrypt before broadcasting
		encryptedMsg, err := Encrypt(string(msg))
		if err != nil {
			log.Println("encryption error:", err)
			continue
		}

		// Broadcast encrypted message
		c.Hub.Broadcast <- []byte(encryptedMsg)
	}
}

func (c *Client) WriteLoop(ctx context.Context) {
	defer c.Conn.Close(websocket.StatusNormalClosure, "write done")

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.Send:
			if !ok {
				c.Conn.Close(websocket.StatusInternalError, "hub closed")
				return
			}

			// Decrypt message before sending to client
			decryptedMsg, err := Decrypt(string(msg))
			if err != nil {
				log.Println("decryption error:", err)
				continue
			}

			err = c.Conn.Write(ctx, websocket.MessageText, []byte(decryptedMsg))
			if err != nil {
				log.Println("write error:", err)
				return
			}
		}
	}
}
