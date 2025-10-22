package server

import (
	"github.com/cmmyare24/go-gRPC/handlers"
	"github.com/cmmyare24/go-gRPC/invoicer"
	// "github.com/cmmyare24/go-gRPC/userpb"
	"google.golang.org/grpc"
)

func RegisterServices(grpcServer *grpc.Server) {
	invoicer.RegisterInvoicerServer(grpcServer, &handlers.InvoicerServer{})
	// userpb.RegisterUserServer(grpcServer, &handlers.UserServer{})
}