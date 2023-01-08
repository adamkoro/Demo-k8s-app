package main

import (
	"log"
	"net/http"
	"time"

	"demo-k8s-app/mq-communicator/env"
	logger "demo-k8s-app/mq-communicator/log"
	"demo-k8s-app/mq-communicator/mq"
	endpoints "demo-k8s-app/mq-communicator/v1"

	docs "demo-k8s-app/mq-communicator/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	serviceName = "/producer"
	apiEndpoint = serviceName + "/v1"
)

func init() {
	env.CheckEnvs()

	var err error
	endpoints.Connection, err = mq.ConnectToMq(mq.CreateConnUrl(env.Username, env.Password, env.MqHost, env.Port, env.Vhost))
	errorExist := logger.IsError(err)
	if errorExist {
		logger.ErrorLogger.Printf("%s: %s", err, endpoints.ConnectionFailed)
	}

}

func main() {
	router := gin.Default()

	// Root "/" redirect to default service route
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, serviceName)
	})

	// Api V1 with swagger json host
	docs.SwaggerInfo.BasePath = apiEndpoint
	v1 := router.Group(apiEndpoint)
	{
		v1.GET("/ping", endpoints.Ping)
		v1.GET("/health", endpoints.Health)
		v1.POST("/push", endpoints.SendMessageToMq)
		v1.StaticFile("/docs/swagger.json", "docs/swagger.json")
	}

	url := ginSwagger.URL("http://localhost:8081" + apiEndpoint + "/docs/swagger.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Http server config
	srv := &http.Server{
		Addr:           ":" + env.HttpPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
