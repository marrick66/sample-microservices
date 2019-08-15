package events

import "reflect"

//EventHandler is a generic interface that takes
//some event and does whatever's required of it.
type EventHandler interface {
	ForType() reflect.Type
	Handle(event interface{}) error
}
