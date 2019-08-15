package storage

import (
	"context"
	"time"

	"github.com/marrick66/sample-microservices/preprocessor/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//JobEventDataRepository has some shared data and operations
//with JobRegistrationRepository, should refactor the MongoDB
//specific elements out so they can be shared by both.
type JobEventDataRepository struct {
	client              *mongo.Client
	eventDataCollection *mongo.Collection
	defaultTimeout      time.Duration
}

//NewJobEventDataRepository creates a MongoDb client based on the passed in connection string.
func NewJobEventDataRepository(connection string) (*JobEventDataRepository, error) {

	client, err := mongo.NewClient(
		options.Client().ApplyURI(connection))

	if err != nil {
		return nil, err
	}

	return &JobEventDataRepository{
		client:         client,
		defaultTimeout: 100 * time.Millisecond}, nil
}

//Connect attempts to connect to the server and assign the collections used.
func (repo *JobEventDataRepository) Connect() error {

	context, cancel := context.WithTimeout(context.Background(), repo.defaultTimeout)
	defer cancel()

	if err := repo.client.Connect(context); err != nil {
		return err
	}

	db := repo.client.Database("service")
	repo.eventDataCollection = db.Collection("jobEventData")
	return nil
}

//Get attempts to retrieve the job event data document by the Id field.
func (repo *JobEventDataRepository) Get(id string) (*data.JobEventData, error) {

	var result data.JobEventData
	docid, err := primitive.ObjectIDFromHex(id)

	if err == nil {
		filter := bson.D{{"_id", docid}}
		context, cancel := context.WithTimeout(context.Background(), repo.defaultTimeout)
		defer cancel()

		err = repo.eventDataCollection.FindOne(context, filter).Decode(&result)
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

//Increment adds one to the number of events set and sets the last event
//send time to now.
func (repo *JobEventDataRepository) Increment(id string) error {

	docid, err := primitive.ObjectIDFromHex(id)

	if err == nil {
		filter := bson.D{{"_id", docid}}

		//The $inc is a MongoDB operator to increment a field on a document...
		updateDoc := bson.D{
			{"_id", docid},
			{"$inc", bson.D{{"registrationeventsset", 1}}},
			{"lastregistrationeventsent", time.Now()}}

		context, cancel := context.WithTimeout(context.Background(), repo.defaultTimeout)
		defer cancel()

		_, err = repo.eventDataCollection.UpdateOne(context, filter, updateDoc)
	}

	if err != nil {
		return err
	}

	return nil
}
