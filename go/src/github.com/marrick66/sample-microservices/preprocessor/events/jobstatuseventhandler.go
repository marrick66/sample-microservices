package events

import (
	"github.com/marrick66/sample-microservices/preprocessor/data"
	"github.com/marrick66/sample-microservices/preprocessor/storage"
)

type JobStatusEventHandler struct {
	repo storage.JobRegistrationStore
}

func NewJobStatusEventHandler(repository storage.JobRegistrationStore) (*JobStatusEventHandler, error) {
	return &JobStatusEventHandler{
		repo: repository}, nil
}

func (handler *JobStatusEventHandler) DefaultEvent() interface{} {
	return new(JobStatusChanged)
}

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
