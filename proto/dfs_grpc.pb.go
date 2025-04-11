// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.0
// source: dfs.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MasterTracker_RequestUpload_FullMethodName        = "/dfs.MasterTracker/RequestUpload"
	MasterTracker_RequestUploadSuccess_FullMethodName = "/dfs.MasterTracker/RequestUploadSuccess"
	MasterTracker_RequestDownload_FullMethodName      = "/dfs.MasterTracker/RequestDownload"
	MasterTracker_SendHeartbeat_FullMethodName        = "/dfs.MasterTracker/SendHeartbeat"
	MasterTracker_RegisterPortStatus_FullMethodName   = "/dfs.MasterTracker/RegisterPortStatus"
)

// MasterTrackerClient is the client API for MasterTracker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Master Tracker Service: Handles file upload/download requests and heartbeat updates.
type MasterTrackerClient interface {
	// A client requests a Data Keeper to upload a file.
	RequestUpload(ctx context.Context, in *UploadRequest, opts ...grpc.CallOption) (*UploadResponse, error)
	// Data Keeper notifies the Master Tracker that a file has been uploaded successfully.
	RequestUploadSuccess(ctx context.Context, in *FileUploadSuccess, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// A client requests addresses to download a file.
	RequestDownload(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadResponse, error)
	// Data Keeper sends a heartbeat to indicate it is alive.
	SendHeartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
	// Data Keeper registers its port status with the Master Tracker.
	RegisterPortStatus(ctx context.Context, in *PortRegistrationRequest, opts ...grpc.CallOption) (*PortRegistrationResponse, error)
}

type masterTrackerClient struct {
	cc grpc.ClientConnInterface
}

func NewMasterTrackerClient(cc grpc.ClientConnInterface) MasterTrackerClient {
	return &masterTrackerClient{cc}
}

func (c *masterTrackerClient) RequestUpload(ctx context.Context, in *UploadRequest, opts ...grpc.CallOption) (*UploadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadResponse)
	err := c.cc.Invoke(ctx, MasterTracker_RequestUpload_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterTrackerClient) RequestUploadSuccess(ctx context.Context, in *FileUploadSuccess, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, MasterTracker_RequestUploadSuccess_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterTrackerClient) RequestDownload(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DownloadResponse)
	err := c.cc.Invoke(ctx, MasterTracker_RequestDownload_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterTrackerClient) SendHeartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, MasterTracker_SendHeartbeat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterTrackerClient) RegisterPortStatus(ctx context.Context, in *PortRegistrationRequest, opts ...grpc.CallOption) (*PortRegistrationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PortRegistrationResponse)
	err := c.cc.Invoke(ctx, MasterTracker_RegisterPortStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MasterTrackerServer is the server API for MasterTracker service.
// All implementations must embed UnimplementedMasterTrackerServer
// for forward compatibility.
//
// Master Tracker Service: Handles file upload/download requests and heartbeat updates.
type MasterTrackerServer interface {
	// A client requests a Data Keeper to upload a file.
	RequestUpload(context.Context, *UploadRequest) (*UploadResponse, error)
	// Data Keeper notifies the Master Tracker that a file has been uploaded successfully.
	RequestUploadSuccess(context.Context, *FileUploadSuccess) (*emptypb.Empty, error)
	// A client requests addresses to download a file.
	RequestDownload(context.Context, *DownloadRequest) (*DownloadResponse, error)
	// Data Keeper sends a heartbeat to indicate it is alive.
	SendHeartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
	// Data Keeper registers its port status with the Master Tracker.
	RegisterPortStatus(context.Context, *PortRegistrationRequest) (*PortRegistrationResponse, error)
	mustEmbedUnimplementedMasterTrackerServer()
}

// UnimplementedMasterTrackerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMasterTrackerServer struct{}

func (UnimplementedMasterTrackerServer) RequestUpload(context.Context, *UploadRequest) (*UploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestUpload not implemented")
}
func (UnimplementedMasterTrackerServer) RequestUploadSuccess(context.Context, *FileUploadSuccess) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestUploadSuccess not implemented")
}
func (UnimplementedMasterTrackerServer) RequestDownload(context.Context, *DownloadRequest) (*DownloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestDownload not implemented")
}
func (UnimplementedMasterTrackerServer) SendHeartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendHeartbeat not implemented")
}
func (UnimplementedMasterTrackerServer) RegisterPortStatus(context.Context, *PortRegistrationRequest) (*PortRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterPortStatus not implemented")
}
func (UnimplementedMasterTrackerServer) mustEmbedUnimplementedMasterTrackerServer() {}
func (UnimplementedMasterTrackerServer) testEmbeddedByValue()                       {}

// UnsafeMasterTrackerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MasterTrackerServer will
// result in compilation errors.
type UnsafeMasterTrackerServer interface {
	mustEmbedUnimplementedMasterTrackerServer()
}

func RegisterMasterTrackerServer(s grpc.ServiceRegistrar, srv MasterTrackerServer) {
	// If the following call pancis, it indicates UnimplementedMasterTrackerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MasterTracker_ServiceDesc, srv)
}

func _MasterTracker_RequestUpload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterTrackerServer).RequestUpload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterTracker_RequestUpload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterTrackerServer).RequestUpload(ctx, req.(*UploadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterTracker_RequestUploadSuccess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileUploadSuccess)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterTrackerServer).RequestUploadSuccess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterTracker_RequestUploadSuccess_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterTrackerServer).RequestUploadSuccess(ctx, req.(*FileUploadSuccess))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterTracker_RequestDownload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterTrackerServer).RequestDownload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterTracker_RequestDownload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterTrackerServer).RequestDownload(ctx, req.(*DownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterTracker_SendHeartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterTrackerServer).SendHeartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterTracker_SendHeartbeat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterTrackerServer).SendHeartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterTracker_RegisterPortStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PortRegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterTrackerServer).RegisterPortStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterTracker_RegisterPortStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterTrackerServer).RegisterPortStatus(ctx, req.(*PortRegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MasterTracker_ServiceDesc is the grpc.ServiceDesc for MasterTracker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MasterTracker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dfs.MasterTracker",
	HandlerType: (*MasterTrackerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestUpload",
			Handler:    _MasterTracker_RequestUpload_Handler,
		},
		{
			MethodName: "RequestUploadSuccess",
			Handler:    _MasterTracker_RequestUploadSuccess_Handler,
		},
		{
			MethodName: "RequestDownload",
			Handler:    _MasterTracker_RequestDownload_Handler,
		},
		{
			MethodName: "SendHeartbeat",
			Handler:    _MasterTracker_SendHeartbeat_Handler,
		},
		{
			MethodName: "RegisterPortStatus",
			Handler:    _MasterTracker_RegisterPortStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dfs.proto",
}

const (
	DataKeeper_ReplicateFile_FullMethodName = "/dfs.DataKeeper/ReplicateFile"
)

// DataKeeperClient is the client API for DataKeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataKeeperClient interface {
	ReplicateFile(ctx context.Context, in *ReplicationRequest, opts ...grpc.CallOption) (*FileUploadSuccess, error)
}

type dataKeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewDataKeeperClient(cc grpc.ClientConnInterface) DataKeeperClient {
	return &dataKeeperClient{cc}
}

func (c *dataKeeperClient) ReplicateFile(ctx context.Context, in *ReplicationRequest, opts ...grpc.CallOption) (*FileUploadSuccess, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FileUploadSuccess)
	err := c.cc.Invoke(ctx, DataKeeper_ReplicateFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataKeeperServer is the server API for DataKeeper service.
// All implementations must embed UnimplementedDataKeeperServer
// for forward compatibility.
type DataKeeperServer interface {
	ReplicateFile(context.Context, *ReplicationRequest) (*FileUploadSuccess, error)
	mustEmbedUnimplementedDataKeeperServer()
}

// UnimplementedDataKeeperServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDataKeeperServer struct{}

func (UnimplementedDataKeeperServer) ReplicateFile(context.Context, *ReplicationRequest) (*FileUploadSuccess, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplicateFile not implemented")
}
func (UnimplementedDataKeeperServer) mustEmbedUnimplementedDataKeeperServer() {}
func (UnimplementedDataKeeperServer) testEmbeddedByValue()                    {}

// UnsafeDataKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataKeeperServer will
// result in compilation errors.
type UnsafeDataKeeperServer interface {
	mustEmbedUnimplementedDataKeeperServer()
}

func RegisterDataKeeperServer(s grpc.ServiceRegistrar, srv DataKeeperServer) {
	// If the following call pancis, it indicates UnimplementedDataKeeperServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DataKeeper_ServiceDesc, srv)
}

func _DataKeeper_ReplicateFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplicationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataKeeperServer).ReplicateFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataKeeper_ReplicateFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataKeeperServer).ReplicateFile(ctx, req.(*ReplicationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DataKeeper_ServiceDesc is the grpc.ServiceDesc for DataKeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataKeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dfs.DataKeeper",
	HandlerType: (*DataKeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReplicateFile",
			Handler:    _DataKeeper_ReplicateFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dfs.proto",
}

const (
	Client_NotifyUploadCompletion_FullMethodName = "/dfs.Client/NotifyUploadCompletion"
)

// ClientClient is the client API for Client service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientClient interface {
	NotifyUploadCompletion(ctx context.Context, in *UploadSuccessResponse, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type clientClient struct {
	cc grpc.ClientConnInterface
}

func NewClientClient(cc grpc.ClientConnInterface) ClientClient {
	return &clientClient{cc}
}

func (c *clientClient) NotifyUploadCompletion(ctx context.Context, in *UploadSuccessResponse, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Client_NotifyUploadCompletion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientServer is the server API for Client service.
// All implementations must embed UnimplementedClientServer
// for forward compatibility.
type ClientServer interface {
	NotifyUploadCompletion(context.Context, *UploadSuccessResponse) (*emptypb.Empty, error)
	mustEmbedUnimplementedClientServer()
}

// UnimplementedClientServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedClientServer struct{}

func (UnimplementedClientServer) NotifyUploadCompletion(context.Context, *UploadSuccessResponse) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyUploadCompletion not implemented")
}
func (UnimplementedClientServer) mustEmbedUnimplementedClientServer() {}
func (UnimplementedClientServer) testEmbeddedByValue()                {}

// UnsafeClientServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientServer will
// result in compilation errors.
type UnsafeClientServer interface {
	mustEmbedUnimplementedClientServer()
}

func RegisterClientServer(s grpc.ServiceRegistrar, srv ClientServer) {
	// If the following call pancis, it indicates UnimplementedClientServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Client_ServiceDesc, srv)
}

func _Client_NotifyUploadCompletion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadSuccessResponse)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).NotifyUploadCompletion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Client_NotifyUploadCompletion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).NotifyUploadCompletion(ctx, req.(*UploadSuccessResponse))
	}
	return interceptor(ctx, in, info, handler)
}

// Client_ServiceDesc is the grpc.ServiceDesc for Client service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Client_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dfs.Client",
	HandlerType: (*ClientServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NotifyUploadCompletion",
			Handler:    _Client_NotifyUploadCompletion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dfs.proto",
}
