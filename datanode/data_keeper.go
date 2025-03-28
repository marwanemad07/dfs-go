package main

import (
	"bufio"
	"context"
	"dfs/config"
	pb "dfs/proto"
	"dfs/utils"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	masterAddress := *flag.String("m", "localhost", "master address")

	if len(os.Args) < 5 {
		fmt.Println("Usage: go run main.go <master_address_type> <number_of_ports> <start_tcp_ports> <node_name> ...")
		return
	}

	// Start TCP Server for receiving files
	nodeAddress := ""

	switch os.Args[1] {
	case "local":
		nodeAddress = "localhost"
	case "docker":
		masterAddress = "host.docker.internal"
		nodeAddress = "host.docker.internal"
	case "network":
		nodeAddress,_ = utils.GetWiFiIPv4()
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
		for !utils.IsPortAvailable(port) {
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
		for !utils.IsPortAvailable(port) {
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

	nodeName := os.Args[4]
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
	filePath := path.Join("storage", req.Filename)
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

	if _, err := io.Copy(writer, file); err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute * 10)
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
func HandleFileUpload(filename string, reader *bufio.Reader, conn net.Conn, isUpload bool, clientAddress ...string) {
	if !isUpload {
		filename = globals.nodeName + filename
	}
	filePath := GetFilePath(filename)
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file %s: %v\n", filename, err)
		return
	}
	defer file.Close()
	if size, err := io.Copy(file, reader); err != nil {
		log.Printf("Failed to save file %s: %v\n", file.Name(), err)
		return
	} else {
		log.Print("File size:", size)
	}
	log.Printf("File %s received and saved successfully!\n", filename)

	// Notify the master that the file has been uploaded
	if isUpload {
		remoteAddr := ""
		if(len(clientAddress) > 0) {
			remoteAddr = clientAddress[0]
		} else{
			log.Println("Client address not provided")
			return
		}

		relativePath := path.Join("storage", filepath.Base(file.Name()))
		SendUploadSuccessResponse(filename, relativePath, remoteAddr)
	}
}

func SendUploadSuccessResponse(filename, filePath string, address string) {
	portNumber, _ := utils.ExtractPort(address)
	fmt.Println("address:", address)	
	conn, err := grpc.Dial(globals.masterAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := pb.NewMasterTrackerClient(conn)
	_, err = client.RequestUploadSuccess(context.Background(), &pb.FileUploadSuccess{
		DataKeeperName: globals.nodeName,
		Filename:       filename,
		FilePath:       filePath,
		PortNumber:     int32(portNumber),
		ClientAddress:        address,
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
	if err != nil {
		log.Println("Failed to read request type:", err)
		return
	}
	requestType = strings.TrimSpace(requestType)
	log.Println("Request type:", requestType)

	filename, err := readFilename(reader)
	if err != nil {
		log.Println(err)
		return
	}

	// Check if it's an upload or download request
	if strings.HasPrefix(requestType, "UPLOAD") {
		address := strings.Split(requestType, "#")[1]
		HandleFileUpload(filename, reader, conn, true, address)
	} else if requestType == "REPLICATE" {
		HandleFileUpload(filename, reader, conn, false)
	} else if requestType == "DOWNLOAD" {
		HandleFileDownload(filename, conn)
	} else {
		log.Println("Invalid request type:", requestType)
	}
	conn.Close()
}

// HandleFileDownload processes file download requests
func HandleFileDownload(filename string, conn net.Conn) {
	filePath := GetFilePath(filename)
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("File not found: %s\n", filename)
		utils.SendResponse(conn, "ERROR: File not found")
		return
	}
	defer file.Close()

	log.Printf("Sending file: %s\n", file.Name())
	fileInfo, _ := file.Stat()
	fmt.Println("File size:", fileInfo.Size())

	// hasher := sha256.New()
	// io.Copy(hasher, file)
	// checksum := fmt.Sprintf("%x", hasher.Sum(nil))
	file.Seek(0, io.SeekStart)

	utils.SendResponse(conn, strconv.FormatInt(fileInfo.Size(), 10))
	utils.WriteFileToConnection(file, conn)
}

// readFilename reads and trims the filename from the connection
func readFilename(reader *bufio.Reader) (string, error) {
	filename, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read filename: %w", err)
	}
	return strings.TrimSpace(filename), nil
}

func GetFilePath(filename string) string {
	utils.EnsureStorageFolder("storage")
	filePath := path.Join(utils.GetWorkingDir(), "storage", filename)
	return filePath
}
