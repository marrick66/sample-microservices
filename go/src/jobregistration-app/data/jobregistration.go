package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//JobRegistration is the simple persisted local data that represents
//a job. Since the RPC interface defines messages that hide the MongoDB
//details, making the ID a BSON object for ease is probably ok.
type JobRegistration struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"` //Map this to the default MongoDb document ID.
	Name   string
	Status JobStatus
}
