package events

import (
	"github.com/marrick66/sample-microservices/preprocessor/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//JobStatusChanged is received from the event bus to notify
//the service that the local registration should be updated.
type JobStatusChanged struct {
	ID     primitive.ObjectID
	Status data.JobStatus
}
