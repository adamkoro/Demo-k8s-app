package v1

import (
	"context"
	"demo-k8s-app/mq-communicator/env"
	logger "demo-k8s-app/mq-communicator/log"
	"demo-k8s-app/mq-communicator/mq"
	"demo-k8s-app/mq-communicator/structs"
	"fmt"
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

	errorExist := logger.IsError(err)
	if errorExist {
		msgError.Status = StatusUnhealthy
		msgError.Message = ConnectionFailed
		logger.ErrorLogger.Printf("%s: %s", err, msgError.Message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, msgError)
		return
	}

	conn.Close()
	msgOK.Status = StatusHealthy
	msgOK.Message = ConnectionSuccsess
	c.JSON(http.StatusOK, msgOK)
}

// HTTP Endpoint for sending data to RabbitMQ
// @BasePath /producer/v1
// @Schemes http
// @Description sending data to RabbitMQ
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseMessageOk
// @Failure 500 {object} structs.ResponseMessageError
// @Router /producer/v1/push [post]
func SendMessageToMq(c *gin.Context) {
	var msgOk structs.ResponseMessageOk
	var msgError structs.ResponseMessageError
	var msg structs.QueueMessage

	// Check application type header
	if c.GetHeader("Content-Type") != "application/json" {
		c.AbortWithStatusJSON(http.StatusBadRequest, InvalidHeader)
		return
	}

	// Check body struct
	err := c.ShouldBindJSON(&msg)
	errorExist := logger.IsError(err)
	if errorExist {
		msgError.Message = JsonValidationError
		c.AbortWithStatusJSON(http.StatusBadRequest, msgError)
		return
	}
	if msg.Queue == "" {
		msgError.Message = "Missing Queue field or empty"
		c.AbortWithStatusJSON(http.StatusBadRequest, msgError.Message)
		return
	}
	if msg.Data == "" {
		msgError.Message = "Missing Data field or empty"
		c.AbortWithStatusJSON(http.StatusBadRequest, msgError.Message)
		return
	}

	fmt.Println(env.Queues)
	// Check requested queue name
	for _, queueName := range env.Queues {
		if msg.Queue == queueName {
			break
		} else {
			msgError.Message = "Requested qeueu not found in avaiable list"
			c.AbortWithStatusJSON(http.StatusBadRequest, msgError.Message)
			return
		}
	}

	ch, err := mq.CreateChannel(*Connection)
	errorExist = logger.IsError(err)
	if errorExist {
		msgError.Message = ChannelError
		logger.ErrorLogger.Printf("%s: %s", err, msgError.Message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, msgError)
		return
	}
	msgOk.Message = ChannelCreate
	logger.InfoLogger.Printf("%s", msgOk.Message)

	defer mq.CloseChannel(*ch)

	q, err := mq.DeclareQueue(*ch, msg.Queue)
	errorExist = logger.IsError(err)
	if errorExist {
		msgError.Message = QueueError
		logger.ErrorLogger.Printf("%s: %s", err, msgError.Message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, msgError)
		return
	}
	msgOk.Message = QueueSuccessfull
	logger.InfoLogger.Printf("%s", msgOk.Message)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(msg.Data),
		},
	)

	errorExist = logger.IsError(err)
	if errorExist {
		msgError.Message = MessagePushFailed
		logger.ErrorLogger.Printf("%s: %s", err, msgError.Message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, msgError)
		return
	}
	msgOk.Message = MessagePushSuccess
	c.JSON(http.StatusAccepted, msgOk)
}
