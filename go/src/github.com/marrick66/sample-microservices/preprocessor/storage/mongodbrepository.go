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

//MongoDbRepository is a terrible name for this, will
//probably refactor this later.
type MongoDbRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

//NewMongoDbRepository creates a MongoDb client based on the passed in connection string.
func NewMongoDbRepository(connection string) (*MongoDbRepository, error) {
	client, err := mongo.NewClient(
		options.Client().ApplyURI(connection))

	if err != nil {
		return nil, err
	}

	return &MongoDbRepository{
		client: client}, nil
}

//Connect attempts to setup the background context and connect to the server.
func (repo *MongoDbRepository) Connect() error {
	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := repo.client.Connect(context); err != nil {
		return err
	}

	repo.collection = repo.client.Database("service").Collection("jobs")
	return nil
}

//Get attempts to retrieve the job document by the Id field.
func (repo *MongoDbRepository) Get(id string) (*data.JobRegistration, error) {
	var result data.JobRegistration
	docid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", docid}}
	if err := repo.collection.FindOne(context.Background(), filter).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

//Set will attempt to update a document if it exists, otherwise it inserts a new one.
func (repo *MongoDbRepository) Set(registration *data.JobRegistration) (string, error) {
	var err error
	var id string

	if registration.ID == primitive.NilObjectID {
		var result *mongo.InsertOneResult
		result, err = repo.collection.InsertOne(context.Background(), registration)

		if result != nil {
			id = result.InsertedID.(primitive.ObjectID).Hex()
		}
	} else {
		if err != nil {
			return "", err
		}

		filter := bson.D{{"_id", registration.ID}}
		_, err = repo.collection.ReplaceOne(context.Background(), filter, registration, nil)
		id = registration.ID.Hex()
	}

	if err != nil {
		return "", err
	}

	return id, nil
}

//IsUp is a heartbeat method to make sure we're connected to the database prior to any
//operations.
func (repo *MongoDbRepository) IsUp() bool {
	if err := repo.client.Ping(context.Background(), nil); err != nil {
		return false
	}

	return true
}
