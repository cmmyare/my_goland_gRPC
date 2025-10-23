package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cmmyare24/go-gRPC/invoicer"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9192", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := invoicer.NewInvoicerClient(conn)

	stream, err := client.StreamInvoices(context.Background())
	if err != nil {
		log.Fatalf("could not open stream: %v", err)
	}

	invoices := []*invoicer.CreateRequest{
		{From: "Ali", To: "Omar", Description: "Invoice 1", Amount: &invoicer.Amount{Amount: 100, Currence: "USD"}},
		{From: "Ali", To: "Hassan", Description: "Invoice 2", Amount: &invoicer.Amount{Amount: 200, Currence: "USD"}},
		{From: "Ali", To: "Amina", Description: "Invoice 3", Amount: &invoicer.Amount{Amount: 300, Currence: "USD"}},
		{From: "Farah", To: "Amina", Description: "Invoice 4", Amount: &invoicer.Amount{Amount: 400, Currence: "USD"}},
	}

	for _, inv := range invoices {
		fmt.Printf("Sending invoice: %+v\n", inv)
		if err := stream.Send(inv); err != nil {
			log.Fatalf("send error: %v", err)
		}
		time.Sleep(time.Millisecond * 500)
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("CloseAndRecv error: %v", err)
	}
	fmt.Printf("Server reply: %v\n", reply)
}
