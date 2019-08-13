package messages

//JobRegistrationRequest is a dummy request to start an asynchronous job handled
//outside the service.
type JobRegistrationRequest struct {
	Name string
}

//JobRegistration is the persisted local data that represents
//a job.
type JobRegistration struct {
	Name   string
	Status JobStatus
}
