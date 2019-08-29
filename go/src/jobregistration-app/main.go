package main

import (
	"jobregistration-app/converters"
	"jobregistration-app/events"
	"jobregistration-app/rpc"
	"jobregistration-app/storage"
	"log"
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
	var repository storage.JobRegistrationStore
	var bus events.EventBus
	var srv *rpc.JobRegistrationServerImpl
	var handler events.EventHandler
	var err error

	defaultBus := os.Getenv("EVENT_BUS")
	defaultExchange := os.Getenv("EXCHANGE")
	statusTopic := os.Getenv("STATUS_TOPIC")
	registeredTopic := os.Getenv("REGISTERED_TOPIC")
	defaultDb := os.Getenv("JOBS_DB")
	converter := &converters.JSONByteConverter{}

	//Get the repository:
	if repository, err = storage.NewJobRegistrationRepository(defaultDb); err != nil {
		log.Fatalf("Failed to get job repository: %v", err)
	}

	//Get the event bus:
	if bus, err = events.NewAMQPEventBus(defaultBus, converter); err != nil {
		log.Fatalf("Failed to get event bus: %v", err)
	}

	//Get the RPC server instance:
	if srv, err = rpc.NewJobRegistrationServer(bus, repository, ":8001", defaultExchange, registeredTopic); err != nil {
		log.Fatalf("Failed to get job registration server: %v", err)
	}

	//Create the JobStatus event handler:
	if handler, err = events.NewJobStatusEventHandler(srv.Repo); err != nil {
		log.Fatalf("Failed to get job status event handler: %v", err)
	}

	//Subscribe the even handler to the status topic:
	if err = srv.Bus.Subscribe(defaultExchange, statusTopic, handler); err != nil {
		log.Fatalf("Failed to get subscribe job status event handler: %v", err)
	}

	//Start the server:
	if err = srv.Start(); err != nil {
		panic(err)
	}

}
