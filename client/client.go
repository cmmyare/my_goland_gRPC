package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cmmyare24/go-gRPC/invoicer"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9192", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer conn.Close()

	client := invoicer.NewInvoicerClient(conn)
	stream, err := client.ChatStream(context.Background())
	if err != nil {
		log.Fatalf("❌ Could not open stream: %v", err)
	}

	// Goroutine to listen to server responses
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Server closed stream.")
				return
			}
			if err != nil {
				fmt.Printf("Receive error: %v\n", err)
				return
			}
			fmt.Printf("💬 Server → %s: %s (%d %s)\n",
				resp.To, resp.Description, resp.Amount.Amount, resp.Amount.Currence)
		}
	}()

	// Read user input and send messages
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type invoices in format: to description amount (or 'exit' to quit)")

	for {
		fmt.Print("You → ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "exit" {
			stream.CloseSend()
			break
		}

		parts := strings.SplitN(text, " ", 3)
		if len(parts) < 3 {
			fmt.Println("Usage: to description amount")
			continue
		}

		to := parts[0]
		desc := parts[1]
		var amt int64
		fmt.Sscan(parts[2], &amt)

		req := &invoicer.CreateRequest{
			From:        "Ali",
			To:          to,
			Description: desc,
			Amount:      &invoicer.Amount{Amount: amt, Currence: "USD"},
		}

		if err := stream.Send(req); err != nil {
			fmt.Printf("Send error: %v\n", err)
			break
		}
	}
}



// static bidirectional streaming
// func main() {
// 	conn, err := grpc.Dial("localhost:9192", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()

// 	client := invoicer.NewInvoicerClient(conn)

// 	// Start bidirectional streaming
// 	stream, err := client.ChatStream(context.Background())
// 	if err != nil {
// 		log.Fatalf("could not open stream: %v", err)
// 	}

// 	// Use goroutine to send messages to the server
// 	go func() {
// 		for i := 1; i <= 5; i++ {
// 			req := &invoicer.CreateRequest{
// 				From:        "Ali",
// 				To:          fmt.Sprintf("User%d", i),
// 				Description: fmt.Sprintf("Invoice %d", i),
// 				Amount:      &invoicer.Amount{Amount: int64(i * 100), Currence: "USD"},
// 			}
// 			fmt.Printf("Client sending: %+v\n", req)
// 			if err := stream.Send(req); err != nil {
// 				log.Fatalf("Send error: %v", err)
// 			}
// 			time.Sleep(time.Second)
// 		}
// 		stream.CloseSend()
// 	}()

// 	// Receive stream responses from the server
// 	for {
// 		resp, err := stream.Recv()
// 		if err == io.EOF {
// 			fmt.Println("Server closed stream")
// 			break
// 		}
// 		if err != nil {
// 			log.Fatalf("Receive error: %v", err)
// 		}
// 		fmt.Printf("Client received: %+v\n", resp)
// 	}
// }





// func main() {
// 	conn, err := grpc.Dial("localhost:9192", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()

// 	client := invoicer.NewInvoicerClient(conn)

// 	stream, err := client.StreamInvoices(context.Background())
// 	if err != nil {
// 		log.Fatalf("could not open stream: %v", err)
// 	}

// 	invoices := []*invoicer.CreateRequest{
// 		{From: "Ali", To: "Omar", Description: "Invoice 1", Amount: &invoicer.Amount{Amount: 100, Currence: "USD"}},
// 		{From: "Ali", To: "Hassan", Description: "Invoice 2", Amount: &invoicer.Amount{Amount: 200, Currence: "USD"}},
// 		{From: "Ali", To: "Amina", Description: "Invoice 3", Amount: &invoicer.Amount{Amount: 300, Currence: "USD"}},
// 		{From: "Farah", To: "Amina", Description: "Invoice 4", Amount: &invoicer.Amount{Amount: 400, Currence: "USD"}},
// 	}

// 	for _, inv := range invoices {
// 		fmt.Printf("Sending invoice: %+v\n", inv)
// 		if err := stream.Send(inv); err != nil {
// 			log.Fatalf("send error: %v", err)
// 		}
// 		time.Sleep(time.Millisecond * 500)
// 	}

// 	reply, err := stream.CloseAndRecv()
// 	if err != nil {
// 		log.Fatalf("CloseAndRecv error: %v", err)
// 	}
// 	fmt.Printf("Server reply: %v\n", reply)
// }
