package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Create conntection url string
func CreateConnUrl(username, password, host, port, vhost string) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s%s", username, password, host, port, vhost)
}

// Create connection to RabbitMQ
func ConnectToMq(connUrl string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(connUrl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Create channel
func CreateChannel(conn amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// Close channel
func CloseChannel(ch amqp.Channel) {
	ch.Close()
}

func CloseConnection(conn amqp.Connection) {
	defer conn.Close()
}

// Create queue
func DeclareQueue(ch amqp.Channel, chName string) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		chName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return q, err
	}
	return q, nil
}
