package main

import (
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	pb "dfs/proto"
)

// server implements the Master Tracker service.
type server struct {
	pb.UnimplementedMasterTrackerServer
	files map[string]string
	mu    sync.Mutex
}

func NewServer() *server {
	return &server{
		files: make(map[string]string),
	}
}

// RequestUpload assigns a Data Keeper for file upload.
func (s *server) RequestUpload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	address := "localhost:50051" // Hardcoded for now

	s.mu.Lock()
	s.files[req.Filename] = address
	s.mu.Unlock()

	log.Printf("[UPLOAD] File: %s assigned to Data Keeper: %s", req.Filename, address)
	return &pb.UploadResponse{DataKeeperAddress: address}, nil
}

// RequestDownload returns a Data Keeper storing the file.
func (s *server) RequestDownload(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	s.mu.Lock()
	address, exists := s.files[req.Filename]
	s.mu.Unlock()

	if !exists {
		log.Printf("[ERROR] File %s not found!", req.Filename)
		return nil, grpc.Errorf(404, "File not found")
	}

	log.Printf("[DOWNLOAD] File %s is available at %s", req.Filename, address)
	return &pb.DownloadResponse{DataKeeperAddresses: []string{address}}, nil
}

// SendHeartbeat updates Master Tracker with Data Keeper status.
func (s *server) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	log.Printf("[HEARTBEAT] Data Keeper: %s is alive", req.DataKeeperId)
	return &pb.HeartbeatResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMasterTrackerServer(grpcServer, NewServer())
	log.Println("Master Tracker is running on port 50050 🚀")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
