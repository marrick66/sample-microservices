package events

type EventBus interface {
	Publish(topic string, event interface{}) error
	Subscribe(topic string, handler EventHandler) error
}