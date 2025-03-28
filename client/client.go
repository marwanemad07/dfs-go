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
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	pb.UnimplementedClientServer
	done chan struct{} // Add a channel to signal completion
}

func main() {
	// Define flags
	outputPath := flag.String("o", "", "Output path for downloaded file")
	grpcPort := flag.String("p", "6800", "client port")
	networkType := flag.String("n", "localhost", "type of network")
	flag.Parse()

	// Ensure there are enough arguments
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: client [-o outputPath] [-p grpcPort] <upload/download> <filename> ")
		return
	}

	command := args[0]
	filename := args[1]

	address := "localhost"
	if *networkType == "network" {
	
	address,_ = utils.GetWiFiIPv4()
	}
	// Connect to the Master Tracker
	serverPort := config.LoadConfig("config.json").Server.Port
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", serverPort), grpc.WithInsecure())
	
	// Start the gRPC server
	
	responseAddress := address + ":" + *grpcPort
	lis, err := net.Listen("tcp", responseAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	clienObject := &Client{done: make(chan struct{})}
	grpcServer := grpc.NewServer()
	pb.RegisterClientServer(grpcServer, clienObject)
	go grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Could not connect to Master Tracker: %v", err)
	}
	defer conn.Close()
	client := pb.NewMasterTrackerClient(conn)


	switch command {
	case "upload":
		uploadFile(client, filename, responseAddress)
	case "download":
		if *outputPath == "" {
			*outputPath = "downloads" 
		}
		
		downloadFile(client, filename, *outputPath)
		close(clienObject.done)

	default:
		fmt.Println("Invalid command. Use 'upload' or 'download'.")
	}
		
	<-clienObject.done
    log.Printf("Client shutting down gracefully")
}


func (c *Client) NotifyUploadCompletion(ctx context.Context, req *pb.UploadSuccessResponse) (*emptypb.Empty, error) {
	if(req.Success){
		log.Printf("Upload completed successfully for file")
		close(c.done)
	}else{
		log.Printf("Upload failed for file")
		// Retry logic I think will be implemented here
	}
	return &emptypb.Empty{}, nil
}

// Upload logic
func uploadFile(master pb.MasterTrackerClient, filename string, responseAddress string) {
	fmt.Println("Uploading file:", filename)

	// Ensure only MP4 files are allowed
	if !strings.HasSuffix(strings.ToLower(filename), ".mp4") {
		log.Fatalf("Only MP4 files are allowed for upload")
	}


	// Request Data Keeper from Master Tracker
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	uploadResp, err := master.RequestUpload(ctx, &pb.UploadRequest{Filename: filename})
	if err != nil {
		log.Fatalf("Error requesting upload: %v", err)
	}

	log.Printf("Uploading to Data Keeper at %s", uploadResp.DataKeeperAddress)

	// Send file to Data Keeper via TCP
	conn, err := net.Dial("tcp", uploadResp.DataKeeperAddress)
	if err != nil {
		log.Fatalf("Failed to connect to Data Keeper: %v", err)
	}
	SendFile(filename, conn, responseAddress)
	conn.Close()
}

// Download logic
func downloadFile(client pb.MasterTrackerClient, filename string,filePath string) {
	fmt.Println("Downloading file:", filename)

	// Request file locations from Master Tracker
	ctx, cancel := context.WithTimeout(context.Background(),10*time.Minute)
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
func SendFile(filename string, conn net.Conn, responseAddress string) {
	// Open the file
	filePath := filepath.Join(utils.GetWorkingDir(), filename)
	
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
		utils.SendResponse(conn, "ERROR: File not found")
		return
	}
	defer file.Close()
	
	// Send request header
	if err := utils.SendRequest(conn, "UPLOAD#" + responseAddress, filename); err != nil {
		log.Fatalf("%v", err)
		return
	}

	// Send file data
	n, err := utils.WriteFileToConnection(file, conn)
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		return
	}
	log.Printf("Upload complete! Sent %d bytes", n)
}

// ReceiveFile downloads a file from the server
func ReceiveFile(address string, filename string, filePath string) {
	conn, err := net.Dial("tcp", address)
	reader := bufio.NewReader(conn)

	if err != nil {
		log.Fatalf("Failed to connect to Data Keeper: %v", err)
	}
	defer conn.Close()

	utils.EnsureStorageFolder(filePath)

	fullFilePath := filepath.Join(filePath, filename)

	if err := utils.SendRequest(conn, "DOWNLOAD", filename); err != nil {
		log.Fatalf("Failed to send download request: %v", err)
	}

	file, err := os.Create(fullFilePath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}

	totalSizeStr, err := getResponse(reader)
	totalSize,_ := strconv.ParseInt(totalSizeStr, 10, 64)
	log.Printf("File size: %d bytes\n", totalSize)
	if err != nil {
		log.Fatalf("Failed to read file size: %v", err)
	}

	receiveAndSaveFile(reader, file, totalSize)
}

// getFileSize reads the expected file size from the server.
func getResponse(reader *bufio.Reader) (string, error) {
	str, err := reader.ReadString('\n')
	if err != nil {
		return "0", fmt.Errorf("failed to read Response: %w", err)
	}
	str = strings.TrimSpace(str)
	return str, err
}

func receiveAndSaveFile(reader *bufio.Reader, file *os.File, totalSize int64) error {
	progress := make(chan int64)
	defer file.Close()

	// Show download progress in a separate goroutine
	go utils.ShowProgress(progress, totalSize)

	buffer := make([]byte, 4096) // Use a reasonable buffer size (4 KB)
	var received int64

	for received < totalSize {
		toRead := int64(len(buffer))
		if remaining := totalSize - received; remaining < toRead {
			toRead = remaining
		}

		// Read exactly `toRead` bytes
		n, err := io.ReadFull(reader, buffer[:toRead])
		if n > 0 {
			if _, writeErr := file.Write(buffer[:n]); writeErr != nil {
				return fmt.Errorf("failed to write to file: %w", writeErr)
			}
			received += int64(n)
			progress <- received
		}

		if err == io.EOF {
			break // End of file reached
		}

		if err != nil && err != io.ErrUnexpectedEOF {
			return fmt.Errorf("error while receiving file: %w", err)
		}
	}

	close(progress)

	// Validate that we received the full file
	if received != totalSize {
		return fmt.Errorf("incomplete file received. Expected %d bytes, got %d bytes", totalSize, received)
	}

	return nil
}