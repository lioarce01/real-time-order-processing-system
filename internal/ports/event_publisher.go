package ports

type EventPublisher interface {
	Publish(topic string, message interface{}) error
}
