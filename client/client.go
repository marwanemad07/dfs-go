package main

import (
	"bufio"
	"context"
	"dfs/config"
	pb "dfs/proto"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: client <upload/download> <filename>")
		return
	}

	command := os.Args[1]
	filename := os.Args[2]

	// Connect to the Master Tracker
	serverPort := config.LoadConfig("config.json").Server.Port
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", serverPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to Master Tracker: %v", err)
	}
	defer conn.Close()
	client := pb.NewMasterTrackerClient(conn)

	switch command {
	case "upload":
		uploadFile(client, filename)
	case "download":
		downloadFile(client, filename)
	default:
		fmt.Println("Invalid command. Use 'upload' or 'download'.")
	}
}

// Upload logic
func uploadFile(master pb.MasterTrackerClient, filename string) {
	fmt.Println("Uploading file:", filename)

	// Ensure only MP4 files are allowed
	if !strings.HasSuffix(strings.ToLower(filename), ".mp4") {
		log.Fatalf("Only MP4 files are allowed for upload")
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	fullPath := filepath.Join(dir, filename)

	// Request Data Keeper from Master Tracker
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	uploadResp, err := master.RequestUpload(ctx, &pb.UploadRequest{Filename: fullPath})
	if err != nil {
		log.Fatalf("Error requesting upload: %v", err)
	}

	dataKeeperPort := uploadResp.DataKeeperAddress
	dataKeeperAddr := getAddress(dataKeeperPort) // "localhost:" need to parameter read from config file
	log.Printf("Uploading to Data Keeper at %s", dataKeeperAddr)

	// Send file to Data Keeper via TCP
	SendFile(dataKeeperAddr, filename)
}

// Download logic
func downloadFile(client pb.MasterTrackerClient, filename string) {
	fmt.Println("Downloading file:", filename)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	fullPath := filepath.Join(dir, filename)

	// Request file locations from Master Tracker
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	downloadResp, err := client.RequestDownload(ctx, &pb.DownloadRequest{Filename: fullPath})
	if err != nil {
		log.Fatalf("Error requesting download: %v", err)
	}

	if len(downloadResp.DataKeeperAddresses) == 0 {
		log.Fatalf("No Data Keeper has the requested file: %s", filename)
	}

	// Select the first available Data Keeper
	dataKeeperPort := downloadResp.DataKeeperAddresses[0]
	dataKeeperAddr := getAddress(dataKeeperPort) // "localhost:" need to parameter read from config file
	log.Printf("Downloading from Data Keeper at %s", dataKeeperAddr)

	// Request and receive file via TCP
	ReceiveFile(dataKeeperAddr, filename)
}

func getAddress(port string) string {
	return "localhost:" + port
}

// SendFile uploads a file to the server
func SendFile(address, filename string) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Connect to Data Keeper
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect to Data Keeper: %v", err)
	}
	defer conn.Close()

	// Send request header
	if err := sendRequest(conn, "UPLOAD", filename); err != nil {
		log.Fatalf("%v", err)
	}

	// Send file data
	writer := bufio.NewWriter(conn)
	if _, err := io.Copy(writer, file); err != nil {
		log.Fatalf("Failed to send file data: %v", err)
	}

	// Ensure all data is sent
	if err := writer.Flush(); err != nil {
		log.Fatalf("Failed to flush file data: %v", err)
	}

	log.Println("Upload complete!")
}

// ReceiveFile downloads a file from the server
func ReceiveFile(address, filename string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect to Data Keeper: %v", err)
	}
	defer conn.Close()

	// Send request header
	if err := sendRequest(conn, "DOWNLOAD", filename); err != nil {
		log.Fatalf("%v", err)
	}

	// Read server response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read server response: %v", err)
	}

	// Check if the server returned an error
	if strings.HasPrefix(response, "ERROR") {
		log.Printf("Server error: %s", response)
		return
	}

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Receive and save file data
	if _, err := io.Copy(file, conn); err != nil {
		log.Fatalf("Error receiving file: %v", err)
	}

	log.Println("File received successfully.")
}

// sendRequest sends a request type and filename to the server
func sendRequest(conn net.Conn, requestType, filename string) error {
	writer := bufio.NewWriter(conn)

	// Send request type
	if _, err := writer.WriteString(requestType + "\n"); err != nil {
		return fmt.Errorf("failed to send request type: %w", err)
	}

	// Send filename
	if _, err := writer.WriteString(filename + "\n"); err != nil {
		return fmt.Errorf("failed to send filename: %w", err)
	}

	// Flush data
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush request data: %w", err)
	}

	return nil
}
