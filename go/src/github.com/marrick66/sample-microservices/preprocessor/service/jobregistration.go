package main

import (
	"context"
	"net"

	"github.com/marrick66/sample-microservices/preprocessor/messages"
	"google.golang.org/grpc"
)

//JobRegistrationServerImpl is the tcp listener and gRPC server implementation:
type JobRegistrationServerImpl struct {
	listener *net.Listener
	rpcSrv   *grpc.Server
	port     string
}

//Register is the actual implementation of the respective RPC call to register a job.
func (srv *JobRegistrationServerImpl) Register(ctx context.Context, request *messages.RegistrationRequest) (*messages.RegistrationReply, error) {
	return &messages.RegistrationReply{
		Id: 1}, nil
}

//GetRegistration is the actual implementation of the respective RPC call to get a registered job.
func (srv *JobRegistrationServerImpl) GetRegistration(context.Context, *messages.GetRegistrationRequest) (*messages.GetRegistrationReply, error) {
	return &messages.GetRegistrationReply{
		Id:     -1,
		Status: messages.GetRegistrationReply_NOTFOUND}, nil
}

//NewJobRegistrationServer creates the server object and gRPC dependencies.
func NewJobRegistrationServer(port string) *JobRegistrationServerImpl {
	srv := JobRegistrationServerImpl{
		listener: nil,
		port:     port,
		rpcSrv:   grpc.NewServer()}

	messages.RegisterJobRegistrationServer(srv.rpcSrv, &srv)

	return &srv
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
