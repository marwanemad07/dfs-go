package main

import (
	"bufio"
	"context"
	"dfs/config"
	pb "dfs/proto"
	"dfs/utils"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataKeeperServer struct {
	pb.UnimplementedDataKeeperServer
	portsTcp  []*pb.PortStatus
	portsGrpc []*pb.PortStatus
}

type PortStatus struct {
	PortNumber  int  `json:"portNumber"`
	IsAvailable bool `json:"isAvailable"`
}

type Globals struct {
	masterAddress string
	nodeName      string
	nodeAddress   string
}

var globals = &Globals{}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <master_address_type> <number_of_ports> <start_tcp_ports>  ...")
		return
	}

	// Start TCP Server for receiving files
	masterAddress := ""
	nodeAddress := ""

	switch os.Args[1] {
	case "local":
		masterAddress = "localhost"
		nodeAddress = "localhost"
	case "docker":
		masterAddress = "host.docker.internal"
		nodeAddress = "host.docker.internal"
	case "network":
		masterAddress, _ = utils.GetLocalIP()
		nodeAddress, _ = utils.GetLocalIP()
	default:
		fmt.Println("Invalid mode. Use 'local' or 'docker'.")
		os.Exit(1)
	}

	numberOfdataNodePorts, _ := strconv.Atoi(os.Args[2])
	startTcpPort, err := strconv.Atoi(os.Args[3])

	if err != nil {
		log.Fatalf("Invalid number of ports: %v", err)
		return
	}

	tcpPorts := make([]*pb.PortStatus, 0)
	counter := 0 // counter for ports
	for range numberOfdataNodePorts {
		port := startTcpPort + counter
		for !isPortAvailable(port) {
			counter += 1
			port = startTcpPort + counter
		}

		tcpPorts = append(tcpPorts, &pb.PortStatus{
			PortNumber:  int32(port),
			IsAvailable: true,
		})
		counter++
	}

	dataKeeperServer := &DataKeeperServer{}

	cfg := config.LoadConfig("config.json")

	masterAddress = masterAddress + ":" + strconv.Itoa(cfg.Server.Port)

	startGrpcPort := startTcpPort + counter
	grpcPorts := make([]*pb.PortStatus, 0)
	for i := range numberOfdataNodePorts {
		log.Printf("Port: %d", tcpPorts[i].PortNumber)
		go startTCPServer(tcpPorts[i])

		// Start gRPC server on the port tcpPort + i
		port := startGrpcPort + counter
		for !isPortAvailable(port) {
			counter += 1
			port = startGrpcPort + counter
		}

		grpcPorts = append(grpcPorts, &pb.PortStatus{
			PortNumber:  int32(port),
			IsAvailable: true,
		})

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterDataKeeperServer(grpcServer, dataKeeperServer)
		log.Printf("Data Keeper gRPC server started on port %d", port)
		go grpcServer.Serve(lis)
		counter++
	}

	dataKeeperServer.SetPorts(tcpPorts, grpcPorts)

	nodeName := uuid.New().String()
	globals.masterAddress = masterAddress
	globals.nodeName = nodeName
	globals.nodeAddress = nodeAddress

	go dataKeeperServer.startHeartbeat(masterAddress, nodeAddress, nodeName)
	select {} // Keep the process running
}

func (s *DataKeeperServer) SetPorts(tcpPorts []*pb.PortStatus, grpcPorts []*pb.PortStatus) {
	s.portsTcp = tcpPorts
	s.portsGrpc = grpcPorts
}

func (s *DataKeeperServer) ReplicateFile(ctx context.Context, req *pb.ReplicationRequest) (*pb.FileUploadSuccess, error) {
	log.Printf("[DATA KEEPER] Replicating file %s to %s", req.Filename, req.DestinationAddress)
	filename := req.Filename
	destinationAddress := req.DestinationAddress
	wd, _ := os.Getwd()
	filePath := path.Join( "storage", req.Filename)
	absolutepath := path.Join(wd, filePath)

	log.Printf("[DATA KEEPER] Replicating file %s to %s", filename, destinationAddress)

	// Check if the file exists.
	if _, err := os.Stat(absolutepath); os.IsNotExist(err) {
		log.Printf("[REPLICATION ERROR] File %s not found", filename)
		return &pb.FileUploadSuccess{}, err
	}

	// Open connection to destination Data Keeper.
	conn, err := net.Dial("tcp", destinationAddress)
	if err != nil {
		log.Printf("[REPLICATION ERROR] Cannot connect to destination %s: %v", destinationAddress, err)
		return &pb.FileUploadSuccess{}, err
	}
	defer conn.Close()

	// Notify the destination about the upcoming file upload.
	writer := bufio.NewWriter(conn)
	writer.WriteString("REPLICATE\n")
	writer.WriteString(filename + "\n")
	writer.Flush()

	// Send the file.
	file, err := os.Open(absolutepath)
	if err != nil {
		log.Printf("[REPLICATION ERROR] Cannot open file %s: %v", filename, err)
		return &pb.FileUploadSuccess{}, err
	}
	defer file.Close()

	if _, err := io.Copy(conn, file); err != nil {
		log.Printf("[REPLICATION ERROR] Failed to send file: %v", err)
		return &pb.FileUploadSuccess{}, err
	}

	log.Printf("[REPLICATION SUCCESS] File %s replicated to %s", filename, destinationAddress)
	port, err := utils.ExtractPort(destinationAddress)
	if err != nil {
		return &pb.FileUploadSuccess{}, err
	}
	return &pb.FileUploadSuccess{DataKeeperName: req.DestinationName, Filename: filename, FilePath: filePath, PortNumber: int32(port)}, nil
}

// create the gRPC connection and start sending heartbeat
func (s *DataKeeperServer) startHeartbeat(masterAddress string, dataNodeAddress string, name string) {
	// create connection
	conn, err := grpc.Dial(masterAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Could not connect to Master Tracker: %v", err)
		return
	}

	// free resources
	defer conn.Close()

	// create gRPC client
	client := pb.NewMasterTrackerClient(conn)

	// send heartbeat every second till process stops
	// TODO: We need here to use alternative to time.Sleep (events) :: EMAD
	for {
		sendHeartbeat(client, name, dataNodeAddress, s.portsTcp, s.portsGrpc)
		time.Sleep(1 * time.Second)
	}
}

func startTCPServer(portStatus *pb.PortStatus) {
	address := fmt.Sprintf(":%d", portStatus.PortNumber)
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

		go handleTcpRequest(conn)
	}
}

// notify master that this data node is alive
func sendHeartbeat(client pb.MasterTrackerClient, name string, dataNodeAddress string, portsTcp []*pb.PortStatus, portsGrpc []*pb.PortStatus) {
	// Context with 1-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Ensure resources are freed

	// Send heartbeat message
	_, err := client.SendHeartbeat(ctx, &pb.HeartbeatRequest{
		DataKeeperName:    name,
		DataKeeperAddress: dataNodeAddress,
		PortsTCP:          portsTcp,
		PortsGRPC:         portsGrpc,
	})
	if err != nil {
		log.Printf("Heartbeat error: %v", err)
	} else {
		log.Println("Heartbeat sent successfully")
	}
}

// HandleFileUpload processes file uploads from the client
func HandleFileUpload(reader *bufio.Reader, conn net.Conn, isUpload bool) {
	filename, err := readFilename(reader)
	if err != nil {
		log.Println(err)
		return
	}

	utils.EnsureStorageFolder()
	filePath := path.Join("storage", filename)

	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()
	defer conn.Close()
	if _, err = io.Copy(file, conn); err != nil {
		log.Printf("Failed to save file %s: %v\n", filePath, err)
		return
	}
	log.Printf("File %s received and saved successfully!\n", filename)
	remoteAddr := conn.RemoteAddr().String()
	remotePort, _ := utils.ExtractPort(remoteAddr)

	// Notify the master that the file has been uploaded
	if isUpload {
		SendUploadSuccessResponse(filename, filePath, int32(remotePort))
	}
}

func SendUploadSuccessResponse(filename, filePath string, portNumber int32) {
	conn, err := grpc.Dial(globals.masterAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := pb.NewMasterTrackerClient(conn)
	_, err = client.RequestUploadSuccess(context.Background(), &pb.FileUploadSuccess{
		DataKeeperName: globals.nodeName,
		Filename:       filename,
		FilePath:       filePath,
		PortNumber:     portNumber,
	})
	if err != nil {
		log.Printf("Failed to notify master about file upload: %v\n", err)
		return
	}
	log.Printf("Notified master about file upload: %s\n", filename)
	defer conn.Close()
}

func handleTcpRequest(conn net.Conn) {
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
		HandleFileUpload(reader, conn, true)
	} else if requestType == "REPLICATE" {
		HandleFileUpload(reader, conn, false)
	} else if requestType == "DOWNLOAD" {
		HandleFileDownload(reader, conn)
	} else {
		log.Println("Invalid request type:", requestType)
	}
	defer conn.Close()
}

// HandleFileDownload processes file download requests
func HandleFileDownload(reader *bufio.Reader, conn net.Conn) {
	filename, err := readFilename(reader)
	if err != nil {
		log.Println(err)
		return
	}

	filePath := path.Join("storage", filename)
	fmt.Printf("filePath: %s\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("File not found: %s\n", filename)
		sendResponse(conn, "ERROR: File not found")
		return
	}

	log.Printf("Sending file: %s\n", filePath)
	writer := bufio.NewWriter(conn)
	// Send file data
	if _, err = io.Copy(writer, file); err != nil {
		log.Printf("Error sending file %s: %v\n", filename, err)
	} else {
		log.Println("File sent successfully!")
	}

	defer file.Close()
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

func isPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false // Port is in use
	}
	listener.Close() // Close immediately after checking
	return true      // Port is available
}
