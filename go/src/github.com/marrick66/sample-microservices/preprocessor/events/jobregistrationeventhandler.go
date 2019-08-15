package events

import (
	"encoding/json"
	"reflect"

	"github.com/streadway/amqp"
)

//JobRegistrationEventHandler is the designated handler for JobRegistered events.
type JobRegistrationEventHandler struct {
	connectionString string
	contentType      string
}

//NewJobRegistrationEventHandler creates a new instance of the handler.
func NewJobRegistrationEventHandler(connectionString string) (*JobRegistrationEventHandler, error) {
	return &JobRegistrationEventHandler{
		connectionString: connectionString,
		contentType:      "application/json"}, nil
}

//ForType returns the reflection type of the event that the handler
//is responsible for.
func (handler *JobRegistrationEventHandler) ForType() reflect.Type {
	return reflect.TypeOf(JobRegisteredEvent{})
}

//Handle takes the event passed and attempts to publish it to the event queue with
//the event specific topic.
func (handler *JobRegistrationEventHandler) Handle(event interface{}) error {
	var err error

	conn, err := amqp.Dial(handler.connectionString)

	if err == nil {
		defer conn.Close()

		channel, err := conn.Channel()

		if err == nil {
			defer channel.Close()
			body, err := json.Marshal(event)

			if err == nil {
				err = channel.Publish(
					"jobevents",
					"jobevents.jobregistered",
					false,
					false,
					amqp.Publishing{
						ContentType: handler.contentType,
						Body:        body})
			}

			return err
		}
	}

	return err
}
