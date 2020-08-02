package services

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
)

var maxId = 0

const channelBufSize = 100

type Client struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *Message
	doneCh chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil {
		panic("ws can't be nil")
	}
	if server == nil {
		panic("server can't be nil")
	}

	maxId++
	ch := make(chan *Message, channelBufSize)
	doneCh := make(chan bool)

	return &Client{
		id:     maxId,
		ws:     ws,
		server: server,
		ch:     ch,
		doneCh: doneCh,
	}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg *Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("Client %d is disconnected. ", c.id)
		c.server.Err(err)
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *Client) listenWrite() {
	for {
		select {
		// отправить сообщение клиенту
		case msg := <-c.ch:
			websocket.JSON.Send(c.ws, msg)

		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true
			return
		}
	}
}

func (c *Client) listenRead() {
	for {
		select {
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true
			return

		default:
			var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				c.server.SendAll(&msg)
			}
		}
	}
}
