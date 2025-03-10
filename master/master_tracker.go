package main

import (
	"context"
	"dfs/config"
	pb "dfs/proto"
	"dfs/utils"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"google.golang.org/grpc"
)

// MasterTracker struct manages file storage using gota DataFrame
type MasterTracker struct {
	pb.UnimplementedMasterTrackerServer
	mu               sync.Mutex
	fileTable        dataframe.DataFrame
	dataKeepersHeartbeat      map[string]time.Time // Tracks last heartbeat timestamp
	dataKeepers map[string]bool      // Tracks only alive data keepers
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
		dataKeepersHeartbeat: make(map[string]time.Time),
		dataKeepers: make(map[string]bool),
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

// SendHeartbeat updates the Data Keeper status
func (s *MasterTracker) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update last heartbeat timestamp
	s.dataKeepersHeartbeat[req.DataKeeperId] = time.Now()

	// Mark Data Keeper as alive
	s.dataKeepers[req.DataKeeperId] = true

	log.Printf("[HEARTBEAT] Data Keeper: %s is alive (Updated at %v)", req.DataKeeperId, s.dataKeepersHeartbeat[req.DataKeeperId])
	return &pb.HeartbeatResponse{Success: true}, nil
}

// CheckInactiveDataKeepers runs every 1 seconds and marks nodes as down
func (s *MasterTracker) CheckInactiveDataKeepers() {
	heartbeatTimeout := time.Duration(config.LoadConfig("config.json").DataKeeper.HeartbeatTimeout)
	ticker := time.NewTicker(heartbeatTimeout * time.Second) // Run every 1 second
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()

		now := time.Now()

		for dk, lastHeartbeat := range s.dataKeepersHeartbeat {
			if time.Duration(now.Sub(lastHeartbeat).Seconds()) >= heartbeatTimeout {
				log.Printf("[WARNING] Data Keeper %s is DOWN! (Last heartbeat: %v)", dk, lastHeartbeat)

				// Remove from aliveDataKeepers
				delete(s.dataKeepers, dk)

				// Update "is_alive" column in fileTable to "false" for this Data Keeper
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

func (s *MasterTracker) ReplicationCheck() {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
        s.mu.Lock()
        s.performReplication()
        s.mu.Unlock()
    }
}

func (s *MasterTracker) performReplication() {
	filenameSeries := s.fileTable.Col("filename")
	filenameMap := make(map[string]struct{})
	for _, filename := range filenameSeries.Records() {
		filenameMap[filename] = struct{}{}
	}
	var filenames []string
	for filename := range filenameMap {
		filenames = append(filenames, filename)
	}
	log.Printf("[REPLICATION] Replicating files names: %v", filenames)
	for _, filename := range filenames {
		filtered := s.fileTable.Filter(
			dataframe.F{Colname: "filename", Comparator: "==", Comparando: filename},
			dataframe.F{Colname: "is_alive", Comparator: "==", Comparando: "true"},
		)
		currentCount := filtered.Nrow()
		if currentCount < 3 {
			sources := filtered.Col("data_keeper").Records()
			if len(sources) == 0 {
				continue
			}
			source := sources[0]
			needed := 3 - currentCount
			for i := 0; i < needed; i++ {
				possibleDests := s.getPossibleDestinations(filtered)
				if len(possibleDests) == 0 {
					break
				}
				destination := possibleDests[0] //TODO: pick random 
				// TODO: NO REPEATED DESTINATION
				log.Printf("[REPLICATION] Replicating file: %s from Data Keeper: %s to Data Keeper: %s", filename, source, destination)
				if err := s.notifyMachineDataTransfer(source, destination, filename); err != nil {
					log.Printf("[ERROR] Replication failed for file: %s from Data Keeper: %s to Data Keeper: %s, error: %v", filename, source, destination, err)
				} else {
					log.Printf("[SUCCESS] Replication succeeded for file: %s from Data Keeper: %s to Data Keeper: %s", filename, source, destination)
				}
			}
		}
	}
}
func (s *MasterTracker) notifyMachineDataTransfer(source, destination, filename string) error {
	sourcePortInt, err := strconv.Atoi(source)
    if err != nil {
        log.Printf("Invalid source port: %s", source)
        return err
    }
    grpcPort := strconv.Itoa(sourcePortInt + 1000)
	conn, err := grpc.Dial("localhost:"+grpcPort, grpc.WithInsecure())
	if err != nil {
		log.Printf("[ERROR] Failed to connect to Data Keeper: %v", err)
		return err
	}
	dataKeeper := pb.NewDataKeeperClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	if _, err := dataKeeper.ReplicateFile(ctx, &pb.ReplicationRequest{DestinationAddress: destination, Filename: filename}); err != nil {
		log.Printf("[ERROR] Failed to replicate file: %v", err)
	}
    defer conn.Close()
	defer cancel()
    return err
}
func (s *MasterTracker) getPossibleDestinations(filtered dataframe.DataFrame) []string {
    currentDKs := filtered.Col("data_keeper").Records()
    possible := make([]string, 0)
    for dk := range s.dataKeepers {
        if !contains(currentDKs, dk) {
            possible = append(possible, dk)
        }
    }
    return possible
}
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
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
	go masterTracker.CheckInactiveDataKeepers()
	go masterTracker.ReplicationCheck()


	log.Printf("Master Tracker is running on port %d🚀", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
