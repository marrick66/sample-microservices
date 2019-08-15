package events

import "go/types"

//The JobEventCoordinator should be run asynchronously at startup to
//take events received from some communication channel and map them
//to the appropriate event handler based on type.
type JobEventCoordinator struct {
	EventChannel chan interface{}
	Handlers     map[types.Type]EventHandler
}

//NewJobEventCoordinator takes all of the handlers passed and adds them to
//the map by the type they're for. Since we're not checking for uniqueness between
//handler types, the last one in wins.
func NewJobEventCoordinator(handlers ...EventHandler) JobEventCoordinator {

	coordinator := JobEventCoordinator{
		EventChannel: make(chan interface{}),
		Handlers:     make(map[types.Type]EventHandler)}

	for _, handler := range handlers {
		coordinator.Handlers[handler.ForType()] = handler
	}

	return coordinator
}
