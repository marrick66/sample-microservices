package main

import (
	"jobregistration-app/events"
	"jobregistration-app/rpc"
	"os"
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
	var srv *rpc.JobRegistrationServerImpl
	var handler events.EventHandler
	var err error

	//Get the RPC server instance:
	if srv, err = rpc.NewJobRegistrationServer(":8001"); err != nil {
		panic(err)
	}

	//Create the JobStatus event handler:
	if handler, err = events.NewJobStatusEventHandler(srv.Repo); err != nil {
		panic(err)
	}

	//Subscribe the even handler to the topic:
	if err = srv.Bus.Subscribe(os.Getenv("STATUS_TOPIC"), handler); err != nil {
		panic(err)
	}

	//Start the server:
	if err = srv.Start(); err != nil {
		panic(err)
	}

}
