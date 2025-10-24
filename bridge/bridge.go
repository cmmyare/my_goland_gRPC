package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/cmmyare24/go-gRPC/invoicer"
	"google.golang.org/grpc"
)

var (
	stream invoicer.Invoicer_ChatStreamClient
	mu     sync.Mutex
)

func main() {
	// 1️⃣ Connect to the gRPC server
	conn, err := grpc.Dial("localhost:9192", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC: %v", err)
	}
	defer conn.Close()

	client := invoicer.NewInvoicerClient(conn)
	stream, err = client.ChatStream(context.Background())
	if err != nil {
		log.Fatalf("Failed to start ChatStream: %v", err)
	}

	// 2️⃣ Start goroutine to continuously receive server messages
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Server closed the stream")
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

	// 3️⃣ Expose HTTP endpoint
	http.HandleFunc("/send", handleSend)
	fmt.Println("🚀 Bridge running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// HTTP handler for POST /send
func handleSend(w http.ResponseWriter, r *http.Request) {
	var req invoicer.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Send to gRPC stream
	if err := stream.Send(&req); err != nil {
		http.Error(w, fmt.Sprintf("Stream send error: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("📤 Sent from API → gRPC: %+v\n", req)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message sent to gRPC stream"))
}
