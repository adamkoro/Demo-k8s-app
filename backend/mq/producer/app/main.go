package main

import (
	"log"
	"net/http"
	"time"

	"demo-k8s-app/mq-communicator/env"
	"demo-k8s-app/mq-communicator/messageHandler"
	"demo-k8s-app/mq-communicator/mq"
	endpoints "demo-k8s-app/mq-communicator/v1"

	"github.com/gin-gonic/gin"
)

func init() {
	env.CheckEnvs()

	var err error
	endpoints.Connection, err = mq.ConnectToMq(mq.CreateConnUrl(env.Username, env.Password, env.MqHost, env.Port, env.Vhost))
	errorExist := messageHandler.IsError(err)
	if errorExist {
		msg := "Failed to connect to RabbitMQ"
		messageHandler.ErrorLogger.Printf("%s: %s", err, msg)
	}
}

func main() {
	router := gin.Default()

	// Defautl route
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/producer")
	})
	// Api V1
	v1 := router.Group("/:servicename/v1")
	{
		v1.GET("/ping", endpoints.Ping)
		v1.GET("/health", endpoints.Health)
		v1.GET("/push", endpoints.SendMessageToMq)
	}

	// Http server config
	srv := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
