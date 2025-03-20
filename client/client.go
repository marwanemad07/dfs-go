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
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/grpc"
)
func main() {
	// Define flags
	outputPath := flag.String("o", "", "Output path for downloaded file")
	flag.Parse()

	// Ensure there are enough arguments
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: client <upload/download> <filename> [-o outputPath]")
		return
	}

	command := args[0]
	filename := args[1]

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
		if *outputPath == "" {
			*outputPath = "downloads" 
		}
		downloadFile(client, filename, *outputPath)
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


	// Request Data Keeper from Master Tracker
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	uploadResp, err := master.RequestUpload(ctx, &pb.UploadRequest{Filename: filename})
	if err != nil {
		log.Fatalf("Error requesting upload: %v", err)
	}

	log.Printf("Uploading to Data Keeper at %s", uploadResp.DataKeeperAddress)

	// Send file to Data Keeper via TCP
	SendFile(uploadResp.DataKeeperAddress, filename)
}

// Download logic
func downloadFile(client pb.MasterTrackerClient, filename string,filePath string) {
	fmt.Println("Downloading file:", filename)

	// Request file locations from Master Tracker
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	downloadResp, err := client.RequestDownload(ctx, &pb.DownloadRequest{Filename: filename})
	if err != nil {
		log.Fatalf("Error requesting download: %v", err)
	}
	lenOFDataKeeperAddresses := len(downloadResp.DataKeeperAddresses)
	if lenOFDataKeeperAddresses == 0 {
		// TODO: Should request to download after some time or tell user to try again later
		log.Fatalf("No Data Keeper has the requested file: %s", filename)
	}

	// Choose a random Data Keeper uniformaly to download from
	index := utils.GetRandomIndex(lenOFDataKeeperAddresses)
	dataKeeperAddress := downloadResp.DataKeeperAddresses[index]
	log.Printf("Downloading from Data Keeper at %s", dataKeeperAddress)

	// Request and receive file via TCP
	ReceiveFile(dataKeeperAddress, filename,filePath)
}


// SendFile uploads a file to the server
func SendFile(address, filename string) {
	fmt.Println("Uploading file:", filename)
	// Open the file

	filePath := filepath.Join(utils.GetWorkingDir(), filename)
	fmt.Println("Uploading file:", filePath,address)
	file, err := os.Open(filePath)
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
func ReceiveFile(address, filename, filePath string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect to Data Keeper: %v", err)
	}
	defer conn.Close()

	// Ensure the directory exists
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Full path where the file will be saved
	fullFilePath := filepath.Join(filePath, filename)

	// Send request header
	if err := sendRequest(conn, "DOWNLOAD", filename); err != nil {
		log.Fatalf("%v", err)
	}

	// Create the file in the specified directory
	file, err := os.Create(fullFilePath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Receive and save file data
	if _, err := io.Copy(file, conn); err != nil {
		log.Fatalf("Error receiving file: %v", err)
	}

	log.Printf("File received successfully: %s\n", fullFilePath)
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
