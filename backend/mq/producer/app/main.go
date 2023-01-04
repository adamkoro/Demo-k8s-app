package main

import (
	"log"
	"net/http"

	docs "demo-k8s-app/mq-communicator/docs"
	"demo-k8s-app/mq-communicator/env"
	endpoints "demo-k8s-app/mq-communicator/v1"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title MQ-Producer Swagger
// @version 1.0
// @description RabbitMQ Producer(sender) API
func main() {
	env.CheckEnvs()
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/v1"

	// Api V1
	v1 := router.Group("/v1")
	{
		v1.GET("/ping", endpoints.Ping)
		v1.GET("/push", endpoints.SendMessageToMq)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Http server config
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
