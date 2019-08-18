package storage

import "jobregistration-app/data"

//JobRegistrationStore is a generic interface for reading
//and writing the struct from/to a distributed key/value store.
type JobRegistrationStore interface {
	Connect() error
	Disconnect() error
	Get(id string) (*data.JobRegistration, error)
	Set(registration *data.JobRegistration) (string, error)
	IsUp() bool
}

//JobEventDataStore is a generic interface for operations to
//access the struct from a distributed key/value store.
type JobEventDataStore interface {
	Connect() error
	Disconnect() error
	Get(id string) (*data.JobEventData, error)
	Increment(id string) error
	IsUp() bool
}
