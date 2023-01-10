package main

import (
	"fmt"
	"net/http"
	"time"

	endpoints "demo-k8s-app/mq-communicator/api/v1"
	docs "demo-k8s-app/mq-communicator/docs"
	"demo-k8s-app/mq-communicator/env"
	logger "demo-k8s-app/mq-communicator/log"
	"demo-k8s-app/mq-communicator/mq"

	"github.com/gin-gonic/gin"
)

var (
	serviceName = "/producer"
	apiEndpoint = serviceName + "/v1"
)

func init() {
	logger.InfoLogger.Println("Init phase started")
	env.CheckEnvs()
	var err error
	endpoints.Connection, err = mq.ConnectToMq(mq.CreateConnUrl(env.Username, env.Password, env.MqHost, env.Port, env.Vhost))
	errorExist := logger.IsError(err)
	if errorExist {
		logger.ErrorLogger.Fatalf("%s: %s", err, endpoints.ConnectionFailed)
	}
	logger.InfoLogger.Println("RabbitMQ connection successfully established")
	logger.InfoLogger.Println("Init phase finished")
}

func main() {
	gin.DisableConsoleColor()
	logger.InfoLogger.Println("Main phase started")

	router := gin.New()

	// Custom logger: [HTTP] 2023/01/08 18:47:25 | Code: 404 | Method: GET | IP: 127.0.0.1 | Path: /producer/v1/test
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[HTTP] %s | Code: %d | Method: %s | IP: %s | Path: %s\n",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.StatusCode,
			param.Method,
			param.ClientIP,
			param.Path,
		)
	}))

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

	// Http server config
	srv := &http.Server{
		Addr:         ":" + env.HttpPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	logger.InfoLogger.Println("Setup http routers/endpoints were successfully finish")
	logger.InfoLogger.Println("Server start at port: " + env.HttpPort)
	err := srv.ListenAndServe()

	errorExist := logger.IsError(err)
	if errorExist {
		logger.ErrorLogger.Fatalf("%s: %s", err, "Could not start the webserver")
	}
}
