package v1

import (
	"context"
	"demo-k8s-app/mq-communicator/env"
	"demo-k8s-app/mq-communicator/messageHandler"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Create conntection url string
func createConnUrl(username, password, host, port, vhost string) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s%s", username, password, host, port, vhost)
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
	conn, err := connectToMq(createConnUrl(env.Username, env.Password, env.MqHost, env.Port, env.Vhost))
	errorExist := messageHandler.IsError(err)
	if errorExist {
		msg := "Failed to connect to RabbitMQ"
		messageHandler.ErrorLogger.Printf("%s: %s", err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": msg,
		})
		return
	}

	ch, err := createChannel(*conn)
	errorExist = messageHandler.IsError(err)
	if errorExist {
		msg := "Failed to open a channel"
		messageHandler.ErrorLogger.Printf("%s: %s", err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": msg,
		})
		return
	}

	defer closeChannel(*ch)

	var chName = "Test1"

	q, err := declareQueue(*ch, chName)
	errorExist = messageHandler.IsError(err)
	if errorExist {
		msg := "Failed to delcare the queue"
		messageHandler.ErrorLogger.Printf("%s: %s", err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": msg,
		})
		return
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

	errorExist = messageHandler.IsError(err)
	if errorExist {
		msg := "Failed to push message to the queue"
		messageHandler.ErrorLogger.Printf("%s: %s", err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Successfully published message to the queue",
	})
}
