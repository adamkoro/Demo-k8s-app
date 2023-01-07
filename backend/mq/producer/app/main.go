package main

import (
	"log"
	"net/http"
	"time"

	"demo-k8s-app/mq-communicator/env"
	"demo-k8s-app/mq-communicator/messageHandler"
	"demo-k8s-app/mq-communicator/mq"
	endpoints "demo-k8s-app/mq-communicator/v1"

	docs "demo-k8s-app/mq-communicator/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var serviceEndpoint = "/producer/v1"

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
	docs.SwaggerInfo.BasePath = serviceEndpoint
	v1 := router.Group(serviceEndpoint)
	{
		v1.GET("/ping", endpoints.Ping)
		v1.GET("/health", endpoints.Health)
		v1.GET("/push", endpoints.SendMessageToMq)
		v1.StaticFile("/docs/swagger.json", "docs/swagger.json")
	}

	url := ginSwagger.URL("http://localhost:8081" + serviceEndpoint + "/docs/swagger.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Http server config
	srv := &http.Server{
		Addr:           ":8081",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
