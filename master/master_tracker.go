package main

import (
	"context"
	pb "dfs/proto"
	"fmt"
	"log"

	//"math/rand"
	"net"
	"sync"
	"time"

	"dfs/config"
	"dfs/utils"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"google.golang.org/grpc"
)

type MasterTracker struct {
	pb.UnimplementedMasterTrackerServer
	mu                   sync.Mutex
	fileTable            dataframe.DataFrame
	dataKeepersHeartbeat map[string]time.Time   // Tracks last heartbeat timestamp
	dataKeepers          map[string]bool        // Tracks alive data keepers
	nodeTimers           map[string]*time.Timer // Per-node timers
	heartbeatTimeout     time.Duration          // Configured heartbeat timeout
}

func NewMasterTracker() *MasterTracker {
	cfg := config.LoadConfig("config.json")
	heartbeatTimeout := time.Duration(cfg.DataKeeper.HeartbeatTimeout) * time.Second

	df := dataframe.New(
		series.New([]string{}, series.String, "filename"),
		series.New([]string{}, series.String, "data_keeper"),
		series.New([]string{}, series.String, "file_path"),
		series.New([]string{}, series.String, "is_alive"),
	)

	return &MasterTracker{
		fileTable:            df,
		dataKeepersHeartbeat: make(map[string]time.Time),
		dataKeepers:          make(map[string]bool),
		nodeTimers:           make(map[string]*time.Timer),
		heartbeatTimeout:     heartbeatTimeout,
	}
}

// RequestUpload assigns a Data Keeper for file upload and updates DataFrame
func (s *MasterTracker) RequestUpload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure we have at least one alive Data Keeper
	if len(s.dataKeepers) == 0 {
		return nil, fmt.Errorf("no available data keepers for upload")
	}

	// Select the first alive Data Keeper
	var selectedDataKeeper string
	for dk := range s.dataKeepers {
		selectedDataKeeper = dk
		break
	}

	filePath := fmt.Sprintf("/storage/%s", req.Filename)

	// Append to DataFrame
	newRow := dataframe.New(
		series.New([]string{req.Filename}, series.String, "filename"),
		series.New([]string{selectedDataKeeper}, series.String, "data_keeper"),
		series.New([]string{filePath}, series.String, "file_path"),
		series.New([]string{"true"}, series.String, "is_alive"),
	)
	s.fileTable = s.fileTable.RBind(newRow)

	log.Printf("[UPLOAD] File: %s assigned to Data Keeper: %s", req.Filename, selectedDataKeeper)
	utils.PrintDataFrame(s.fileTable)
	return &pb.UploadResponse{DataKeeperAddress: selectedDataKeeper}, nil
}

// RequestDownload returns all Data Keepers that store the requested file
func (s *MasterTracker) RequestDownload(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Filter file table to get all Data Keepers storing the requested file
	filtered := s.fileTable.Filter(
		dataframe.F{Colname: "filename", Comparator: "==", Comparando: req.Filename},
		dataframe.F{Colname: "is_alive", Comparator: "==", Comparando: "true"}, // Ensure Data Keeper is alive
	)

	if filtered.Nrow() == 0 {
		log.Printf("[ERROR] File %s not found or no active Data Keeper!", req.Filename)
		return nil, fmt.Errorf("file %s not found or no active Data Keeper", req.Filename)
	}

	dataKeepers := filtered.Col("data_keeper").Records()
	return &pb.DownloadResponse{DataKeeperAddresses: dataKeepers}, nil
}

// marks the given Data Keeper as down if its last heartbeat is older than the timeout.
func (s *MasterTracker) markDataKeeperDown(dk string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	lastHeartbeat, exists := s.dataKeepersHeartbeat[dk]
	if !exists {
		return
	}

	// Check if the elapsed time is less than the timeout.
	if time.Since(lastHeartbeat) < s.heartbeatTimeout {
		return
	}

	log.Printf("[WARNING] Data Keeper %s is DOWN! (Last heartbeat: %v)", dk, lastHeartbeat)

	// Remove the Data Keeper from the alive list.
	delete(s.dataKeepers, dk)

	// Update "is_alive" column in fileTable to "false" for this Data Keeper.
	for i := 0; i < s.fileTable.Nrow(); i++ {
		if s.fileTable.Elem(i, 1).String() == dk { // Column 1 = "data_keeper"
			s.fileTable.Elem(i, 3).Set("false") // Column 3 = "is_alive"
		}
	}

	// Remove and stop the timer for this Data Keeper.
	if timer, exists := s.nodeTimers[dk]; exists {
		timer.Stop()
		delete(s.nodeTimers, dk)
	}
}

// updates the Data Keeper status and resets its timer.
func (s *MasterTracker) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.dataKeepersHeartbeat[req.DataKeeperId] = now
	s.dataKeepers[req.DataKeeperId] = true

	// Stop any existing timer for this Data Keeper.
	if timer, exists := s.nodeTimers[req.DataKeeperId]; exists {
		timer.Stop()
	}

	// Define a small buffer (e.g., 100ms) to account for jitter.
	buffer := 100 * time.Millisecond

	// Reset the timer for this Data Keeper to fire after the heartbeat timeout plus the buffer.
	s.nodeTimers[req.DataKeeperId] = time.AfterFunc(s.heartbeatTimeout+buffer, func() {
		s.markDataKeeperDown(req.DataKeeperId)
	})

	log.Printf("[HEARTBEAT] Data Keeper: %s is alive (Updated at %v)", req.DataKeeperId, now)
	return &pb.HeartbeatResponse{Success: true}, nil
}

func main() {
	port := config.LoadConfig("config.json").Server.Port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a single instance of MasterTracker
	masterTracker := NewMasterTracker()

	// Register it with the gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterMasterTrackerServer(grpcServer, masterTracker)

	// Start checking for inactive Data Keepers
	// go masterTracker.CheckInactiveDataKeepers()

	log.Printf("Master Tracker is running on port %d🚀", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
