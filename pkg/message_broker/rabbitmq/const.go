package rabbitmq

const (
	AmqpUrl = "amqp://%s:%s@%s:%s/"

	ErrConnFailed         = "Failed to connect to RabbitMQ"
	ErrChanFailed         = "Failed to open a channel"
	ErrQueueFailed        = "Failed to declare a queue"
	ErrConsumerFailed     = "Failed to register a consumer"
	ErrProducerFailed     = "Failed to publish a message"
	ErrClosingConnection  = "Error closing RabbitMQ connection"
	ErrClosingChannel     = "Error closing RabbitMQ channel"
	ErrGetChannelFromPool = "Error getting channel from pool"

	MsgConnSucceed         = "Connected to RabbitMQ"
	MsgChanCreated         = "Channel created"
	MsgQueueCreated        = "Queue '%s' created\n"
	MsgConsumerSucceed     = "Waiting for messages on queue '%s'.\n"
	MsgProducerSucceed     = "Message '%s' published to queue '%s'\n"
	MsgConnectionPoolEmpty = "connection pool is empty"
)
