package data

//JobStatus is the current state that the job is in.
type JobStatus int

//The possible values of JobStatus.
const (
	//Set this to match the protobuf spec:
	Registered JobStatus = iota + 1
	Running
	Failed
	Completed
)

