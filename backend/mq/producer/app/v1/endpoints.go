package v1

import (
	"context"
	"demo-k8s-app/mq-communicator/env"
	"demo-k8s-app/mq-communicator/messageHandler"
	"demo-k8s-app/mq-communicator/mq"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

// readyness check
func Ping(c *gin.Context) {
	fmt.Fprintf(c.Writer, "pong")
}

// liveness check
func Health(c *gin.Context) {
	conn, err := mq.ConnectToMq(mq.CreateConnUrl(env.Username, env.Password, env.MqHost, env.Port, env.Vhost))
	errorExist := messageHandler.IsError(err)
	if errorExist {
		msg := "Failed to connect to RabbitMQ"
		messageHandler.ErrorLogger.Printf("%s: %s", err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "unhealty",
			"message": msg,
		})
		return
	}
	conn.Close()
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

// Gin HTTP Endpoint for sending data to RabbitMQ
func SendMessageToMq(c *gin.Context) {
	ch, err := mq.CreateChannel(*Connection)
	errorExist := messageHandler.IsError(err)
	if errorExist {
		msg := "Failed to open a channel"
		messageHandler.ErrorLogger.Printf("%s: %s", err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": msg,
		})
		return
	}

	defer mq.CloseChannel(*ch)

	var chName = "Test1"

	q, err := mq.DeclareQueue(*ch, chName)
	errorExist = messageHandler.IsError(err)
	if errorExist {
		msg := "Failed to delcare the queue"
		messageHandler.ErrorLogger.Printf("%s: %s", err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": msg,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte("{id: 1}"),
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
	//mq.CloseConnection(*conn)
	c.JSON(200, gin.H{
		"message": "Successfully published message to the queue",
	})
}
