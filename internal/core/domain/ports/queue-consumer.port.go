package ports

type MessageHandler func(messageBody []byte) error

type IQueueConsumer interface {
	ConsumeQueue(queueName string, handler MessageHandler) error
}
