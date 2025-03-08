package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"dfs/utils"
	pb "dfs/proto"

	"google.golang.org/grpc"
)

func main() {
	// Connect to the Master Tracker
	conn, err := grpc.Dial("localhost:50050", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to Master Tracker: %v", err)
	}
	defer conn.Close()
	client := pb.NewMasterTrackerClient(conn)

	// File to upload
	filename := "example.mp4"
	fmt.Println("Uploading file:", filename)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	fullPath := filepath.Join(dir, filename)


	// Request a Data Keeper for upload
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	uploadResp, err := client.RequestUpload(ctx, &pb.UploadRequest{Filename: fullPath})
	if err != nil {
		log.Fatalf("Error requesting upload: %v", err)
	}
	dataKeeperAddr := uploadResp.DataKeeperAddress
	log.Printf("Uploading to Data Keeper at %s", dataKeeperAddr)

	// Send file via TCP
	utils.SendFile(dataKeeperAddr, filename)
}

// sendFile connects to Data Keeper and transfers the file over TCP
