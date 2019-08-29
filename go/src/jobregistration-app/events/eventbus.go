package events

//EventBus is the interface representing publishing and subscribing events.
type EventBus interface {
	Publish(exchange string, topic string, event interface{}) error
	Subscribe(exchange string, topic string, handler EventHandler) error
}
