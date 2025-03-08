package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "dfs/proto"
)



func main() {
	// Start TCP Server for receiving files
	go startTCPServer(":50051")

	// Start gRPC heartbeat mechanism
	go startHeartbeat()

	select {} // Keep the process running
}

// startTCPServer listens for file uploads
func startTCPServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
	log.Println("Data Keeper is listening for file uploads on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go handleFileUpload(conn)
	}
}

// handleFileUpload receives and saves a file
func handleFileUpload(conn net.Conn) {
	defer conn.Close()

	// Read filename from client
	reader := bufio.NewReader(conn)
	filename, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Failed to read filename:", err)
		return
	}
	filename = filename[:len(filename)-1] // Remove newline character

	// Create the file
	file, err := os.Create("storage/" + filename)
	if err != nil {
		log.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	// Copy file data
	_, err = io.Copy(file, conn)
	if err != nil {
		log.Println("Failed to save file:", err)
		return
	}

	fmt.Printf("File %s received and saved successfully!\n", filename)
}

// startHeartbeat sends a periodic heartbeat to the Master Tracker
func startHeartbeat() {
	for {
		sendHeartbeat("DataKeeper-1")
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
