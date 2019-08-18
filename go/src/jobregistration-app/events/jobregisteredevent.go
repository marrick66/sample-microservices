package events

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//JobRegisteredEvent is sent to the event bus when a successful 
//registration occurs.
type JobRegisteredEvent struct {
	ID   primitive.ObjectID
	Name string
}
