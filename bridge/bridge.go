
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	invoicer "github.com/cmmyare24/go-gRPC/invoicer"

)


var (
	stream invoicer.Invoicer_ChatStreamClient
	streamMu sync.Mutex
	wsClients = make(map[*websocket.Conn]bool)
	wsMu      sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	conn, err := grpc.Dial("localhost:9192", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := invoicer.NewInvoicerClient(conn)

	stream, err = client.ChatStream(context.Background())
	if err != nil {
		log.Fatalf("Error creating gRPC stream: %v", err)
	}

	// Listen for messages coming from the gRPC server
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				fmt.Println("⚠️ Stream closed:", err)
				break
			}
			msg := fmt.Sprintf("💬 %s → %s: %s (%d %s)",
				"Server",
				resp.To,
				resp.Description,
				resp.Amount.Amount,
				resp.Amount.Currence,
			)
			fmt.Println(msg)

			// Broadcast message to all connected WebSocket clients
			broadcastToWS(msg)
		}
	}()

	http.HandleFunc("/send", handleSend)
	http.HandleFunc("/listen", handleListen)

	fmt.Println("🚀 Bridge running on:")
	fmt.Println("- REST: http://localhost:8080/send")
	fmt.Println("- WS:   ws://localhost:8080/listen")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSend(w http.ResponseWriter, r *http.Request) {
    // Decode JSON
    var req struct {
        From, To, Description string
        Amount invoicer.Amount
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Connect gRPC
    conn, err := grpc.Dial("localhost:9192", grpc.WithInsecure())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer conn.Close()

    client := invoicer.NewInvoicerClient(conn)
    stream, err := client.ChatStream(context.Background())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    msg := &invoicer.CreateRequest{
        From: req.From,
        To: req.To,
        Description: req.Description,
        Amount: &invoicer.Amount{
            Amount: req.Amount.Amount,
            Currence: req.Amount.Currence,
        },
    }

    // Send and receive response immediately
    if err := stream.Send(msg); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    resp, err := stream.Recv()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return response to HTTP client
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "✅ Message sent and response received",
        "response": resp,
    })
}





func handleListen(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	wsMu.Lock()
	wsClients[ws] = true
	wsMu.Unlock()

	fmt.Println("🌐 New WebSocket client connected")

	for {
		// Keep connection alive
		time.Sleep(time.Minute)
	}
}

func broadcastToWS(message string) {
	wsMu.Lock()
	defer wsMu.Unlock()

	for conn := range wsClients {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("⚠️ Removing disconnected WebSocket client:", err)
			conn.Close()
			delete(wsClients, conn)
		}
	}
}






























// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"sync"

// 	"github.com/cmmyare24/go-gRPC/invoicer"
// 	"google.golang.org/grpc"
// )

// var (
// 	stream invoicer.Invoicer_ChatStreamClient
// 	mu     sync.Mutex
// )

// func main() {
// 	// 1️⃣ Connect to the gRPC server
// 	conn, err := grpc.Dial("localhost:9192", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("Failed to connect to gRPC: %v", err)
// 	}
// 	defer conn.Close()

// 	client := invoicer.NewInvoicerClient(conn)
// 	stream, err = client.ChatStream(context.Background())
// 	if err != nil {
// 		log.Fatalf("Failed to start ChatStream: %v", err)
// 	}

// 	// 2️⃣ Start goroutine to continuously receive server messages
// 	go func() {
// 		for {
// 			resp, err := stream.Recv()
// 			if err == io.EOF {
// 				fmt.Println("Server closed the stream")
// 				return
// 			}
// 			if err != nil {
// 				fmt.Printf("Receive error: %v\n", err)
// 				return
// 			}
// 			fmt.Printf("💬 Server → %s: %s (%d %s)\n",
// 				resp.To, resp.Description, resp.Amount.Amount, resp.Amount.Currence)
// 		}
// 	}()

// 	// 3️⃣ Expose HTTP endpoint
// 	http.HandleFunc("/send", handleSend)
// 	fmt.Println("🚀 Bridge running on http://localhost:8080")
// 	http.ListenAndServe(":8080", nil)
// }

// // HTTP handler for POST /send
// func handleSend(w http.ResponseWriter, r *http.Request) {
// 	var req invoicer.CreateRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	mu.Lock()
// 	defer mu.Unlock()

// 	// Send to gRPC stream
// 	if err := stream.Send(&req); err != nil {
// 		http.Error(w, fmt.Sprintf("Stream send error: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Printf("📤 Sent from API → gRPC: %+v\n", req)
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Message sent to gRPC stream"))
// }
