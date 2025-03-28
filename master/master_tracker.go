package main

import (
	"context"
	"dfs/config"
	pb "dfs/proto"
	"dfs/utils"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"net"
	"strconv"
	"sync"
	"time"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DataKeeperInfo struct {
	Address   string           `json:"address"`
	PortsTCP  []*pb.PortStatus `json:"portsTcp"`
	PortsGRPC []*pb.PortStatus `json:"portsGrpc"`
}

// ENUM for port types
const (
	GRPC = iota // 0
	TCP         // 1
)

type PortStatus struct {
	PortNumber  int  `json:"portNumber"`
	IsAvailable bool `json:"isAvailable"`
}
type MasterTracker struct {
	pb.UnimplementedMasterTrackerServer // Embed the generated struct
	mu                                  sync.Mutex
	fileTable                           dataframe.DataFrame
	dataKeeperInfo                      map[string]DataKeeperInfo // Data Keeper name -> Data Keeper info (AVILABLE DATA KEEPER )
	dataKeepersHeartbeat                map[string]time.Time
	nodeTimers                          map[string]*time.Timer
	heartbeatTimeout                    time.Duration
}

func NewMasterTracker() *MasterTracker {
	cfg := config.LoadConfig("config.json")
	heartbeatTimeout := time.Duration(cfg.DataKeeper.HeartbeatTimeout) * time.Second + 200*time.Millisecond

	df := dataframe.New(
		series.New([]string{}, series.String, "dataKeeperName"),
		series.New([]string{}, series.String, "filename"),
		series.New([]string{}, series.String, "filePath"),
		series.New([]bool{}, series.Bool, "isAlive"),
	)

	return &MasterTracker{
		fileTable:            df,
		dataKeeperInfo:       make(map[string]DataKeeperInfo),
		dataKeepersHeartbeat: make(map[string]time.Time),
		nodeTimers:           make(map[string]*time.Timer),
		heartbeatTimeout:     heartbeatTimeout,
	}
}

func (s *MasterTracker) GetRandomAvailablePort(dataKeeperName string, portType int) (int, error) {
	dataKeeper, exists := s.dataKeeperInfo[dataKeeperName]
	if !exists {
		return 0, fmt.Errorf("dataKeeper %s not found", dataKeeperName)
	}
	var ports []*pb.PortStatus
	if portType == TCP {
		ports = dataKeeper.PortsTCP
	} else {
		ports = dataKeeper.PortsGRPC
	}
	availablePorts := []int{}
	for _, port := range ports {
		if port.IsAvailable {
			availablePorts = append(availablePorts, int(port.PortNumber))
		}
	}

	if len(availablePorts) == 0 {
		return 0, fmt.Errorf("no available ports for dataKeeper %s", dataKeeperName)
	}

	index := utils.GetRandomIndex(len(availablePorts))
	return availablePorts[index], nil
}

func (s *MasterTracker) RequestUpload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !strings.HasSuffix(strings.ToLower(req.Filename), ".mp4") {
		return nil, fmt.Errorf("only MP4 files are allowed for upload")
	}
	// Ensure we have at least one alive Data Keeper
	if len(s.dataKeeperInfo) == 0 {
		return nil, fmt.Errorf("no available data keepers for upload")
	}

	// Select the first alive Data Keeper
	var selectedDataKeeperName string
	for dk := range s.dataKeeperInfo {
		selectedDataKeeperName = dk
		break
	}

	log.Printf("[UPLOAD] File assigned to Data Keeper: %s", selectedDataKeeperName)
	selectedDataKeeperPort, _ := s.GetRandomAvailablePort(selectedDataKeeperName, TCP)
	// s.SetPortAvailability(selectedDataKeeperName, selectedDataKeeperPort, TCP, false)
	selectedDataKeeperAdress := s.dataKeeperInfo[selectedDataKeeperName].Address + ":" + strconv.Itoa(selectedDataKeeperPort)
	return &pb.UploadResponse{DataKeeperAddress: selectedDataKeeperAdress}, nil
}

func (s *MasterTracker) SetPortAvailability(dataKeeperName string, portNumber int, portType int, isAvailable bool) error {
	// Retrieve the DataKeeperInfo from the map
	dataKeeper, exists := s.dataKeeperInfo[dataKeeperName]
	if !exists {
		return fmt.Errorf("dataKeeper %s not found", dataKeeperName)
	}

	// Select the appropriate port list based on the port type
	var ports *[]*pb.PortStatus
	switch portType {
	case TCP:
		ports = &dataKeeper.PortsTCP
	case GRPC:
		ports = &dataKeeper.PortsGRPC
	default:
		return errors.New("invalid port type")
	}
	index := -1
	// Find and update the port status
	for i := range *ports {
		if (*ports)[i].PortNumber == int32(portNumber) {
			index = i
		}
	}

	if index == -1 {
		return fmt.Errorf("port %d not found in dataKeeper %s", portNumber, dataKeeperName)
	}
	switch portType {
	case TCP:
		s.dataKeeperInfo[dataKeeperName].PortsTCP[index].IsAvailable = isAvailable
	case GRPC:
		s.dataKeeperInfo[dataKeeperName].PortsGRPC[index].IsAvailable = isAvailable
	default:
		return errors.New("invalid port type")
	}
	return nil
}

func (s *MasterTracker) RequestUploadSuccess(ctx context.Context, req *pb.FileUploadSuccess) (*emptypb.Empty, error) {
	s.mu.Lock()

	s.AddFile(req.DataKeeperName, req.Filename, req.FilePath)
	// s.SetPortAvailability(req.DataKeeperName, int(req.PortNumber), TCP, true)
	fmt.Printf("[UPLOAD SUCCESS] Selected port: %v\n", s.dataKeeperInfo)
	s.performReplication()

	s.mu.Unlock()

	notifyClient(true,req.ClientAddress);

	return &emptypb.Empty{}, nil
}
func (s *MasterTracker) AddFile(dataKeeperName, filename, filePath string) {
	newRow := dataframe.New(
		series.New([]string{dataKeeperName}, series.String, "dataKeeperName"),
		series.New([]string{filename}, series.String, "filename"),
		series.New([]string{filePath}, series.String, "filePath"),
		series.New([]bool{true}, series.Bool, "isAlive"),
	)

	s.fileTable = s.fileTable.RBind(newRow)

	log.Printf("[UPLOAD SUCCESS] File: %s uploaded successfully to Data Keeper: %s", filename, dataKeeperName)
	utils.PrintDataFrame(s.fileTable)
}

// RequestDownload returns all Data Keepers that store the requested file
func (s *MasterTracker) RequestDownload(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Filter file table to get all Data Keepers storing the requested file

	filtered := s.fileTable.FilterAggregation(
		dataframe.And,
		dataframe.F{Colname: "filename", Comparator: series.Eq, Comparando: req.Filename},
		dataframe.F{Colname: "isAlive", Comparator: series.Eq, Comparando: true}, // Ensure Data Keeper is alive
	)

	if filtered.Nrow() == 0 {
		log.Printf("[ERROR] File %s not found or no active Data Keeper!", req.Filename)
		return nil, fmt.Errorf("file %s not found or no active Data Keeper", req.Filename)
	}

	dataKeepers := filtered.Col("dataKeeperName").Records()
	log.Printf("Data Keepers: %v %v", dataKeepers,s.dataKeeperInfo)
	dataKeepersAddresses := s.FormatNodeAdresses(dataKeepers)
	log.Printf("[DOWNLOAD] File %s found on Data Keepers: %v", req.Filename, dataKeepersAddresses)
	return &pb.DownloadResponse{DataKeeperAddresses: dataKeepersAddresses}, nil
}

func (s *MasterTracker) FormatNodeAdresses(dataKeepers []string) ([]string) {
	dataKeeperAddresses := make ([]string, 0)
	for _, dataKeeper := range dataKeepers {
		address := s.dataKeeperInfo[dataKeeper].Address
		for _, port := range s.dataKeeperInfo[dataKeeper].PortsTCP {
			if port.IsAvailable {
				dataKeeperAddresses = append(dataKeeperAddresses, address + ":" + strconv.Itoa(int(port.PortNumber)))
				break
			}
		}
	}
	return dataKeeperAddresses
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
	delete(s.dataKeeperInfo, dk)

	// Update "isAlive" column in fileTable to "false" for this Data Keeper.
	for i := 0; i < s.fileTable.Nrow(); i++ {
		if s.fileTable.Elem(i, 0).String() == dk { // Column 0 = "dataKeeperName"
			s.fileTable.Elem(i, 3).Set(false) // Column 3 = "isAlive"
		}
	}
	log.Printf("[INFO] Updated fileTable after marking Data Keeper %s as DOWN", s.fileTable)
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
	s.dataKeepersHeartbeat[req.DataKeeperName] = now

	// Stop any existing timer for this Data Keeper.
	if timer, exists := s.nodeTimers[req.DataKeeperName]; exists {
		timer.Stop()
	}

	// Define a small buffer (e.g., 100ms) to account for jitter.
	buffer := 100 * time.Millisecond
	for i := range s.fileTable.Nrow() {
		if s.fileTable.Elem(i, 0).String() == req.DataKeeperName { // Column 0 = "dataKeeperName"
			s.fileTable.Elem(i, 3).Set(true) // Column 3 = "isAlive"
		}
	}
	// Reset the timer for this Data Keeper to fire after the heartbeat timeout plus the buffer.
	s.nodeTimers[req.DataKeeperName] = time.AfterFunc(s.heartbeatTimeout+buffer, func() {
		s.markDataKeeperDown(req.DataKeeperName)
	})
	if _, exists := s.dataKeeperInfo[req.DataKeeperName]; !exists {
		s.dataKeeperInfo[req.DataKeeperName] = DataKeeperInfo{Address: req.DataKeeperAddress, PortsTCP: req.PortsTCP, PortsGRPC: req.PortsGRPC}
	}

	log.Printf("[HEARTBEAT] Data Keeper: %s is alive (Updated at %v)", req.DataKeeperName, now)
	return &pb.HeartbeatResponse{Success: true}, nil
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
	// Extract all unique filenames from the fileTable
	filenameSeries := s.fileTable.Col("filename")
	filenameMap := make(map[string]struct{})
	// Use a map to ensure uniqueness
	for _, filename := range filenameSeries.Records() {
		filenameMap[filename] = struct{}{}
	}

	// Convert the map keys to a slice
	var filenames []string
	for filename := range filenameMap {
		filenames = append(filenames, filename)
	}
	log.Printf("[REPLICATION] Replicating files names: %v", filenames)
	for _, filename := range filenames {
		filtered := s.fileTable.FilterAggregation(
			dataframe.And,
			dataframe.F{Colname: "filename", Comparator: series.Eq, Comparando: filename},
			dataframe.F{Colname: "isAlive", Comparator: series.Eq, Comparando: true},
		)
		currentCount := filtered.Nrow()
		if currentCount < 3 {
			sources := filtered.Col("dataKeeperName").Records()
			if (len(sources) == 0) {
				continue
			}
			source := sources[0]
			remainingDataKeepers := 3 - currentCount
			possibleDests := s.getPossibleDestinations(filtered)

			if len(possibleDests) == 0 {
				continue
			}
			
			if len(possibleDests) < remainingDataKeepers {
				remainingDataKeepers = len(possibleDests) // Prevent out-of-bounds errors
			}

			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(possibleDests), func(i, j int) {
				possibleDests[i], possibleDests[j] = possibleDests[j], possibleDests[i]
			})

			selectedDests := possibleDests[:remainingDataKeepers]

			for _, destination := range selectedDests {
				log.Printf("[REPLICATION] Replicating file: %s from Data Keeper: %s to Data Keeper: %s", filename, source, destination)

				if err := s.notifyMachineDataTransfer(source, destination, filename); err == nil {
					log.Printf("[SUCCESS] Replication succeeded for file: %s from Data Keeper: %s to Data Keeper: %s", filename, source, destination)
				} else {
					log.Printf("[ERROR] Replication failed for file: %s from Data Keeper: %s to Data Keeper: %s, error: %v", filename, source, destination, err)
				}
			}
		}
	}
}
func (s *MasterTracker) notifyMachineDataTransfer(sourceNodeName, destinationNodeName, filename string) error {
	grpcPortSrc, _ := s.GetRandomAvailablePort(sourceNodeName, GRPC)
	tcpPortDest, err := s.GetRandomAvailablePort(destinationNodeName, TCP)

	if err != nil {
		log.Printf("Invalid source port: %s", sourceNodeName)
		return err
	}
	conn, err := grpc.Dial(s.dataKeeperInfo[sourceNodeName].Address+":"+strconv.Itoa(grpcPortSrc), grpc.WithInsecure())
	if err != nil {
		log.Printf("[ERROR] Failed to connect to Data Keeper: %v", err)
		return err
	}
	dataKeeper := pb.NewDataKeeperClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	destinationAddress := s.dataKeeperInfo[destinationNodeName].Address + ":" + strconv.Itoa(tcpPortDest)
	response, err := dataKeeper.ReplicateFile(ctx, &pb.ReplicationRequest{DestinationAddress: destinationAddress, Filename: filename,DestinationName: destinationNodeName});
	if err != nil {
		log.Printf("[ERROR] Failed to replicate file: %v", err)
	}
	s.AddFile(response.DataKeeperName, response.Filename, response.FilePath)
	// s.SetPortAvailability(response.DataKeeperName, int(response.PortNumber), TCP, true)
	defer conn.Close()
	defer cancel()
	return err
}
func (s *MasterTracker) getPossibleDestinations(filtered dataframe.DataFrame) []string {
	currentDKs := filtered.Col("dataKeeperName").Records()
	possible := make([]string, 0)
	for dk := range s.dataKeeperInfo {
		if !contains(currentDKs, dk) {
			possible = append(possible, dk)
		}
	}
	return possible
}

func main() {
	port := config.LoadConfig("config.json").Server.Port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	masterTracker := NewMasterTracker()
	grpcServer := grpc.NewServer()
	pb.RegisterMasterTrackerServer(grpcServer, masterTracker)
	go masterTracker.ReplicationCheck()

	log.Printf("Master Tracker is running on port %d🚀", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func notifyClient(isUploaded bool, address string) {
	// Notify client
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to client: %v", err)
	}

	client := pb.NewClientClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	if _, err := client.NotifyUploadCompletion(ctx, &pb.UploadSuccessResponse{Success: isUploaded}); err != nil {
		log.Printf("[ERROR] Failed to notify client: %v", err)
	}
	cancel()
	conn.Close()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}