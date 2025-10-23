package handlers

import (
	"context"
	"fmt"
	"io"
	"time"

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

func (s *InvoicerServer) StreamInvoices(stream invoicer.Invoicer_StreamInvoicesServer) error {
	fmt.Println("Client streaming started...")
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&invoicer.StreamResponse{
				Message: fmt.Sprintf("Received %d invoices from client", count),
			})
		}
		if err != nil {
			return err
		}

		fmt.Printf("Received invoice: From=%s To=%s Amount=%d\n", req.From, req.To, req.Amount.Amount)
		count++
	}
}

func (s *InvoicerServer) GetInvoiceStream(req *invoicer.CreateRequest, stream invoicer.Invoicer_GetInvoiceStreamServer) error {
	for i := 1; i <= 5; i++ {
		resp := &invoicer.CreateResponse{
			From:        req.From,
			To:          req.To,
			Description: fmt.Sprintf("Invoice stream item %d", i),
			Amount:      req.Amount,
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

