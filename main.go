package main

import (
	"log"
	"net"

	// "github.com/cmmyare24/go-gRPC/handlers"
	// "github.com/cmmyare24/go-gRPC/invoicer"
	"github.com/cmmyare24/go-gRPC/models"
	"google.golang.org/grpc"
	"github.com/cmmyare24/go-gRPC/server"
)



func main() {
	models.InitMongoDB()

	list, err := net.Listen("tcp", ":9192")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// invoicer.RegisterInvoicerServer(grpcServer, &handlers.InvoicerServer{})
	server.RegisterServices(grpcServer)

	log.Printf("server listening at %v", list.Addr())
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
