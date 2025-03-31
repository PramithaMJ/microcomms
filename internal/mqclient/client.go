package mqclient

import (
	"fmt"
	"log"
	"time"
)

// MessageQueue is a simple message queue interface
type MessageQueue struct {
	channel chan string
}

// NewMessageQueue creates a new MessageQueue instance
func NewMessageQueue(size int) *MessageQueue {
	return &MessageQueue{
		channel: make(chan string, size),
	}
}

// SendMessage sends a message to the queue
func (mq *MessageQueue) SendMessage(message string) error {
	select {
	case mq.channel <- message:
		log.Printf("Message sent: %v", message)
		return nil
	case <-time.After(2 * time.Second):
		return fmt.Errorf("timeout sending message")
	}
}

// ReceiveMessage receives a message from the queue
func (mq *MessageQueue) ReceiveMessage() (string, error) {
	select {
	case message := <-mq.channel:
		log.Printf("Message received: %v", message)
		return message, nil
	case <-time.After(2 * time.Second):
		return "", fmt.Errorf("timeout receiving message")
	}
}
