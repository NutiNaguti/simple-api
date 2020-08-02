package services

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type Server struct {
	pattern   string
	messages  []*Message
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
}

// Новый сервер чата
func NewServer(pattern string) *Server {
	messages := []*Message{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{pattern: pattern, messages: messages, clients: clients, addCh: addCh, delCh: delCh, sendAllCh: sendAllCh, doneCh: doneCh, errCh: errCh}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendAll(msg *Message) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) SendPastMessages(c *Client) {
	for _, msg := range s.messages {
		c.Write(msg)
	}
}

func (s *Server) Listen() {

	log.Println("Server listening...")

	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {
		// Добавить нового клиента
		case c := <-s.addCh:
			s.clients[c.id] = c
			s.SendPastMessages(c)
			log.Println("New client")
		// Удалить клиента
		case c := <-s.delCh:
			delete(s.clients, c.id)
			log.Println("Client deleted")
		// Сообщение для всех
		case msg := <-s.sendAllCh:
			s.messages = append(s.messages, msg)
			log.Println("Message for all")
		case err := <-s.errCh:
			log.Println("Error: ", err)
		case <-s.doneCh:
			return
		}
	}
}
