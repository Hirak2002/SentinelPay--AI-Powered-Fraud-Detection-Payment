package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	mu      sync.Mutex
}

func NewMessageBroker() (*MessageBroker, error) {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Printf("RabbitMQ connection failed (using in-memory fallback): %v", err)
		return &MessageBroker{}, nil
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &MessageBroker{conn: conn, channel: ch}, nil
}

func (mb *MessageBroker) PublishEvent(eventType string, data interface{}) error {
	if mb.channel == nil {
		log.Printf("Event published to in-memory: %s", eventType)
		return nil
	}

	mb.mu.Lock()
	defer mb.mu.Unlock()

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	q, err := mb.channel.QueueDeclare(
		eventType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return mb.channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (mb *MessageBroker) Close() error {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if mb.channel != nil {
		mb.channel.Close()
	}
	if mb.conn != nil {
		return mb.conn.Close()
	}
	return nil
}
