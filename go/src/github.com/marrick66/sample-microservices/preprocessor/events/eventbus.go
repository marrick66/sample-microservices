package events

//EventBus is the interface representing publishing and subscribing events.
type EventBus interface {
	Publish(topic string, event interface{}) error
	Subscribe(topic string, handler EventHandler) error
}
