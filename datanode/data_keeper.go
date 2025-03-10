package main

import (
	"bufio"
	"context"
	pb "dfs/proto"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
type DataKeeperServer struct {
	pb.UnimplementedDataKeeperServer
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <arg1> ")
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

// create the gRPC connection and start sending heartbeat
func startHeartbeat(id string) {
	// create connection
	conn, err := grpc.Dial("localhost:50050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Could not connect to Master Tracker: %v", err)
		return
	}

	// free resources
	defer conn.Close()

	// create gRPC client
	client := pb.NewMasterTrackerClient(conn)

	// send heartbeat every second till process stops
	for {
		sendHeartbeat(client, id)
		time.Sleep(5 * time.Second)
	}
}

// notify master that this data node is alive
func sendHeartbeat(client pb.MasterTrackerClient, id string) {
	// context with 1 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	// free resources
	defer cancel()

	// send heartbeat message
	_, err := client.SendHeartbeat(ctx, &pb.HeartbeatRequest{DataKeeperId: id})
	if err != nil {
		log.Printf("Heartbeat error: %v", err)
	} else {
		log.Println("Heartbeat sent successfully")
	}
}

func handleClient(conn net.Conn) {
	reader := bufio.NewReader(conn)

	// read request type (Upload or Download)
	requestType, err := reader.ReadString('\n')
	log.Printf("Request type: %s", requestType)
	if err != nil {
		log.Println("Failed to read request type:", err)
		return
	}
	requestType = strings.TrimSpace(requestType)
	log.Println("Request type:", requestType)

	// handle each request type
	if requestType == "UPLOAD" {
		HandleFileUpload(reader, conn)
	} else if requestType == "DOWNLOAD" {
		HandleFileDownload(reader, conn)
	} else {
		log.Println("Invalid request type:", requestType)
	}
	defer conn.Close() // free resources

}

// HandleFileUpload processes file uploads from the client
func HandleFileUpload(reader *bufio.Reader, conn net.Conn) {
	// get filename from request
	filename, err := readFilename(reader)
	if err != nil {
		log.Println(err)
		return
	}

	// create storage if doesn't exist
	storageDir := "storage"
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		log.Printf("Failed to create storage directory: %v\n", err)
		return
	}

	// create file
	filePath := storageDir + "/" + filename
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file %s: %v\n", filePath, err)
		return
	}
	defer file.Close() // free resources

	// copy data from connection to the file
	if _, err = io.Copy(file, conn); err != nil {
		log.Printf("Failed to save file %s: %v\n", filePath, err)
		return
	}
	defer conn.Close() // free resources

	log.Printf("File %s received and saved successfully!\n", filename)
}

// HandleFileDownload processes file download requests
func HandleFileDownload(reader *bufio.Reader, conn net.Conn) {
	// get file name from request
	filename, err := readFilename(reader)
	if err != nil {
		log.Println(err)
		return
	}

	// get file
	filePath := "storage/" + filename
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("File not found: %s\n", filename)
		sendResponse(conn, "ERROR: File not found")
		return
	}

	log.Printf("Sending file: %s\n", filename)

	// notify client before sending file data
	if err := sendResponse(conn, "OK"); err != nil {
		log.Println(err)
		return
	}

	// send file
	if _, err = io.Copy(conn, file); err != nil {
		log.Printf("Error sending file %s: %v\n", filename, err)
	} else {
		log.Println("File sent successfully!")
	}

	// free resources
	defer conn.Close()
	defer file.Close()
}

// utility function that reads and trims the filename from the connection
func readFilename(reader *bufio.Reader) (string, error) {
	filename, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read filename: %w", err)
	}
	return strings.TrimSpace(filename), nil
}

// utility function that writes a message to the connection and flushes it
func sendResponse(conn net.Conn, message string) error {
	writer := bufio.NewWriter(conn)
	_, err := writer.WriteString(message + "\n")
	if err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}
	return writer.Flush()
}
