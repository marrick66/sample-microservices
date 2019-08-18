package events

import (
	"jobregistration-app/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//JobStatusChanged is received from the event bus to notify
//the service that the local registration should be updated.
type JobStatusChanged struct {
	ID     primitive.ObjectID
	Status data.JobStatus
}
