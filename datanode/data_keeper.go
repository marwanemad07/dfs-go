package main

import (
	"bufio"
	"context"
	pb "dfs/proto"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataKeeperServer struct {
	pb.UnimplementedDataKeeperServer
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run data_keeper.go <master_address_type> <data_node_port> <master_port>")
		return
	}

	// Start TCP Server for receiving files
	masterAddress := ""
	dataNodePort := os.Args[2]
	masterPort := os.Args[3]

	switch os.Args[1] {
	case "local":
		masterAddress = "localhost"
	case "docker":
		masterAddress = "host.docker.internal"
	default:
		fmt.Println("Invalid mode. Use 'local' or 'docker'.")
		os.Exit(1)
	}

	go startTCPServer(dataNodePort)

	// Start gRPC heartbeat mechanism
	go startHeartbeat(masterAddress, masterPort, dataNodePort)

	tcpPort := os.Args[2]
	tcpPortInt, err := strconv.Atoi(tcpPort)
	if err != nil {
		log.Fatalf("Invalid TCP port: %v", err)
	}
	grpcPort := tcpPortInt + 1000

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDataKeeperServer(grpcServer, &DataKeeperServer{})
	log.Printf("Data Keeper gRPC server started on port %d", grpcPort)
	go grpcServer.Serve(lis)
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
func startHeartbeat(masterAddress string, masterPort string, id string) {
	// create connection
	conn, err := grpc.Dial(masterAddress+":"+masterPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		time.Sleep(1 * time.Second)
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
		HandleFileUpload(reader, conn)
	} else if requestType == "DOWNLOAD" {
		HandleFileDownload(reader, conn)
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
	filePath := filename

	if len(filename) < 30 {
		filePath = filepath.Join("storage", filename)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Printf("Failed to create storage directory: %v\n", err)
			return
		}
	}

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
	log.Printf("[DATA KEEPER] Replicating file %s to %s", req.Filename, req.DestinationAddress)
	filename := req.Filename
	destinationAddress := req.DestinationAddress
	filePath := filename
	log.Printf("[DATA KEEPER] Replicating file %s to %s", filename, destinationAddress)

	// Check if the file exists.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("[REPLICATION ERROR] File %s not found", filename)
		return &pb.ReplicationResponse{Success: false}, err
	}

	// Open connection to destination Data Keeper.
	conn, err := net.Dial("tcp", "localhost:"+destinationAddress)
	if err != nil {
		log.Printf("[REPLICATION ERROR] Cannot connect to destination %s: %v", destinationAddress, err)
		return &pb.ReplicationResponse{Success: false}, err
	}
	defer conn.Close()

	// Notify the destination about the upcoming file upload.
	writer := bufio.NewWriter(conn)
	writer.WriteString("UPLOAD\n")
	writer.WriteString(filename + strconv.Itoa(rand.Intn(1000)) + "\n")
	writer.Flush()

	// Send the file.
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

func (s *DataKeeperServer) StartReplication(ctx context.Context, req *pb.ReplicationRequest) (*pb.ReplicationResponse, error) {
	log.Printf("[Replication] Starting replication of %s to %s", req.Filename, req.DestinationAddress)

	// Send the file from this Data Keeper to the destination
	err := sendFile(req.DestinationAddress, req.Filename)
	if err != nil {
		log.Printf("[Replication] File transfer failed: %v", err)
		return &pb.ReplicationResponse{Success: false}, err
	}

	log.Printf("[Replication] File %s replicated successfully to %s", req.Filename, req.DestinationAddress)
	return &pb.ReplicationResponse{Success: true}, nil
}

func sendFile(destination, filename string) error {
	srcPath := "storage/" + filename
	dstAddr := destination + ":6000" // Destination listens on port 6000

	// Open the source file
	file, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// Connect to the destination Data Keeper
	conn, err := net.Dial("tcp", dstAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to destination: %w", err)
	}
	defer conn.Close()

	// Send the filename first
	_, err = conn.Write([]byte(filename + "\n"))
	if err != nil {
		return fmt.Errorf("failed to send filename: %w", err)
	}

	// Copy file data over TCP
	_, err = io.Copy(conn, file)
	if err != nil {
		return fmt.Errorf("file transfer failed: %w", err)
	}

	log.Printf("[Replication] File %s sent successfully to %s", filename, destination)
	return nil
}
