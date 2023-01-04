package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Create conntection url string
func createConnUrl(username, password, host, port string) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", username, password, host, port)
}

// Create connection to RabbitMQ
func connectToMq(connUrl string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(connUrl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Create channel
func createChannel(conn amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// Close channel
func closeChannel(ch amqp.Channel) {
	ch.Close()
}

// Create queue
func declareQueue(ch amqp.Channel, chName string) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		chName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-queue-type": "quorum",
		},
	)
	if err != nil {
		return q, err
	}
	return q, nil
}

// Gin HTTP Endpoint for sending data to RabbitMQ
func SendMessageToMq(c *gin.Context) {
	conn, err := connectToMq(createConnUrl("adamkoro", "legostarwars99", "192.168.1.37", "5672"))
	if err != nil {
		msg := "Failed to connect to RabbitMQ"
		fmt.Printf("%s: %s", err, msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": msg,
		})
	}

	ch, err := createChannel(*conn)
	if err != nil {
		msg := "Failed to open a channel"
		fmt.Printf("%s: %s", err, msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": msg,
		})
	}

	defer closeChannel(*ch)

	var chName = "Test1"
	q, err := declareQueue(*ch, chName)

	if err != nil {
		msg := "Failed to delcare the queue"
		fmt.Printf("%s: %s", err, msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": msg,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello World"),
		},
	)

	if err != nil {
		msg := "Failed to push message to the queue"
		fmt.Printf("%s: %s", err, msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": msg,
		})
	}

	c.JSON(200, gin.H{
		"message": "Successfully published message to the queue",
	})
}
