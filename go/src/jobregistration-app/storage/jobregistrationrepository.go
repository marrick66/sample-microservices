package storage

import (
	"context"
	"time"

	"jobregistration-app/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//JobRegistrationRepository is a terrible name for this, will
//probably refactor this later.
type JobRegistrationRepository struct {
	client                    *mongo.Client
	jobRegistrationCollection *mongo.Collection
	eventDataCollection       *mongo.Collection
	defaultTimeout            time.Duration
}

//NewJobRegistrationRepository creates a MongoDb client based on the passed in connection string.
func NewJobRegistrationRepository(connection string) (*JobRegistrationRepository, error) {

	client, err := mongo.NewClient(
		options.Client().ApplyURI(connection))

	if err != nil {
		return nil, err
	}

	//The default timeout is a bit of a wildcard at this point,
	//need to get vendor SLA for key/value storage and modify this to suit.
	//Configuration based will be preferred if this PoC is expanded.
	return &JobRegistrationRepository{
		client:         client,
		defaultTimeout: 100 * time.Millisecond}, nil
}

//Connect attempts to connect to the server and assign the collections used.
func (repo *JobRegistrationRepository) Connect() error {

	context, cancel := context.WithTimeout(context.Background(), repo.defaultTimeout)
	defer cancel()

	if err := repo.client.Connect(context); err != nil {
		return err
	}

	db := repo.client.Database("service")
	repo.jobRegistrationCollection = db.Collection("jobRegistrations")
	repo.eventDataCollection = db.Collection("jobEventData")
	return nil
}

//Disconnect passes the disconnect call through to the client.
func (repo *JobRegistrationRepository) Disconnect() error {
	return repo.client.Disconnect(context.Background())
}

//Get attempts to retrieve the job document by the Id field.
func (repo *JobRegistrationRepository) Get(id string) (*data.JobRegistration, error) {

	var result data.JobRegistration
	var docid primitive.ObjectID
	var err error

	if docid, err = primitive.ObjectIDFromHex(id); err != nil {
		return nil, err
	}

	filter := bson.M{"_id": docid}
	ctx, cancel := context.WithTimeout(context.Background(), repo.defaultTimeout)
	defer cancel()

	if err := repo.jobRegistrationCollection.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

//Set will attempt to update a document if it exists, otherwise it inserts a new one.
func (repo *JobRegistrationRepository) Set(registration *data.JobRegistration) (string, error) {

	if registration.ID == primitive.NilObjectID {
		return repo.insertNew(registration)
	}

	return repo.updateExisting(registration)
}

//IsUp is a heartbeat method to make sure we're connected to the database prior to any
//operations.
func (repo *JobRegistrationRepository) IsUp() bool {

	if err := repo.client.Ping(context.Background(), nil); err != nil {
		return false
	}

	return true
}

//insertNew is a helper method that creates a new transaction and inserts both the
//job registration object and its initial event data.
func (repo *JobRegistrationRepository) insertNew(registration *data.JobRegistration) (string, error) {

	var result *mongo.InsertOneResult
	var session mongo.Session
	var err error

	if session, err = repo.client.StartSession(); err != nil {
		return "", err
	}

	if err = session.StartTransaction(); err != nil {
		return "", err
	}

	//Attempt to insert the job registration, get the generated ID, and use it to
	//insert the job event data.
	context, cancel := context.WithTimeout(context.Background(), repo.defaultTimeout)
	defer cancel()
	defer session.EndSession(context)

	if result, err = repo.jobRegistrationCollection.InsertOne(context, registration); err != nil {
		session.AbortTransaction(context)
		return "", nil
	}

	id := result.InsertedID.(primitive.ObjectID)
	eventData := data.JobEventData{ID: id}

	if _, err = repo.eventDataCollection.InsertOne(context, eventData); err != nil {
		session.AbortTransaction(context)
		return "", err
	}

	session.CommitTransaction(context)

	return id.Hex(), nil
}

//updateExisting is a helper method that wraps all the boilerplate of the client when updating a registration.
func (repo *JobRegistrationRepository) updateExisting(registration *data.JobRegistration) (string, error) {

	filter := bson.M{"_id": registration.ID}
	updateDoc := bson.M{"$set": bson.M{"status": registration.Status}}
	context, cancel := context.WithTimeout(context.Background(), repo.defaultTimeout)
	defer cancel()

	if _, err := repo.jobRegistrationCollection.UpdateOne(context, filter, updateDoc, nil); err != nil {
		return "", err
	}

	return registration.ID.Hex(), nil
}
