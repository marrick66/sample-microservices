package main

import (
	"os"

	"github.com/marrick66/sample-microservices/preprocessor/events"
)

/*
This is a small prototype of an asynchronous job scheduler registration microservice, in which
a dummy job is registered and stored locally for later execution. On successful registration,
an integration event with the details is sent across a service bus for some endpoint to execute
asynchronously.  Eventually, another integration event will be received with the status of the job,
which is updated locally for querying. This is purely for educational purposes and uses a mix of
enterprise and cloud components. I'm not using any frameworks unless absolutely necessary, for
simplicity.
*/

func main() {
	srv, err := NewJobRegistrationServer(":8001")

	if err == nil {
		handler, err := events.NewJobStatusEventHandler(srv.repo)

		if err == nil {
			err = srv.bus.Subscribe(os.Getenv("STATUS_TOPIC"), handler)
		}
	}
	err = srv.Start()
	if err != nil {
		panic("Unable to start RPC server.")
	}

}
