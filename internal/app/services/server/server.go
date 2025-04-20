package server

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Venukishore-R/chat_app/internal/app/services/client"
)

type Server struct {
	Client    map[string]*client.Client
	BroadCast chan string
	Join      chan *client.Client
	Leave     chan string
	Mu        sync.Mutex
}

// NewServer creates and returns a new Server instance
func NewServer() *Server {
	return &Server{
		Client:    make(map[string]*client.Client),
		BroadCast: make(chan string),
		Join:      make(chan *client.Client),
		Leave:     make(chan string),
	}
}

func (s *Server) Run() {
	fmt.Println("Server is running...")

	for {
		select {
		case client := <-s.Join:
			s.handleJoin(client)

		case clientID := <-s.Leave:
			s.handleLeave(clientID)

		case msg := <-s.BroadCast:
			s.handleBroadcast(msg)
		}
	}
}

func (s *Server) handleJoin(client *client.Client) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Client[client.Id] = client
	fmt.Printf("Client %s joined the chat\n", client.Id)
}

func (s *Server) handleLeave(clientID string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if c, exists := s.Client[clientID]; exists {
		close(c.MsgChan)
		delete(s.Client, clientID)
		fmt.Printf("Client %s left the chat\n", clientID)
	} else {
		fmt.Printf("Client %s not found to leave\n", clientID)
	}
}

func (s *Server) handleBroadcast(message string) {
	fmt.Println("Broadcasting message...")

	s.Mu.Lock()
	defer s.Mu.Unlock()

	for _, client := range s.Client {
		select {
		case client.MsgChan <- message:
		default:
			fmt.Printf("Client %s is not ready to receive the message\n", client.Id)
		}
	}
	fmt.Println("Broadcast completed")
}

func (s *Server) JoinRoom(id string) *client.Client {
	newClient := client.NewClient(id)

	select {
	case s.Join <- newClient:
		return newClient
	case <-time.After(2 * time.Second):
		fmt.Printf("Timeout while joining client %s\n", id)
		return nil
	}
}

func (s *Server) LeaveRoom(id string) error {
	select {
	case s.Leave <- id:
		return nil
	case <-time.After(2 * time.Second):
		return fmt.Errorf("timeout while trying to remove client %s", id)
	}
}

func (s *Server) SendMessage(id, msg string) error {
	formatted := fmt.Sprintf("[%s]: %s", id, msg)

	select {
	case s.BroadCast <- formatted:
		return nil
	case <-time.After(2 * time.Second):
		return errors.New("timeout while sending broadcast message")
	}
}
