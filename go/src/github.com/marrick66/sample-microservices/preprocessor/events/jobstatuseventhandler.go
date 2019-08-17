package events

import (
	"github.com/marrick66/sample-microservices/preprocessor/data"
	"github.com/marrick66/sample-microservices/preprocessor/storage"
)

//JobStatusEventHandler is the type that implements the EventHandler interface for the
//JobStatus event.
type JobStatusEventHandler struct {
	repo storage.JobRegistrationStore
}

//NewJobStatusEventHandler constructs a default instance.
func NewJobStatusEventHandler(repository storage.JobRegistrationStore) (*JobStatusEventHandler, error) {
	return &JobStatusEventHandler{
		repo: repository}, nil
}

//DefaultEvent returns a new JobStatusChanged event that messages are deserialized to.
func (handler *JobStatusEventHandler) DefaultEvent() interface{} {
	return new(JobStatusChanged)
}

//Handle takes a JobStatusChanged event, and updates the repository with the new status.
func (handler *JobStatusEventHandler) Handle(event interface{}) error {

	statusEvent, ok := event.(*JobStatusChanged)
	if ok {
		_, err := handler.repo.Set(&data.JobRegistration{ID: statusEvent.ID, Status: statusEvent.Status})

		if err != nil {
			return err
		}
	}

	return nil
}
