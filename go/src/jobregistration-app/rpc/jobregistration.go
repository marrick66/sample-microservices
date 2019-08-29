package rpc

import (
	"context"
	"net"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"

	"jobregistration-app/data"
	"jobregistration-app/events"
	"jobregistration-app/storage"

	"google.golang.org/grpc"
)

//Instead of using REST/HTTP, I'm using gRPC over TCP for service calls. This is just to mix it up, but in production this wouldn't be exposed to external clients.
//We'd need an API gateway to make it friendly and do things like TLS termination, authentication tokens, request throttling, etc. Load balancing is another element,
//but it looks like it's handled differently between providers. For further authentication/authorization, we'd want to handle mediating the token in this service as well.
//In a future phase, I'd integrate this in with an IDP (Google Auth is supported out of the box).

//JobRegistrationServerImpl is the tcp listener and gRPC server implementation:
type JobRegistrationServerImpl struct {
	listener *net.Listener
	rpcSrv   *grpc.Server
	port     string
	Repo     storage.JobRegistrationStore
	Bus      events.EventBus
	exchange string
	topic    string
}

//Register is the actual implementation of the respective RPC call to register a job.
func (srv *JobRegistrationServerImpl) Register(ctx context.Context, request *RegistrationRequest) (*RegistrationReply, error) {
	if !srv.Repo.IsUp() {
		if err := srv.Repo.Connect(); err != nil {
			return nil, err
		}
	}

	id, err := srv.Repo.Set(
		&data.JobRegistration{Name: request.Name, Status: data.Registered})

	if err == nil {
		//Asynchronously send the JobRegisteredEvent to the coordinator for handling. It's possible
		//that the channel blocks, so a lot of goroutines could exist here.
		rawID, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			go srv.Bus.Publish(srv.topic, events.JobRegisteredEvent{ID: rawID, Name: request.Name})
		}

		return &RegistrationReply{Id: id}, nil
	}

	return nil, err

}

//GetRegistration is the actual implementation of the respective RPC call to get a registered job.
func (srv *JobRegistrationServerImpl) GetRegistration(ctx context.Context, request *GetRegistrationRequest) (*GetRegistrationReply, error) {
	if !srv.Repo.IsUp() {
		if err := srv.Repo.Connect(); err != nil {
			return nil, err
		}
	}

	//Query mongoDb for the document, if ErrNoDocuments is returned,
	//send the custom not found reply back.
	registration, err := srv.Repo.Get(request.Id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &GetRegistrationReply{Id: "", Status: GetRegistrationReply_NOTFOUND}, nil
		}

		return nil, err
	}

	//Since the type enum for status and the RPC are different, need
	//to map them here. There's probably a better way to do this, need
	//to research protobuf more...
	var status GetRegistrationReply_Status
	switch registration.Status {
	case data.Registered:
		status = GetRegistrationReply_REGISTERED
		break
	case data.Running:
		status = GetRegistrationReply_RUNNING
		break
	case data.Completed:
		status = GetRegistrationReply_COMPLETED
		break
	case data.Failed:
		status = GetRegistrationReply_FAILED
	}

	id := registration.ID.Hex()

	return &GetRegistrationReply{Id: id, Status: status}, nil
}

//NewJobRegistrationServer creates the server object and gRPC dependencies.
func NewJobRegistrationServer(bus *events.EventBus, repository *storage.JobRegistrationStore, port string, exchange string, topic string) (*JobRegistrationServerImpl, error) {
	srv := JobRegistrationServerImpl{
		listener: nil,
		port:     port,
		rpcSrv:   grpc.NewServer(),
		exchange: exchange
		topic:    topic,
		repo:	  repository,
		bus:	  bus}

	RegisterJobRegistrationServer(srv.rpcSrv, &srv)

	return &srv, nil
}

//Start sets up the tcp listener and starts serving requests.
func (srv *JobRegistrationServerImpl) Start() error {
	if srv.listener == nil {
		listener, err := net.Listen("tcp", srv.port)
		if err != nil {
			return err
		}

		srv.listener = &listener
	}

	if err := srv.rpcSrv.Serve(*srv.listener); err != nil {
		return err
	}

	return nil
}
