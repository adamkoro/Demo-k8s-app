package structs

type LivenessResponse struct {
	Message string `json:"message" example:"pong"`
}

type ReadinessResponseStatusOk struct {
	Status string `json:"status" example:"healthy"`
}

type ReadinessResponseStatusError struct {
	Status  string `json:"status" example:"unhealthy"`
	Message string `json:"message" example:"Failed to connect to RabbitMQ"`
}
