// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v6.30.0
// source: dfs.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReplicationRequest struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	Filename           string                 `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	DestinationAddress string                 `protobuf:"bytes,2,opt,name=destinationAddress,proto3" json:"destinationAddress,omitempty"`
	DestinationName    string                 `protobuf:"bytes,3,opt,name=destinationName,proto3" json:"destinationName,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *ReplicationRequest) Reset() {
	*x = ReplicationRequest{}
	mi := &file_dfs_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReplicationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplicationRequest) ProtoMessage() {}

func (x *ReplicationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dfs_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplicationRequest.ProtoReflect.Descriptor instead.
func (*ReplicationRequest) Descriptor() ([]byte, []int) {
	return file_dfs_proto_rawDescGZIP(), []int{0}
}

func (x *ReplicationRequest) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *ReplicationRequest) GetDestinationAddress() string {
	if x != nil {
		return x.DestinationAddress
	}
	return ""
}

func (x *ReplicationRequest) GetDestinationName() string {
	if x != nil {
		return x.DestinationName
	}
	return ""
}

// schemas for file upload/download and heartbeat.
// Message definitions for file upload/download.
type UploadRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Filename      string                 `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadRequest) Reset() {
	*x = UploadRequest{}
	mi := &file_dfs_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadRequest) ProtoMessage() {}

func (x *UploadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dfs_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadRequest.ProtoReflect.Descriptor instead.
func (*UploadRequest) Descriptor() ([]byte, []int) {
	return file_dfs_proto_rawDescGZIP(), []int{1}
}

func (x *UploadRequest) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

type FileUploadSuccess struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	DataKeeperName string                 `protobuf:"bytes,1,opt,name=dataKeeperName,proto3" json:"dataKeeperName,omitempty"`
	Filename       string                 `protobuf:"bytes,2,opt,name=filename,proto3" json:"filename,omitempty"`
	FilePath       string                 `protobuf:"bytes,3,opt,name=filePath,proto3" json:"filePath,omitempty"`
	PortNumber     int32                  `protobuf:"varint,4,opt,name=portNumber,proto3" json:"portNumber,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *FileUploadSuccess) Reset() {
	*x = FileUploadSuccess{}
	mi := &file_dfs_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileUploadSuccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadSuccess) ProtoMessage() {}

func (x *FileUploadSuccess) ProtoReflect() protoreflect.Message {
	mi := &file_dfs_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadSuccess.ProtoReflect.Descriptor instead.
func (*FileUploadSuccess) Descriptor() ([]byte, []int) {
	return file_dfs_proto_rawDescGZIP(), []int{2}
}

func (x *FileUploadSuccess) GetDataKeeperName() string {
	if x != nil {
		return x.DataKeeperName
	}
	return ""
}

func (x *FileUploadSuccess) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *FileUploadSuccess) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *FileUploadSuccess) GetPortNumber() int32 {
	if x != nil {
		return x.PortNumber
	}
	return 0
}

type UploadSuccessResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadSuccessResponse) Reset() {
	*x = UploadSuccessResponse{}
	mi := &file_dfs_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadSuccessResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadSuccessResponse) ProtoMessage() {}

func (x *UploadSuccessResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dfs_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadSuccessResponse.ProtoReflect.Descriptor instead.
func (*UploadSuccessResponse) Descriptor() ([]byte, []int) {
	return file_dfs_proto_rawDescGZIP(), []int{3}
}

func (x *UploadSuccessResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type UploadResponse struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	DataKeeperAddress string                 `protobuf:"bytes,1,opt,name=dataKeeperAddress,proto3" json:"dataKeeperAddress,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *UploadResponse) Reset() {
	*x = UploadResponse{}
	mi := &file_dfs_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadResponse) ProtoMessage() {}

func (x *UploadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dfs_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadResponse.ProtoReflect.Descriptor instead.
func (*UploadResponse) Descriptor() ([]byte, []int) {
	return file_dfs_proto_rawDescGZIP(), []int{4}
}

func (x *UploadResponse) GetDataKeeperAddress() string {
	if x != nil {
		return x.DataKeeperAddress
	}
	return ""
}

type DownloadRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Filename      string                 `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DownloadRequest) Reset() {
	*x = DownloadRequest{}
	mi := &file_dfs_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DownloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadRequest) ProtoMessage() {}

func (x *DownloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dfs_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadRequest.ProtoReflect.Descriptor instead.
func (*DownloadRequest) Descriptor() ([]byte, []int) {
	return file_dfs_proto_rawDescGZIP(), []int{5}
}

func (x *DownloadRequest) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

type DownloadResponse struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	DataKeeperAddresses []string               `protobuf:"bytes,1,rep,name=dataKeeperAddresses,proto3" json:"dataKeeperAddresses,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *DownloadResponse) Reset() {
	*x = DownloadResponse{}
	mi := &file_dfs_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DownloadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadResponse) ProtoMessage() {}

func (x *DownloadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dfs_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadResponse.ProtoReflect.Descriptor instead.
func (*DownloadResponse) Descriptor() ([]byte, []int) {
	return file_dfs_proto_rawDescGZIP(), []int{6}
}

func (x *DownloadResponse) GetDataKeeperAddresses() []string {
	if x != nil {
		return x.DataKeeperAddresses
	}
	return nil
}

// Message definitions for heartbeat.
type HeartbeatRequest struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	DataKeeperName    string                 `protobuf:"bytes,1,opt,name=dataKeeperName,proto3" json:"dataKeeperName,omitempty"`
	DataKeeperAddress string                 `protobuf:"bytes,2,opt,name=dataKeeperAddress,proto3" json:"dataKeeperAddress,omitempty"`
	PortsTCP          []*PortStatus          `protobuf:"bytes,3,rep,name=portsTCP,proto3" json:"portsTCP,omitempty"`
	PortsGRPC         []*PortStatus          `protobuf:"bytes,4,rep,name=portsGRPC,proto3" json:"portsGRPC,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *HeartbeatRequest) Reset() {
	*x = HeartbeatRequest{}
	mi := &file_dfs_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HeartbeatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatRequest) ProtoMessage() {}

func (x *HeartbeatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dfs_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatRequest.ProtoReflect.Descriptor instead.
func (*HeartbeatRequest) Descriptor() ([]byte, []int) {
	return file_dfs_proto_rawDescGZIP(), []int{7}
}

func (x *HeartbeatRequest) GetDataKeeperName() string {
	if x != nil {
		return x.DataKeeperName
	}
	return ""
}

func (x *HeartbeatRequest) GetDataKeeperAddress() string {
	if x != nil {
		return x.DataKeeperAddress
	}
	return ""
}

func (x *HeartbeatRequest) GetPortsTCP() []*PortStatus {
	if x != nil {
		return x.PortsTCP
	}
	return nil
}

func (x *HeartbeatRequest) GetPortsGRPC() []*PortStatus {
	if x != nil {
		return x.PortsGRPC
	}
	return nil
}

type HeartbeatResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HeartbeatResponse) Reset() {
	*x = HeartbeatResponse{}
	mi := &file_dfs_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HeartbeatResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatResponse) ProtoMessage() {}

func (x *HeartbeatResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dfs_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatResponse.ProtoReflect.Descriptor instead.
func (*HeartbeatResponse) Descriptor() ([]byte, []int) {
	return file_dfs_proto_rawDescGZIP(), []int{8}
}

func (x *HeartbeatResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type PortStatus struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PortNumber    int32                  `protobuf:"varint,1,opt,name=portNumber,proto3" json:"portNumber,omitempty"`
	IsAvailable   bool                   `protobuf:"varint,2,opt,name=isAvailable,proto3" json:"isAvailable,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PortStatus) Reset() {
	*x = PortStatus{}
	mi := &file_dfs_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PortStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortStatus) ProtoMessage() {}

func (x *PortStatus) ProtoReflect() protoreflect.Message {
	mi := &file_dfs_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortStatus.ProtoReflect.Descriptor instead.
func (*PortStatus) Descriptor() ([]byte, []int) {
	return file_dfs_proto_rawDescGZIP(), []int{9}
}

func (x *PortStatus) GetPortNumber() int32 {
	if x != nil {
		return x.PortNumber
	}
	return 0
}

func (x *PortStatus) GetIsAvailable() bool {
	if x != nil {
		return x.IsAvailable
	}
	return false
}

var File_dfs_proto protoreflect.FileDescriptor

var file_dfs_proto_rawDesc = string([]byte{
	0x0a, 0x09, 0x64, 0x66, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x64, 0x66, 0x73,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8a, 0x01,
	0x0a, 0x12, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x2e, 0x0a, 0x12, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x64, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x28, 0x0a, 0x0f, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x2b, 0x0a, 0x0d, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x93, 0x01, 0x0a, 0x11, 0x46, 0x69, 0x6c, 0x65,
	0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x26, 0x0a,
	0x0e, 0x64, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x64, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x65, 0x70, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1e, 0x0a,
	0x0a, 0x70, 0x6f, 0x72, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0a, 0x70, 0x6f, 0x72, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x31, 0x0a,
	0x15, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x22, 0x3e, 0x0a, 0x0e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2c, 0x0a, 0x11, 0x64, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x64,
	0x61, 0x74, 0x61, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x22, 0x2d, 0x0a, 0x0f, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x44, 0x0a, 0x10, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x13, 0x64, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x65, 0x70, 0x65,
	0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x13, 0x64, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x65, 0x73, 0x22, 0xc4, 0x01, 0x0a, 0x10, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62,
	0x65, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x64, 0x61,
	0x74, 0x61, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x64, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x2c, 0x0a, 0x11, 0x64, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x64,
	0x61, 0x74, 0x61, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x2b, 0x0a, 0x08, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x54, 0x43, 0x50, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x66, 0x73, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x08, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x54, 0x43, 0x50, 0x12, 0x2d, 0x0a,
	0x09, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x47, 0x52, 0x50, 0x43, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x64, 0x66, 0x73, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x09, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x47, 0x52, 0x50, 0x43, 0x22, 0x2d, 0x0a, 0x11,
	0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x4e, 0x0a, 0x0a, 0x50,
	0x6f, 0x72, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x6f, 0x72,
	0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70,
	0x6f, 0x72, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x69, 0x73, 0x41,
	0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b,
	0x69, 0x73, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x32, 0x91, 0x02, 0x0a, 0x0d,
	0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x38, 0x0a,
	0x0d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x12,
	0x2e, 0x64, 0x66, 0x73, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x13, 0x2e, 0x64, 0x66, 0x73, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x14, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12,
	0x16, 0x2e, 0x64, 0x66, 0x73, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12,
	0x3e, 0x0a, 0x0f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f,
	0x61, 0x64, 0x12, 0x14, 0x2e, 0x64, 0x66, 0x73, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x64, 0x66, 0x73, 0x2e, 0x44,
	0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3e, 0x0a, 0x0d, 0x53, 0x65, 0x6e, 0x64, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74,
	0x12, 0x15, 0x2e, 0x64, 0x66, 0x73, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x64, 0x66, 0x73, 0x2e, 0x48, 0x65,
	0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32,
	0x4e, 0x0a, 0x0a, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x12, 0x40, 0x0a,
	0x0d, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x17,
	0x2e, 0x64, 0x66, 0x73, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x64, 0x66, 0x73, 0x2e, 0x46, 0x69,
	0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x42,
	0x0b, 0x5a, 0x09, 0x64, 0x66, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_dfs_proto_rawDescOnce sync.Once
	file_dfs_proto_rawDescData []byte
)

func file_dfs_proto_rawDescGZIP() []byte {
	file_dfs_proto_rawDescOnce.Do(func() {
		file_dfs_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_dfs_proto_rawDesc), len(file_dfs_proto_rawDesc)))
	})
	return file_dfs_proto_rawDescData
}

var file_dfs_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_dfs_proto_goTypes = []any{
	(*ReplicationRequest)(nil),    // 0: dfs.ReplicationRequest
	(*UploadRequest)(nil),         // 1: dfs.UploadRequest
	(*FileUploadSuccess)(nil),     // 2: dfs.FileUploadSuccess
	(*UploadSuccessResponse)(nil), // 3: dfs.UploadSuccessResponse
	(*UploadResponse)(nil),        // 4: dfs.UploadResponse
	(*DownloadRequest)(nil),       // 5: dfs.DownloadRequest
	(*DownloadResponse)(nil),      // 6: dfs.DownloadResponse
	(*HeartbeatRequest)(nil),      // 7: dfs.HeartbeatRequest
	(*HeartbeatResponse)(nil),     // 8: dfs.HeartbeatResponse
	(*PortStatus)(nil),            // 9: dfs.PortStatus
	(*emptypb.Empty)(nil),         // 10: google.protobuf.Empty
}
var file_dfs_proto_depIdxs = []int32{
	9,  // 0: dfs.HeartbeatRequest.portsTCP:type_name -> dfs.PortStatus
	9,  // 1: dfs.HeartbeatRequest.portsGRPC:type_name -> dfs.PortStatus
	1,  // 2: dfs.MasterTracker.RequestUpload:input_type -> dfs.UploadRequest
	2,  // 3: dfs.MasterTracker.RequestUploadSuccess:input_type -> dfs.FileUploadSuccess
	5,  // 4: dfs.MasterTracker.RequestDownload:input_type -> dfs.DownloadRequest
	7,  // 5: dfs.MasterTracker.SendHeartbeat:input_type -> dfs.HeartbeatRequest
	0,  // 6: dfs.DataKeeper.ReplicateFile:input_type -> dfs.ReplicationRequest
	4,  // 7: dfs.MasterTracker.RequestUpload:output_type -> dfs.UploadResponse
	10, // 8: dfs.MasterTracker.RequestUploadSuccess:output_type -> google.protobuf.Empty
	6,  // 9: dfs.MasterTracker.RequestDownload:output_type -> dfs.DownloadResponse
	8,  // 10: dfs.MasterTracker.SendHeartbeat:output_type -> dfs.HeartbeatResponse
	2,  // 11: dfs.DataKeeper.ReplicateFile:output_type -> dfs.FileUploadSuccess
	7,  // [7:12] is the sub-list for method output_type
	2,  // [2:7] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_dfs_proto_init() }
func file_dfs_proto_init() {
	if File_dfs_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_dfs_proto_rawDesc), len(file_dfs_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_dfs_proto_goTypes,
		DependencyIndexes: file_dfs_proto_depIdxs,
		MessageInfos:      file_dfs_proto_msgTypes,
	}.Build()
	File_dfs_proto = out.File
	file_dfs_proto_goTypes = nil
	file_dfs_proto_depIdxs = nil
}
