package messages

//JobRegistration is the simple persisted local data that represents
//a job.
type JobRegistration struct {
	Name   string
	Status JobStatus
}
