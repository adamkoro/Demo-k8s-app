package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/go-swagger/go-swagger/swag"
)

func main() {
	// Create a Gin router
	r := gin.Default()

	// Serve the Swagger UI static assets
	r.Static("/swagger-ui", "./swagger-ui/dist")

	// Create a swagger API documentation server
	api := middleware.New(swag.New("1.0.0", "API Documentation"))

	// Load the Swagger specifications for each service
	api.Spec("http://192.168.1.100:8081/producer/v1/docs/swagger.json")

	// Serve the combined Swagger specification
	r.GET("/swagger.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", api.Spec())
	})

	// Start the server
	r.Run(":8080")
}
