package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
	pb "dfs/proto"
	"google.golang.org/grpc"
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
	client := pb.NewMasterTrackerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = client.SendHeartbeat(ctx, &pb.HeartbeatRequest{DataKeeperId: id})
	if err != nil {
		log.Printf("Heartbeat error: %v", err)
		} else {
			log.Println("Heartbeat sent successfully")
		}
	defer conn.Close()
	defer cancel()
}

func handleClient(conn net.Conn) {
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
		HandleFileUpload(reader,conn)
	} else if requestType == "DOWNLOAD" {
		HandleFileDownload(reader,conn)
	} else {
		log.Println("Invalid request type:", requestType)
	}
	defer conn.Close()

}

// HandleFileUpload processes file uploads from the client
func HandleFileUpload(reader *bufio.Reader, conn net.Conn) {

	filename, err := readFilename(reader)
	if err != nil {
		log.Println(err)
		return
	}

	filePath := "storage/" + filename
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	if _, err = io.Copy(file, conn); err != nil {
		log.Printf("Failed to save file %s: %v\n", filePath, err)
		return
	}

	log.Printf("File %s received and saved successfully!\n", filename)
	defer conn.Close()

}

// HandleFileDownload processes file download requests
func HandleFileDownload(reader *bufio.Reader, conn net.Conn) {

	filename, err := readFilename(reader)
	if err != nil {
		log.Println(err)
		return
	}

	filePath := "storage/" + filename
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("File not found: %s\n", filename)
		sendResponse(conn, "ERROR: File not found")
		return
	}

	log.Printf("Sending file: %s\n", filename)

	// Notify client before sending file data
	if err := sendResponse(conn, "OK"); err != nil {
		log.Println(err)
		return
	}

	// Send file data
	if _, err = io.Copy(conn, file); err != nil {
		log.Printf("Error sending file %s: %v\n", filename, err)
	} else {
		log.Println("File sent successfully!")
	}

	defer conn.Close()
	defer file.Close()
}

func (s *DataKeeperServer) ReplicateFile(ctx context.Context, req *pb.ReplicationRequest) (*pb.ReplicationResponse, error) {
	filename := req.Filename
	destinationAddress := req.DestinationAddress

	// Check if the file exists
	filePath := fmt.Sprintf("storage/%s", filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("[REPLICATION ERROR] File %s not found for replication", filename)
		return &pb.ReplicationResponse{Success: false}, err
	}

	// Open connection to destination Data Keeper
	conn, err := net.Dial("tcp", destinationAddress)
	if err != nil {
		log.Printf("[REPLICATION ERROR] Cannot connect to destination %s: %v", destinationAddress, err)
		return &pb.ReplicationResponse{Success: false}, err
	}
	defer conn.Close()

	// Notify destination
	writer := bufio.NewWriter(conn)
	writer.WriteString("UPLOAD\n")
	writer.WriteString(filename + "\n")
	writer.Flush()

	// Send file
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("[REPLICATION ERROR] Cannot open file %s: %v", filename, err)
		return &pb.ReplicationResponse{Success: false}, err
	}
	defer file.Close()

	if _, err := io.Copy(conn, file); err != nil {
		log.Printf("[REPLICATION ERROR] Failed to send file: %v", err)
		return &pb.ReplicationResponse{Success: false}, err
	}

	log.Printf("[REPLICATION SUCCESS] File %s replicated to %s", filename, destinationAddress)
	return &pb.ReplicationResponse{Success: true}, nil
}


// readFilename reads and trims the filename from the connection
func readFilename(reader *bufio.Reader) (string, error) {
	filename, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read filename: %w", err)
	}
	return strings.TrimSpace(filename), nil
}

// sendResponse writes a message to the connection and flushes it
func sendResponse(conn net.Conn, message string) error {
	writer := bufio.NewWriter(conn)
	_, err := writer.WriteString(message + "\n")
	if err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}
	return writer.Flush()
}
