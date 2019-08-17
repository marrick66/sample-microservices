package events

//EventHandler is a generic interface that takes
//some event and does whatever's required of it.
type EventHandler interface {
	Handle(event interface{}) error
}
