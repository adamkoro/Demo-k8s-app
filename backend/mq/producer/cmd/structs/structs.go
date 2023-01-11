package structs

type LivenessResponse struct {
	Message string `json:"message" example:"pong"`
}

type ReadinessResponseStatusOk struct {
	Status  string `json:"status" example:"healthy"`
	Message string `json:"message" example:"Sucessfully connected to RabbitMQ"`
}

type ReadinessResponseStatusError struct {
	Status  string `json:"status" example:"unhealthy"`
	Message string `json:"message" example:"Failed to connect to RabbitMQ"`
}

type ResponseMessageOk struct {
	Message string `json:"message" example:"Successfully published message to the queue"`
}

type ResponseMessageError struct {
	Message string `json:"message" example:"Failed to push message to the queue"`
}

type QueueMessage struct {
	Queue string `json:"queue" binding:"required" example:"test"`
	Data  string `json:"data" binding:"required" example:"{ data: 1 }"`
}
