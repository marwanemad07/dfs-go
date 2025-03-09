package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	pb "dfs/proto"
	"dfs/utils"
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <arg1> <arg2>")
		return
	}

	// Start TCP Server for receiving files
	port := os.Args[1]
	go startTCPServer(port)

	// Start gRPC heartbeat mechanism
	go startHeartbeat(port)

	select {} // Keep the process running
}

func startTCPServer(port string) {
	address := ":" + port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
	log.Println("Data Keeper is listening on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}

		go handleClient(conn)
	}
}

func startHeartbeat(id string) {
	for {
		sendHeartbeat(id)
		time.Sleep(1 * time.Second)
	}
}

// sendHeartbeat notifies Master Tracker that this Data Keeper is alive
func sendHeartbeat(id string) {
	conn, err := grpc.Dial("localhost:50050", grpc.WithInsecure())
	if err != nil {
		log.Printf("Could not connect to Master Tracker: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewMasterTrackerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = client.SendHeartbeat(ctx, &pb.HeartbeatRequest{DataKeeperId: id})
	if err != nil {
		log.Printf("Heartbeat error: %v", err)
	} else {
		log.Println("Heartbeat sent successfully")
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Read request type (Upload or Download)
	requestType, err := reader.ReadString('\n')
	log.Printf("Request type: %s", requestType)
	if err != nil {
		log.Println("Failed to read request type:", err)
		return
	}
	requestType = strings.TrimSpace(requestType)
	log.Println("Request type:", requestType)

	// Check if it's an upload or download request
	if requestType == "UPLOAD" {
		utils.HandleFileUpload(reader,conn)
	} else if requestType == "DOWNLOAD" {
		utils.HandleFileDownload(reader,conn)
	} else {
		log.Println("Invalid request type:", requestType)
	}
}