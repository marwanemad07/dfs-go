package main

import (
	"context"
	pb "dfs/proto"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"dfs/utils"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"google.golang.org/grpc"
)

// MasterTracker struct manages file storage using gota DataFrame
type MasterTracker struct {
	pb.UnimplementedMasterTrackerServer
	mu          sync.Mutex
	fileTable   dataframe.DataFrame
	dataKeepers map[string]time.Time // Tracks the last heartbeat timestamp
}

// NewMasterTracker initializes the Master Tracker
func NewMasterTracker() *MasterTracker {
	// Create an empty DataFrame with defined columns
	df := dataframe.New(
		series.New([]string{}, series.String, "filename"),
		series.New([]string{}, series.String, "data_keeper"),
		series.New([]string{}, series.String, "file_path"),
		series.New([]string{}, series.String, "is_alive"),
	)
	

	return &MasterTracker{
		fileTable:   df,
		dataKeepers: make(map[string]time.Time),
	}
}


// RequestUpload assigns a Data Keeper for file upload and updates DataFrame
func (s *MasterTracker) RequestUpload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dataKeeper := "localhost:50051" // pick a random datakeeper
	filePath := fmt.Sprintf("/storage/%s", req.Filename)

	// Append to DataFrame
	newRow := dataframe.New(
		series.New([]string{req.Filename}, series.String, "filename"),
		series.New([]string{dataKeeper}, series.String, "data_keeper"),
		series.New([]string{filePath}, series.String, "file_path"),
		series.New([]string{"true"}, series.String, "is_alive"),
	)
	s.fileTable = s.fileTable.RBind(newRow)

	log.Printf("[UPLOAD] File: %s assigned to Data Keeper: %s", req.Filename, dataKeeper)
	utils.PrintDataFrame(s.fileTable)
	return &pb.UploadResponse{DataKeeperAddress: dataKeeper}, nil
}

// RequestDownload returns all Data Keepers that store the requested file
func (s *MasterTracker) RequestDownload(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filtered := s.fileTable.Filter(
		dataframe.F{Colname: "filename", Comparator: "==", Comparando: req.Filename},
	)

	if filtered.Nrow() == 0 {
		log.Printf("[ERROR] File %s not found!", req.Filename)
		return nil, fmt.Errorf("file %s not found", req.Filename)
	}

	dataKeepers := filtered.Col("data_keeper").Records()
	return &pb.DownloadResponse{DataKeeperAddresses: dataKeepers}, nil
}

// SendHeartbeat updates the Data Keeper status
func (s *MasterTracker) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.dataKeepers[req.DataKeeperId] = time.Now()

	log.Printf("[HEARTBEAT] Data Keeper: %s is alive (Updated at %v)", req.DataKeeperId, s.dataKeepers[req.DataKeeperId])
	return &pb.HeartbeatResponse{Success: true}, nil
}
// CheckInactiveDataKeepers runs every 1 seconds and marks nodes as down
func (s *MasterTracker) CheckInactiveDataKeepers() {
	ticker := time.NewTicker(1 * time.Second) // Run every 1 seconds
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()

		now := time.Now()

		for dk, lastHeartbeat := range s.dataKeepers {
			if now.Sub(lastHeartbeat) > 1*time.Second { // If last heartbeat was more than 1 sec ago
				log.Printf("[WARNING] Data Keeper %s is DOWN! (Last heartbeat: %v)", dk, lastHeartbeat)

				// Update "is_alive" column to "false" for this Data Keeper
				for i := 0; i < s.fileTable.Nrow(); i++ {
					if s.fileTable.Elem(i, 1).String() == dk { // Column 1 = "data_keeper"
						s.fileTable.Elem(i, 3).Set("false") // Column 3 = "is_alive"
					}
				}
			}
		}

		s.mu.Unlock()
	}
}
func main() {
	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a single instance of MasterTracker
	masterTracker := NewMasterTracker()

	// Register it with the gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterMasterTrackerServer(grpcServer, masterTracker)

	// Start checking for inactive Data Keepers
	go masterTracker.CheckInactiveDataKeepers()

	log.Println("Master Tracker is running on port 50050 🚀")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
