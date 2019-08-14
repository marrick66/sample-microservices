package storage

import "github.com/marrick66/sample-microservices/preprocessor/data"

//JobRegistrationStore is a generic interface for retrieving
//the struct from a distributed key/value store.
type JobRegistrationStore interface {
	Connect() error
	Disconnect() error
	Get(id string) (*data.JobRegistration, error)
	Set(registration *data.JobRegistration) (string, error)
	IsUp() bool
}
