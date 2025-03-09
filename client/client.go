package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"dfs/config"
	pb "dfs/proto"
	"dfs/utils"
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
func uploadFile(client pb.MasterTrackerClient, filename string) {
	fmt.Println("Uploading file:", filename)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	fullPath := filepath.Join(dir, filename)

	// Request Data Keeper from Master Tracker
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	uploadResp, err := client.RequestUpload(ctx, &pb.UploadRequest{Filename: fullPath})
	if err != nil {
		log.Fatalf("Error requesting upload: %v", err)
	}

	dataKeeperPort := uploadResp.DataKeeperAddress
	dataKeeperAddr := "localhost:" + dataKeeperPort // "localhost:" need to parameter read from config file
	log.Printf("Uploading to Data Keeper at %s", dataKeeperAddr)


	// Send file to Data Keeper via TCP
	utils.SendFile(dataKeeperAddr, filename)
}

// Download logic
func downloadFile(client pb.MasterTrackerClient, filename string) {
	fmt.Println("Downloading file:", filename)

	// Request file locations from Master Tracker
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	downloadResp, err := client.RequestDownload(ctx, &pb.DownloadRequest{Filename: filename})
	if err != nil {
		log.Fatalf("Error requesting download: %v", err)
	}

	if len(downloadResp.DataKeeperAddresses) == 0 {
		log.Fatalf("No Data Keeper has the requested file: %s", filename)
	}

	// Select the first available Data Keeper
	dataKeeperAddr := downloadResp.DataKeeperAddresses[0]
	log.Printf("Downloading from Data Keeper at %s", dataKeeperAddr)

	// Request and receive file via TCP
	utils.ReceiveFile(dataKeeperAddr, filename)
}
