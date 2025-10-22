package handlers

import (
	"context"
	"fmt"

	"github.com/cmmyare24/go-gRPC/controllers"
	"github.com/cmmyare24/go-gRPC/invoicer"
)

type InvoicerServer struct {
	invoicer.UnimplementedInvoicerServer
}

func (s *InvoicerServer) Create(ctx context.Context, req *invoicer.CreateRequest) (*invoicer.CreateResponse, error) {
	fmt.Printf("req >>>>>>: %+v\n", req)

	resp, err := controllers.CreateInvoice(ctx, req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response: %+v\n", resp)
	return resp, nil
}

func (s *InvoicerServer) Update(ctx context.Context, req *invoicer.UpdateRequest) (*invoicer.UpdateResponse, error) {
	success, err := controllers.UpdateInvoice(ctx, req)
	if err != nil {
		return nil, err
	}

	return &invoicer.UpdateResponse{Success: success}, nil
}
