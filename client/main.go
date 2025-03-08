package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "dfs/proto" // update this import path
)

func main() {
	// Connect to the Master Tracker.
	conn, err := grpc.Dial("localhost:50050", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to Master Tracker: %v", err)
	}
	defer conn.Close()
	client := pb.NewMasterTrackerClient(conn)

	// Request to upload a file.
	filename := "example.mp4"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	uploadResp, err := client.RequestUpload(ctx, &pb.UploadRequest{Filename: filename})
	if err != nil {
		log.Fatalf("Error during upload request: %v", err)
	}
	log.Printf("Received Data Keeper address for upload: %s", uploadResp.DataKeeperAddress)

	// Simulate the file upload via TCP (the actual TCP logic is omitted for brevity).
	fmt.Printf("Uploading %s to Data Keeper at %s...\n", filename, uploadResp.DataKeeperAddress)
	time.Sleep(2 * time.Second) // simulate file transfer delay
	fmt.Println("Upload complete. Data Keeper will notify the Master Tracker.")

	// In a full implementation, after the TCP transfer, the Data Keeper would notify the Master Tracker.
}
