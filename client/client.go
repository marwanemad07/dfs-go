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
	fmt.Println("Client started"," on port", *grpcPort)
	// Ensure there are enough arguments
	args := flag.Args()
	if len(args) < 3 {
		fmt.Println("Usage: client [-o outputPath] [-p grpcPort] <upload/download> <filename> <masterAddress>")
		return
	}

	command := args[0]
	filename := args[1]
	masterAddress := args[2]

	address := "localhost"
	if *networkType == "network" {
	
	address,_ = utils.GetWiFiIPv4()
	}
	// Connect to the Master Tracker
	serverPort := config.LoadConfig("config.json").Server.Port
	conn, err := grpc.Dial(masterAddress + ":" + strconv.Itoa(serverPort), grpc.WithInsecure())
	log.Println(masterAddress + ":" + strconv.Itoa(serverPort))
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
		clienObject.uploadFile(client, filename, responseAddress)
	case "download":
		if *outputPath == "" {
			*outputPath = "downloads" 
		}
		
		downloadFile(client, filename, *outputPath)
		close(clienObject.done)

	default:
		fmt.Println("Invalid command. Use 'upload' or 'download'.")
		return
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
func (c *Client) uploadFile(master pb.MasterTrackerClient, filename string, responseAddress string) {
    fmt.Println("Uploading file:", filename)

    // Ensure only MP4 files are allowed
    if !strings.HasSuffix(strings.ToLower(filename), ".mp4") {
        log.Fatalf("Only MP4 files are allowed for upload")
    }

    maxAttempts := 3
	attemptUpload := func() (bool, error) {
        // Request Data Keeper from Master Tracker
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
        defer cancel()
        uploadResp, err := master.RequestUpload(ctx, &pb.UploadRequest{Filename: filename})
        if err != nil {
            return false, fmt.Errorf("request upload failed: %v", err)
        }

        log.Printf("Uploading to Data Keeper at %s", uploadResp.DataKeeperAddress)

        conn, err := net.Dial("tcp", uploadResp.DataKeeperAddress)
        if err != nil {
            return false, fmt.Errorf("connect to Data Keeper failed: %v", err)
        }
        defer conn.Close() // Close connection when done

        if !SendFile(filename, conn, responseAddress) {
            return false, fmt.Errorf("send file failed")
        }

        log.Printf("File %s uploaded successfully to %s", filename, uploadResp.DataKeeperAddress)
        return true, nil
    }

    for attempt := 0; attempt < maxAttempts; attempt++ {
        success, err := attemptUpload()
        if success {
            return 
        }

        log.Printf("Attempt %d/%d: %v", attempt+1, maxAttempts, err)
        if attempt == maxAttempts-1 {
            log.Println("All attempts failed, giving up")
            close(c.done) 
            return
        }

        // Delay before next attempt (0s, 5s, 10s)
        delay := time.Duration(attempt*5+5) * time.Second
        log.Printf("Waiting %v before retrying...", delay)
        time.Sleep(delay)
    }
}
// Download logic
func downloadFile(client pb.MasterTrackerClient, filename string, filePath string) {
    fmt.Println("Downloading file:", filename)

    maxAttempts := 3

    attemptDownload := func() (bool, error) {
        // Request file locations from Master Tracker
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
        defer cancel()
        downloadResp, err := client.RequestDownload(ctx, &pb.DownloadRequest{Filename: filename})
        if err != nil {
            return false, fmt.Errorf("error requesting download: %v", err)
        }

        lenOfDataKeeperAddresses := len(downloadResp.DataKeeperAddresses)
        if lenOfDataKeeperAddresses == 0 {
            return false, fmt.Errorf("no Data Keeper has the requested file: %s", filename)
        }

        // Choose a random Data Keeper uniformly to download from
        index := utils.GetRandomIndex(lenOfDataKeeperAddresses)
        dataKeeperAddress := downloadResp.DataKeeperAddresses[index]
        log.Printf("Downloading from Data Keeper at %s", dataKeeperAddress)

        // Request and receive file via TCP
        if ReceiveFile(dataKeeperAddress, filename, filePath) {
			return true, nil
		}
		return false, fmt.Errorf("failed to download file from Data Keeper: %s", dataKeeperAddress)
    }

    for attempt := 0; attempt < maxAttempts; attempt++ {
        success, err := attemptDownload()
        if success {
            log.Printf("Successfully downloaded %s", filename)
            return // Success, exit the function
        }

        log.Printf("Attempt  %d/%d: %v", attempt+1, maxAttempts, err)
        if attempt == maxAttempts-1 {
            log.Printf("All attempts failed to download %s, giving up", filename)
            return
        }
		delay := time.Duration(attempt*5+5) * time.Second

        log.Printf("Waiting %v before retrying...", delay)
        time.Sleep(delay)
    }
}

// SendFile uploads a file to the server
func SendFile(filename string, conn net.Conn, responseAddress string) (bool) {
	// Open the file
	filePath := filepath.Join(utils.GetWorkingDir(), filename)
	
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
		utils.SendResponse(conn, "ERROR: File not found")
		return false
	}
	defer file.Close()
	
	// Send request header
	if err := utils.SendRequest(conn, "UPLOAD#" + responseAddress, filename); err != nil {
		log.Fatalf("%v", err)
		return false
	}
	log.Printf("upload request recived successfully", )

	// Send file data
	n, err := utils.WriteFileToConnection(file, conn)
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		return false
	}
	log.Printf("Upload complete! Sent %d bytes", n)
	return true
}

// ReceiveFile downloads a file from the server
func ReceiveFile(address string, filename string, filePath string)(bool) {
    maxAttempts := 3

    attemptDownload := func(attempt int) (bool, error) {
        // Attempt to connect to the Data Keeper
        conn, err := net.Dial("tcp", address)
        if err != nil {
            return false, fmt.Errorf("failed to connect to Data Keeper: %v", err)
        }
        defer conn.Close()

        reader := bufio.NewReader(conn)
        utils.EnsureStorageFolder(filePath)
        fullFilePath := filepath.Join(filePath, filename)

        // Send download request
        if err := utils.SendRequest(conn, "DOWNLOAD", filename); err != nil {
            return false, fmt.Errorf("failed to send download request: %v", err)
        }

        // Create the file
        file, err := os.Create(fullFilePath)
        if err != nil {
            return false, fmt.Errorf("failed to create file: %v", err)
        }
        defer file.Close()

        // Get file size
        totalSizeStr, err := getResponse(reader)
        if err != nil {
            return false, fmt.Errorf("failed to read file size: %v", err)
        }
        totalSize, err := strconv.ParseInt(totalSizeStr, 10, 64)
        if err != nil {
            return false, fmt.Errorf("failed to parse file size: %v", err)
        }
        log.Printf("File size: %d bytes", totalSize)

        // Receive and save the file
        if err := receiveAndSaveFile(reader, file, totalSize); err != nil {
            return false, fmt.Errorf("failed to receive file: %v", err)
        }

        log.Printf("Successfully downloaded %s to %s", filename, fullFilePath)
        return true, nil
    }

    for attempt := 0; attempt < maxAttempts; attempt++ {
        success, err := attemptDownload(attempt)
        if success {
            return true // Success, exit the function
        }

        log.Printf("Attempt %d/%d: %v", attempt+1, maxAttempts, err)
        if attempt == maxAttempts-1 {
            log.Printf("All attempts failed to download %s, Trying differnt data node", filename)
            return false
        }
		delay := time.Duration(attempt*5+5) * time.Second

        log.Printf("Waiting %v before retrying...", delay)
        time.Sleep(delay)
    }
	return true
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

	buffer := make([]byte, 4096) 
	var received int64

	for received < totalSize {
		toRead := int64(len(buffer))
		if remaining := totalSize - received; remaining < toRead {
			toRead = remaining
		}

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