package v1

var (
	StatusHealthy       = "healthy"
	StatusUnhealthy     = "unhealthy"
	ConnectionSuccess   = "Successfully connected to RabbitMQ"
	ConnectionFailed    = "Failed to connect to RabbitMQ"
	ChannelCreate       = "Channel successfully created"
	ChannelError        = "Failed to open a channel"
	ChannelClose        = "Channel successfully closed"
	QueueSuccessfully   = "Queue successfully declared"
	QueueError          = "Failed to declare the queue"
	QueueNotFound       = "Requested queue not found in available list"
	MessagePushSuccess  = "Successfully published message to the queue"
	MessagePushFailed   = "Failed to push message to the queue"
	JsonValidationError = "Failed to validate the posted data"
	InvalidHeader       = "Invalid Content-Type header, expected: application/json"
	MethodNotAllowed    = "Method not allowed"
)
