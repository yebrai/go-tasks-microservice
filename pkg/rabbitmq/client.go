package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type ClientConfig struct {
	URL      string // "amqp://guest:guest@localhost:5672/"
	Exchange string // "task-events"
	Queue    string // "task-events-queue"
}

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  ClientConfig
}

func NewClient(config ClientConfig) (*Client, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// 1. DECLARAR EXCHANGE
	err = channel.ExchangeDeclare(
		config.Exchange, // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	// 2. DECLARAR QUEUE
	_, err = channel.QueueDeclare(
		config.Queue, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	// 3. BINDING QUEUE AL EXCHANGE
	err = channel.QueueBind(
		config.Queue,    // queue name
		"task.*",        // routing key pattern (captura task.created, task.completed, etc.)
		config.Exchange, // exchange
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	return &Client{
		conn:    conn,
		channel: channel,
		config:  config,
	}, nil
}

func (c *Client) Publish(routingKey string, body []byte) error {
	return c.channel.Publish(
		c.config.Exchange, // exchange
		routingKey,        // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Persistir mensajes
		},
	)
}

func (c *Client) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}
