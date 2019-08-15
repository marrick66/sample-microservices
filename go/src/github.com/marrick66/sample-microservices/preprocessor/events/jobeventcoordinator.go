package events

import (
	"reflect"
)

//The JobEventCoordinator should be run asynchronously at startup to
//take events received from some communication channel and map them
//to the appropriate event handler based on type. Should probably profile
//the reflection methods to see if this has any performance implications.
type JobEventCoordinator struct {
	isRunning    bool
	EventChannel chan interface{}
	Handlers     map[reflect.Type]EventHandler
}

//StartEventLoopEvent is sent to the coordinator's event channel
//to start handling other events. While in that event loop, these
//events are ignored.
type StartEventLoopEvent struct{}

//StopEventLoopEvent is sent to the coordinator to stop the event loop goroutine.
type StopEventLoopEvent struct{}

//EventLoopCompletedEvent tells the coordinator that the event loop has ended,
//and another can be started.
type EventLoopCompletedEvent struct{}

//NewJobEventCoordinator takes all of the handlers passed and adds them to
//the map by the type they're for. Since we're not checking for uniqueness between
//handler types, the last one in wins.
func NewJobEventCoordinator(handlers ...EventHandler) *JobEventCoordinator {

	coordinator := JobEventCoordinator{
		isRunning:    false,
		EventChannel: make(chan interface{}),
		Handlers:     make(map[reflect.Type]EventHandler)}

	for _, handler := range handlers {
		coordinator.Handlers[handler.ForType()] = handler
	}

	go coordinator.initEventLoop()

	return &coordinator
}

//Run sends the start event loop message to the coordinator's channel.
func (crdnr *JobEventCoordinator) Run() {
	crdnr.EventChannel <- StartEventLoopEvent{}
}

//Stop sends the start event loop message to the coordinator's channel.
func (crdnr *JobEventCoordinator) Stop() {
	crdnr.EventChannel <- StopEventLoopEvent{}
}

//initEventLoop should be running when the coordinator is created, or when
//the an existing one is stopped. It waits for a StartEventLoopEvent, starts
//the event loop goroutine, and returns.
func (crdnr *JobEventCoordinator) initEventLoop() {

	for {
		event := <-crdnr.EventChannel

		switch event.(type) {
		case StartEventLoopEvent:
			go crdnr.eventLoop()
			return
		default:
			//ignore any other events
		}
	}

}

//eventLoop waits for an event on the coordinator's channel.
//If it's a StopEventLoopEvent, we start another init goroutine and return.
//Otherwise, attempt to find a matching handler and execute it. Since the handler is
//running synchronously, it's effectively processing one event at a time.  Can augment the
//handlers themselves to have this property and run Handle() as a goroutine to get some concurrency
//between event types.
func (crdnr *JobEventCoordinator) eventLoop() {

	for {
		event := <-crdnr.EventChannel

		switch event.(type) {
		case StopEventLoopEvent:
			go crdnr.initEventLoop()
			return
		default:
			evtType := reflect.TypeOf(event)
			handler, exists := crdnr.Handlers[evtType]
			if exists {
				handler.Handle(event)
			}
		}
	}
}
