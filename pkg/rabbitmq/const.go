package rabbitmq

const (
	AmqpUrl = "amqp://%s:%s@%s:%s/"

	ErrConnFailed         = "Failed to connect to RabbitMQ : %v"
	ErrChanFailed         = "Failed to open a channel : %v"
	ErrQueueFailed        = "Failed to declare a queue : %v"
	ErrConsumerFailed     = "Failed to register a consumer : %v"
	ErrProducerFailed     = "Failed to publish a message : %v"
	ErrClosingConnection  = "Error closing RabbitMQ connection : %v"
	ErrClosingChannel     = "Error closing RabbitMQ channel : %v"
	ErrGetChannelFromPool = "Error getting channel from pool : %v"

	MsgConnSucceed         = "Connected to RabbitMQ"
	MsgChanCreated         = "Channel created"
	MsgQueueCreated        = "Queue '%s' created\n"
	MsgConsumerSucceed     = "Waiting for messages on queue '%s'.\n"
	MsgConsumerStopped     = "Received stop signal. Stopping consumer..."
	MsgProducerSucceed     = "Message '%s' published to queue '%s'\n"
	MsgConnectionPoolEmpty = "connection pool is empty"
)
