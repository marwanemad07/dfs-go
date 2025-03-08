package main

import (
	"context"
	pb "dfs/proto"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"google.golang.org/grpc"
)

// MasterTracker struct manages file storage using gota DataFrame
type MasterTracker struct {
	pb.UnimplementedMasterTrackerServer
	mu          sync.Mutex
	fileTable   dataframe.DataFrame
	dataKeepers map[string]bool // Tracks which Data Keepers are alive
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
		dataKeepers: make(map[string]bool),
	}
}
func printDataFrame(df dataframe.DataFrame) {
	records := df.Records() // Get all rows as [][]string

	if len(records) == 0 {
		fmt.Println("Empty DataFrame")
		return
	}

	// Print header
	header := records[0]
	fmt.Printf("| %-40s | %-20s | %-30s | %-10s |\n", header[0], header[1], header[2], header[3])
	fmt.Println(strings.Repeat("-", 110))

	// Print rows
	for _, row := range records[1:] { // Skip header row
		fmt.Printf("| %-40s | %-20s | %-30s | %-10s |\n", row[0], row[1], row[2], row[3])
	}
	fmt.Println(strings.Repeat("-", 110))
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
	printDataFrame(s.fileTable)
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

	s.dataKeepers[req.DataKeeperId] = true
	log.Printf("[HEARTBEAT] Data Keeper: %s is alive", req.DataKeeperId)

	return &pb.HeartbeatResponse{Success: true}, nil
}

// CheckInactiveDataKeepers runs every 10 seconds and marks nodes as down
func (s *MasterTracker) CheckInactiveDataKeepers() {
	for {
		time.Sleep(10 * time.Second)

		s.mu.Lock()
		for dk, alive := range s.dataKeepers {
			if !alive {
				log.Printf("[WARNING] Data Keeper %s is DOWN!", dk)
			}
			s.dataKeepers[dk] = false // Reset for next check
		}
		s.mu.Unlock()
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMasterTrackerServer(grpcServer, NewMasterTracker())

	// Start checking for inactive Data Keepers
	go func(s *MasterTracker) {
		s.CheckInactiveDataKeepers()
	}(NewMasterTracker())

	log.Println("Master Tracker is running on port 50050 🚀")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
