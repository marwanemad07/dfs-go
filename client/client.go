package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

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
	sendFile(dataKeeperAddr, filename)
}

// sendFile connects to Data Keeper and transfers the file over TCP
func sendFile(address, filename string) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Connect to Data Keeper over TCP
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect to Data Keeper: %v", err)
	}
	defer conn.Close()

	// Send filename first
	_, err = conn.Write([]byte(filename + "\n"))
	if err != nil {
		log.Fatalf("Failed to send filename: %v", err)
	}

	// Send file data
	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatalf("Failed to send file data: %v", err)
	}

	fmt.Println("Upload complete!")
}
