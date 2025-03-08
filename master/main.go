package main

import (
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	pb "dfs/proto/dfs" // update this import path
)

// server implements the MasterTracker sdfervice.
type server struct {
	pb.UnimplementedMasterTrackerServer
	// For simplicity, we use a map to keep track of file storage info.
	// In a complete implementation, you would track file location, replication status, etc.
	files map[string]string
	mu    sync.Mutex
}

// NewServer returns a new Master Tracker server.
func NewServer() *server {
	return &server{
		files: make(map[string]string),
	}
}

// RequestUpload assigns a Data Keeper for the file upload.
func (s *server) RequestUpload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	// For now, we choose a dummy Data Keeper address.
	address := "localhost:50051"
	log.Printf("Upload requested for file: %s, assigning Data Keeper: %s", req.Filename, address)
	return &pb.UploadResponse{DataKeeperAddress: address}, nil
}

// RequestDownload returns a list of Data Keeper addresses that have the file.
func (s *server) RequestDownload(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	// For now, we return a single hardcoded address.
	addresses := []string{"localhost:50051"}
	log.Printf("Download requested for file: %s, returning addresses: %v", req.Filename, addresses)
	return &pb.DownloadResponse{DataKeeperAddresses: addresses}, nil
}

// SendHeartbeat updates the server with a heartbeat from a Data Keeper.
func (s *server) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	log.Printf("Received heartbeat from Data Keeper: %s", req.DataKeeperId)
	// Here you could update the lookup table to mark the node as alive.
	return &pb.HeartbeatResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMasterTrackerServer(grpcServer, NewServer())
	log.Println("Master Tracker is listening on port 50050")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
