package v1

import (
	"context"
	"demo-k8s-app/mq-communicator/env"
	"demo-k8s-app/mq-communicator/messageHandler"
	"demo-k8s-app/mq-communicator/mq"
	"demo-k8s-app/mq-communicator/structs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Liveness check
// @BasePath /producer/v1
// @Schemes http
// @Description check liveness
// @Tags health
// @Produce json
// @Success 200 {object} structs.LivenessResponse
// @Router /producer/v1/ping [get]
func Ping(c *gin.Context) {
	var msg structs.LivenessResponse
	msg.Message = "pong"
	c.JSON(http.StatusOK, msg)

}

// Readiness check
// @BasePath /producer/v1
// @Schemes http
// @Description check readyness
// @Tags health
// @Produce json
// @Success 200 {object} structs.ReadinessResponseStatusOk
// @Failure 500 {object} structs.ReadinessResponseStatusError
// @Router /producer/v1/health [get]
func Health(c *gin.Context) {
	var msgOK structs.ReadinessResponseStatusOk
	var msgError structs.ReadinessResponseStatusError

	conn, err := mq.ConnectToMq(mq.CreateConnUrl(env.Username, env.Password, env.MqHost, env.Port, env.Vhost))

	errorExist := messageHandler.IsError(err)
	if errorExist {
		msgError.Status = "unhealthy"
		msgError.Message = "Failed to connect to RabbitMQ"
		messageHandler.ErrorLogger.Printf("%s: %s", err, msgError.Message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, msgError)
		return
	}

	conn.Close()
	msgOK.Status = "healthy"
	c.JSON(http.StatusOK, msgOK)
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
