package main

import (
	"log"
	"net/http"

	endpoints "demo-k8s-app/mq-communicator/v1"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Api V1
	v1 := router.Group("/v1")
	{
		v1.GET("/ping", endpoints.Ping)
		v1.GET("/push", endpoints.SendMessageToMq)
	}

	// Http server config
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
