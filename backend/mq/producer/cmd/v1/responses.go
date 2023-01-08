package v1

var (
	StatusHealthy       = "healthy"
	StatusUnhealthy     = "unhealthy"
	ConnectionSuccsess  = "Sucessfully connected to RabbitMQ"
	ConnectionFailed    = "Failed to connect to RabbitMQ"
	ChannelCreate       = "Channel successfully created"
	ChannelError        = "Failed to open a channel"
	ChannelClose        = "Channel successfully closed"
	QueueSuccessfull    = "Queue successfully decleared"
	QueueError          = "Failed to delcare the queue"
	MessagePushSuccess  = "Successfully published message to the queue"
	MessagePushFailed   = "Failed to push message to the queue"
	JsonValidationError = "Failed to validate the posted data"
	InvalidHeader       = "Invalid Content-Type header, expected: application/json"
)
