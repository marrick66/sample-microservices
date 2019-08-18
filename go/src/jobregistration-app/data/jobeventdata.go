package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//JobEventData is stored to keep track of how many registration events for this job
//have been sent and how long ago the last one was.  We use this to resend events for those
//jobs that have been sitting in "Registered" status for a while.
type JobEventData struct {
	ID                        primitive.ObjectID `bson:"_id,omitempty"` //Map this to the default MongoDb document ID.
	RegistrationEventsSent    int
	LastRegistrationEventTime time.Time
}
