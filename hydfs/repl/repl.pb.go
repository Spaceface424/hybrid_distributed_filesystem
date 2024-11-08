// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.12.4
// source: hydfs/repl/repl.proto

package repl

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filename string `protobuf:"bytes,1,opt,name=Filename,proto3" json:"Filename,omitempty"`
}

func (x *GetData) Reset() {
	*x = GetData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hydfs_repl_repl_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetData) ProtoMessage() {}

func (x *GetData) ProtoReflect() protoreflect.Message {
	mi := &file_hydfs_repl_repl_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetData.ProtoReflect.Descriptor instead.
func (*GetData) Descriptor() ([]byte, []int) {
	return file_hydfs_repl_repl_proto_rawDescGZIP(), []int{0}
}

func (x *GetData) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

type CreateData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NewFile *File `protobuf:"bytes,1,opt,name=NewFile,proto3" json:"NewFile,omitempty"`
}

func (x *CreateData) Reset() {
	*x = CreateData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hydfs_repl_repl_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateData) ProtoMessage() {}

func (x *CreateData) ProtoReflect() protoreflect.Message {
	mi := &file_hydfs_repl_repl_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateData.ProtoReflect.Descriptor instead.
func (*CreateData) Descriptor() ([]byte, []int) {
	return file_hydfs_repl_repl_proto_rawDescGZIP(), []int{1}
}

func (x *CreateData) GetNewFile() *File {
	if x != nil {
		return x.NewFile
	}
	return nil
}

type RequestFiles struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Files []*File `protobuf:"bytes,1,rep,name=Files,proto3" json:"Files,omitempty"`
}

func (x *RequestFiles) Reset() {
	*x = RequestFiles{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hydfs_repl_repl_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestFiles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestFiles) ProtoMessage() {}

func (x *RequestFiles) ProtoReflect() protoreflect.Message {
	mi := &file_hydfs_repl_repl_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestFiles.ProtoReflect.Descriptor instead.
func (*RequestFiles) Descriptor() ([]byte, []int) {
	return file_hydfs_repl_repl_proto_rawDescGZIP(), []int{2}
}

func (x *RequestFiles) GetFiles() []*File {
	if x != nil {
		return x.Files
	}
	return nil
}

type RequestMissing struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MissingFiles []*File `protobuf:"bytes,1,rep,name=MissingFiles,proto3" json:"MissingFiles,omitempty"`
}

func (x *RequestMissing) Reset() {
	*x = RequestMissing{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hydfs_repl_repl_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestMissing) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestMissing) ProtoMessage() {}

func (x *RequestMissing) ProtoReflect() protoreflect.Message {
	mi := &file_hydfs_repl_repl_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestMissing.ProtoReflect.Descriptor instead.
func (*RequestMissing) Descriptor() ([]byte, []int) {
	return file_hydfs_repl_repl_proto_rawDescGZIP(), []int{3}
}

func (x *RequestMissing) GetMissingFiles() []*File {
	if x != nil {
		return x.MissingFiles
	}
	return nil
}

type RequestData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataFiles []*File `protobuf:"bytes,1,rep,name=DataFiles,proto3" json:"DataFiles,omitempty"`
}

func (x *RequestData) Reset() {
	*x = RequestData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hydfs_repl_repl_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestData) ProtoMessage() {}

func (x *RequestData) ProtoReflect() protoreflect.Message {
	mi := &file_hydfs_repl_repl_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestData.ProtoReflect.Descriptor instead.
func (*RequestData) Descriptor() ([]byte, []int) {
	return file_hydfs_repl_repl_proto_rawDescGZIP(), []int{4}
}

func (x *RequestData) GetDataFiles() []*File {
	if x != nil {
		return x.DataFiles
	}
	return nil
}

type RequestAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OK bool `protobuf:"varint,1,opt,name=OK,proto3" json:"OK,omitempty"`
}

func (x *RequestAck) Reset() {
	*x = RequestAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hydfs_repl_repl_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestAck) ProtoMessage() {}

func (x *RequestAck) ProtoReflect() protoreflect.Message {
	mi := &file_hydfs_repl_repl_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestAck.ProtoReflect.Descriptor instead.
func (*RequestAck) Descriptor() ([]byte, []int) {
	return file_hydfs_repl_repl_proto_rawDescGZIP(), []int{5}
}

func (x *RequestAck) GetOK() bool {
	if x != nil {
		return x.OK
	}
	return false
}

type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filename string       `protobuf:"bytes,1,opt,name=Filename,proto3" json:"Filename,omitempty"`
	Blocks   []*FileBlock `protobuf:"bytes,2,rep,name=Blocks,proto3" json:"Blocks,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hydfs_repl_repl_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_hydfs_repl_repl_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_hydfs_repl_repl_proto_rawDescGZIP(), []int{6}
}

func (x *File) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *File) GetBlocks() []*FileBlock {
	if x != nil {
		return x.Blocks
	}
	return nil
}

type FileBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockNode uint32 `protobuf:"varint,1,opt,name=BlockNode,proto3" json:"BlockNode,omitempty"`
	BlockID   uint32 `protobuf:"varint,2,opt,name=BlockID,proto3" json:"BlockID,omitempty"`
	Data      []byte `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *FileBlock) Reset() {
	*x = FileBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hydfs_repl_repl_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileBlock) ProtoMessage() {}

func (x *FileBlock) ProtoReflect() protoreflect.Message {
	mi := &file_hydfs_repl_repl_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileBlock.ProtoReflect.Descriptor instead.
func (*FileBlock) Descriptor() ([]byte, []int) {
	return file_hydfs_repl_repl_proto_rawDescGZIP(), []int{7}
}

func (x *FileBlock) GetBlockNode() uint32 {
	if x != nil {
		return x.BlockNode
	}
	return 0
}

func (x *FileBlock) GetBlockID() uint32 {
	if x != nil {
		return x.BlockID
	}
	return 0
}

func (x *FileBlock) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_hydfs_repl_repl_proto protoreflect.FileDescriptor

var file_hydfs_repl_repl_proto_rawDesc = []byte{
	0x0a, 0x15, 0x68, 0x79, 0x64, 0x66, 0x73, 0x2f, 0x72, 0x65, 0x70, 0x6c, 0x2f, 0x72, 0x65, 0x70,
	0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x25, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2d,
	0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1f, 0x0a, 0x07,
	0x4e, 0x65, 0x77, 0x46, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e,
	0x46, 0x69, 0x6c, 0x65, 0x52, 0x07, 0x4e, 0x65, 0x77, 0x46, 0x69, 0x6c, 0x65, 0x22, 0x2b, 0x0a,
	0x0c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x1b, 0x0a,
	0x05, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x05, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x3b, 0x0a, 0x0e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x12, 0x29, 0x0a, 0x0c,
	0x4d, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x05, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x0c, 0x4d, 0x69, 0x73, 0x73, 0x69,
	0x6e, 0x67, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x32, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x23, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x46, 0x69,
	0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x09, 0x44, 0x61, 0x74, 0x61, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x1c, 0x0a, 0x0a, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x4f, 0x4b, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x4f, 0x4b, 0x22, 0x46, 0x0a, 0x04, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a,
	0x06, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x46, 0x69, 0x6c, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x06, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x73, 0x22, 0x57, 0x0a, 0x09, 0x46, 0x69, 0x6c, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x1c,
	0x0a, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x32, 0x92, 0x02, 0x0a, 0x0b, 0x52,
	0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2e, 0x0a, 0x0a, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x73, 0x6b, 0x12, 0x0d, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x1a, 0x0f, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x4d, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x22, 0x00, 0x12, 0x2a, 0x0a, 0x0b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x0c, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x0b, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x41, 0x63, 0x6b, 0x22, 0x00, 0x12, 0x2b, 0x0a, 0x0d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x0b, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x1a, 0x0b, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x63,
	0x6b, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x14, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x0b, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x0b, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x41, 0x63, 0x6b, 0x22, 0x00, 0x12, 0x1f, 0x0a, 0x0a, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x47, 0x65, 0x74, 0x12, 0x08, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x1a,
	0x05, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x22, 0x00, 0x12, 0x25, 0x0a, 0x0d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x12, 0x05, 0x2e, 0x46, 0x69, 0x6c, 0x65,
	0x1a, 0x0b, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x63, 0x6b, 0x22, 0x00, 0x42,
	0x08, 0x5a, 0x06, 0x2e, 0x2f, 0x72, 0x65, 0x70, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_hydfs_repl_repl_proto_rawDescOnce sync.Once
	file_hydfs_repl_repl_proto_rawDescData = file_hydfs_repl_repl_proto_rawDesc
)

func file_hydfs_repl_repl_proto_rawDescGZIP() []byte {
	file_hydfs_repl_repl_proto_rawDescOnce.Do(func() {
		file_hydfs_repl_repl_proto_rawDescData = protoimpl.X.CompressGZIP(file_hydfs_repl_repl_proto_rawDescData)
	})
	return file_hydfs_repl_repl_proto_rawDescData
}

var file_hydfs_repl_repl_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_hydfs_repl_repl_proto_goTypes = []any{
	(*GetData)(nil),        // 0: GetData
	(*CreateData)(nil),     // 1: CreateData
	(*RequestFiles)(nil),   // 2: RequestFiles
	(*RequestMissing)(nil), // 3: RequestMissing
	(*RequestData)(nil),    // 4: RequestData
	(*RequestAck)(nil),     // 5: RequestAck
	(*File)(nil),           // 6: File
	(*FileBlock)(nil),      // 7: FileBlock
}
var file_hydfs_repl_repl_proto_depIdxs = []int32{
	6,  // 0: CreateData.NewFile:type_name -> File
	6,  // 1: RequestFiles.Files:type_name -> File
	6,  // 2: RequestMissing.MissingFiles:type_name -> File
	6,  // 3: RequestData.DataFiles:type_name -> File
	7,  // 4: File.Blocks:type_name -> FileBlock
	2,  // 5: Replication.RequestAsk:input_type -> RequestFiles
	4,  // 6: Replication.RequestSend:input_type -> RequestData
	1,  // 7: Replication.RequestCreate:input_type -> CreateData
	1,  // 8: Replication.RequestReplicaCreate:input_type -> CreateData
	0,  // 9: Replication.RequestGet:input_type -> GetData
	6,  // 10: Replication.RequestAppend:input_type -> File
	3,  // 11: Replication.RequestAsk:output_type -> RequestMissing
	5,  // 12: Replication.RequestSend:output_type -> RequestAck
	5,  // 13: Replication.RequestCreate:output_type -> RequestAck
	5,  // 14: Replication.RequestReplicaCreate:output_type -> RequestAck
	6,  // 15: Replication.RequestGet:output_type -> File
	5,  // 16: Replication.RequestAppend:output_type -> RequestAck
	11, // [11:17] is the sub-list for method output_type
	5,  // [5:11] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_hydfs_repl_repl_proto_init() }
func file_hydfs_repl_repl_proto_init() {
	if File_hydfs_repl_repl_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_hydfs_repl_repl_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GetData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hydfs_repl_repl_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreateData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hydfs_repl_repl_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*RequestFiles); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hydfs_repl_repl_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*RequestMissing); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hydfs_repl_repl_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*RequestData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hydfs_repl_repl_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*RequestAck); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hydfs_repl_repl_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*File); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hydfs_repl_repl_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*FileBlock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_hydfs_repl_repl_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_hydfs_repl_repl_proto_goTypes,
		DependencyIndexes: file_hydfs_repl_repl_proto_depIdxs,
		MessageInfos:      file_hydfs_repl_repl_proto_msgTypes,
	}.Build()
	File_hydfs_repl_repl_proto = out.File
	file_hydfs_repl_repl_proto_rawDesc = nil
	file_hydfs_repl_repl_proto_goTypes = nil
	file_hydfs_repl_repl_proto_depIdxs = nil
}
