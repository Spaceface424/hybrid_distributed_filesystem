// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: shared/pingack.proto

package shared

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

type MessageType int32

const (
	MessageType_PING     MessageType = 0
	MessageType_ACK      MessageType = 1
	MessageType_PING_REQ MessageType = 2
	MessageType_IDR_PING MessageType = 3
	MessageType_ACK_REQ  MessageType = 4
)

// Enum value maps for MessageType.
var (
	MessageType_name = map[int32]string{
		0: "PING",
		1: "ACK",
		2: "PING_REQ",
		3: "IDR_PING",
		4: "ACK_REQ",
	}
	MessageType_value = map[string]int32{
		"PING":     0,
		"ACK":      1,
		"PING_REQ": 2,
		"IDR_PING": 3,
		"ACK_REQ":  4,
	}
)

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}

func (x MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_shared_pingack_proto_enumTypes[0].Descriptor()
}

func (MessageType) Type() protoreflect.EnumType {
	return &file_shared_pingack_proto_enumTypes[0]
}

func (x MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageType.Descriptor instead.
func (MessageType) EnumDescriptor() ([]byte, []int) {
	return file_shared_pingack_proto_rawDescGZIP(), []int{0}
}

type NodeState int32

const (
	NodeState_ALIVE   NodeState = 0
	NodeState_FAILED  NodeState = 1
	NodeState_SUSPECT NodeState = 2
)

// Enum value maps for NodeState.
var (
	NodeState_name = map[int32]string{
		0: "ALIVE",
		1: "FAILED",
		2: "SUSPECT",
	}
	NodeState_value = map[string]int32{
		"ALIVE":   0,
		"FAILED":  1,
		"SUSPECT": 2,
	}
)

func (x NodeState) Enum() *NodeState {
	p := new(NodeState)
	*p = x
	return p
}

func (x NodeState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NodeState) Descriptor() protoreflect.EnumDescriptor {
	return file_shared_pingack_proto_enumTypes[1].Descriptor()
}

func (NodeState) Type() protoreflect.EnumType {
	return &file_shared_pingack_proto_enumTypes[1]
}

func (x NodeState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NodeState.Descriptor instead.
func (NodeState) EnumDescriptor() ([]byte, []int) {
	return file_shared_pingack_proto_rawDescGZIP(), []int{1}
}

type SWIMIn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd      string  `protobuf:"bytes,1,opt,name=Cmd,proto3" json:"Cmd,omitempty"`
	DropRate float64 `protobuf:"fixed64,2,opt,name=drop_rate,json=dropRate,proto3" json:"drop_rate,omitempty"`
}

func (x *SWIMIn) Reset() {
	*x = SWIMIn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_pingack_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SWIMIn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SWIMIn) ProtoMessage() {}

func (x *SWIMIn) ProtoReflect() protoreflect.Message {
	mi := &file_shared_pingack_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SWIMIn.ProtoReflect.Descriptor instead.
func (*SWIMIn) Descriptor() ([]byte, []int) {
	return file_shared_pingack_proto_rawDescGZIP(), []int{0}
}

func (x *SWIMIn) GetCmd() string {
	if x != nil {
		return x.Cmd
	}
	return ""
}

func (x *SWIMIn) GetDropRate() float64 {
	if x != nil {
		return x.DropRate
	}
	return 0
}

type SWIMOut struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Output string `protobuf:"bytes,1,opt,name=output,proto3" json:"output,omitempty"`
}

func (x *SWIMOut) Reset() {
	*x = SWIMOut{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_pingack_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SWIMOut) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SWIMOut) ProtoMessage() {}

func (x *SWIMOut) ProtoReflect() protoreflect.Message {
	mi := &file_shared_pingack_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SWIMOut.ProtoReflect.Descriptor instead.
func (*SWIMOut) Descriptor() ([]byte, []int) {
	return file_shared_pingack_proto_rawDescGZIP(), []int{1}
}

func (x *SWIMOut) GetOutput() string {
	if x != nil {
		return x.Output
	}
	return ""
}

type PingAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type         MessageType       `protobuf:"varint,1,opt,name=Type,proto3,enum=MessageType" json:"Type,omitempty"`
	SenderId     int32             `protobuf:"varint,2,opt,name=SenderId,proto3" json:"SenderId,omitempty"`
	Round        int32             `protobuf:"varint,3,opt,name=Round,proto3" json:"Round,omitempty"`
	GossipBuffer map[int32]*Gossip `protobuf:"bytes,4,rep,name=GossipBuffer,proto3" json:"GossipBuffer,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	RequestId    int32             `protobuf:"varint,5,opt,name=RequestId,proto3" json:"RequestId,omitempty"`
	IncNum       int32             `protobuf:"varint,6,opt,name=IncNum,proto3" json:"IncNum,omitempty"`
}

func (x *PingAck) Reset() {
	*x = PingAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_pingack_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingAck) ProtoMessage() {}

func (x *PingAck) ProtoReflect() protoreflect.Message {
	mi := &file_shared_pingack_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingAck.ProtoReflect.Descriptor instead.
func (*PingAck) Descriptor() ([]byte, []int) {
	return file_shared_pingack_proto_rawDescGZIP(), []int{2}
}

func (x *PingAck) GetType() MessageType {
	if x != nil {
		return x.Type
	}
	return MessageType_PING
}

func (x *PingAck) GetSenderId() int32 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *PingAck) GetRound() int32 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *PingAck) GetGossipBuffer() map[int32]*Gossip {
	if x != nil {
		return x.GossipBuffer
	}
	return nil
}

func (x *PingAck) GetRequestId() int32 {
	if x != nil {
		return x.RequestId
	}
	return 0
}

func (x *PingAck) GetIncNum() int32 {
	if x != nil {
		return x.IncNum
	}
	return 0
}

type Gossip struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Member *MemberInfo `protobuf:"bytes,1,opt,name=Member,proto3" json:"Member,omitempty"`
	TTL    int32       `protobuf:"varint,2,opt,name=TTL,proto3" json:"TTL,omitempty"`
}

func (x *Gossip) Reset() {
	*x = Gossip{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_pingack_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Gossip) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Gossip) ProtoMessage() {}

func (x *Gossip) ProtoReflect() protoreflect.Message {
	mi := &file_shared_pingack_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Gossip.ProtoReflect.Descriptor instead.
func (*Gossip) Descriptor() ([]byte, []int) {
	return file_shared_pingack_proto_rawDescGZIP(), []int{3}
}

func (x *Gossip) GetMember() *MemberInfo {
	if x != nil {
		return x.Member
	}
	return nil
}

func (x *Gossip) GetTTL() int32 {
	if x != nil {
		return x.TTL
	}
	return 0
}

type MemberContainer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID         int32         `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Hash       int32         `protobuf:"varint,2,opt,name=Hash,proto3" json:"Hash,omitempty"`
	EnableSus  bool          `protobuf:"varint,3,opt,name=EnableSus,proto3" json:"EnableSus,omitempty"`
	MemberList []*MemberInfo `protobuf:"bytes,4,rep,name=MemberList,proto3" json:"MemberList,omitempty"`
	FailList   []*MemberInfo `protobuf:"bytes,5,rep,name=FailList,proto3" json:"FailList,omitempty"`
}

func (x *MemberContainer) Reset() {
	*x = MemberContainer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_pingack_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberContainer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberContainer) ProtoMessage() {}

func (x *MemberContainer) ProtoReflect() protoreflect.Message {
	mi := &file_shared_pingack_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberContainer.ProtoReflect.Descriptor instead.
func (*MemberContainer) Descriptor() ([]byte, []int) {
	return file_shared_pingack_proto_rawDescGZIP(), []int{4}
}

func (x *MemberContainer) GetID() int32 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *MemberContainer) GetHash() int32 {
	if x != nil {
		return x.Hash
	}
	return 0
}

func (x *MemberContainer) GetEnableSus() bool {
	if x != nil {
		return x.EnableSus
	}
	return false
}

func (x *MemberContainer) GetMemberList() []*MemberInfo {
	if x != nil {
		return x.MemberList
	}
	return nil
}

func (x *MemberContainer) GetFailList() []*MemberInfo {
	if x != nil {
		return x.FailList
	}
	return nil
}

type MemberInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string    `protobuf:"bytes,1,opt,name=Address,proto3" json:"Address,omitempty"`
	ID      int32     `protobuf:"varint,2,opt,name=ID,proto3" json:"ID,omitempty"`
	Hash    uint32    `protobuf:"varint,3,opt,name=Hash,proto3" json:"Hash,omitempty"`
	State   NodeState `protobuf:"varint,4,opt,name=State,proto3,enum=NodeState" json:"State,omitempty"`
	IncNum  int32     `protobuf:"varint,5,opt,name=IncNum,proto3" json:"IncNum,omitempty"`
}

func (x *MemberInfo) Reset() {
	*x = MemberInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_pingack_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberInfo) ProtoMessage() {}

func (x *MemberInfo) ProtoReflect() protoreflect.Message {
	mi := &file_shared_pingack_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberInfo.ProtoReflect.Descriptor instead.
func (*MemberInfo) Descriptor() ([]byte, []int) {
	return file_shared_pingack_proto_rawDescGZIP(), []int{5}
}

func (x *MemberInfo) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *MemberInfo) GetID() int32 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *MemberInfo) GetHash() uint32 {
	if x != nil {
		return x.Hash
	}
	return 0
}

func (x *MemberInfo) GetState() NodeState {
	if x != nil {
		return x.State
	}
	return NodeState_ALIVE
}

func (x *MemberInfo) GetIncNum() int32 {
	if x != nil {
		return x.IncNum
	}
	return 0
}

var File_shared_pingack_proto protoreflect.FileDescriptor

var file_shared_pingack_proto_rawDesc = []byte{
	0x0a, 0x14, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x70, 0x69, 0x6e, 0x67, 0x61, 0x63, 0x6b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x37, 0x0a, 0x06, 0x53, 0x57, 0x49, 0x4d, 0x49, 0x6e,
	0x12, 0x10, 0x0a, 0x03, 0x43, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x43,
	0x6d, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x72, 0x6f, 0x70, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x64, 0x72, 0x6f, 0x70, 0x52, 0x61, 0x74, 0x65, 0x22,
	0x21, 0x0a, 0x07, 0x53, 0x57, 0x49, 0x4d, 0x4f, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70,
	0x75, 0x74, 0x22, 0x9d, 0x02, 0x0a, 0x07, 0x50, 0x69, 0x6e, 0x67, 0x41, 0x63, 0x6b, 0x12, 0x20,
	0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x52, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x52, 0x6f, 0x75,
	0x6e, 0x64, 0x12, 0x3e, 0x0a, 0x0c, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x42, 0x75, 0x66, 0x66,
	0x65, 0x72, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x41,
	0x63, 0x6b, 0x2e, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x42, 0x75, 0x66, 0x66, 0x65, 0x72, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x0c, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x42, 0x75, 0x66, 0x66,
	0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x49, 0x6e, 0x63, 0x4e, 0x75, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x49, 0x6e, 0x63, 0x4e, 0x75, 0x6d, 0x1a, 0x48, 0x0a, 0x11, 0x47, 0x6f, 0x73, 0x73,
	0x69, 0x70, 0x42, 0x75, 0x66, 0x66, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x1d, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07,
	0x2e, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x3f, 0x0a, 0x06, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x12, 0x23, 0x0a, 0x06,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x06, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x10, 0x0a, 0x03, 0x54, 0x54, 0x4c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x54, 0x54, 0x4c, 0x22, 0xa9, 0x01, 0x0a, 0x0f, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f,
	0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x45,
	0x6e, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09,
	0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x75, 0x73, 0x12, 0x2b, 0x0a, 0x0a, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x08, 0x46, 0x61, 0x69, 0x6c, 0x4c, 0x69,
	0x73, 0x74, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x46, 0x61, 0x69, 0x6c, 0x4c, 0x69, 0x73, 0x74, 0x22,
	0x84, 0x01, 0x0a, 0x0a, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18,
	0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x12, 0x20, 0x0a, 0x05,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x49, 0x6e, 0x63, 0x4e, 0x75, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x49, 0x6e, 0x63, 0x4e, 0x75, 0x6d, 0x2a, 0x49, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x49, 0x4e, 0x47, 0x10, 0x00, 0x12,
	0x07, 0x0a, 0x03, 0x41, 0x43, 0x4b, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x49, 0x4e, 0x47,
	0x5f, 0x52, 0x45, 0x51, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x49, 0x44, 0x52, 0x5f, 0x50, 0x49,
	0x4e, 0x47, 0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x43, 0x4b, 0x5f, 0x52, 0x45, 0x51, 0x10,
	0x04, 0x2a, 0x2f, 0x0a, 0x09, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x09,
	0x0a, 0x05, 0x41, 0x4c, 0x49, 0x56, 0x45, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x53, 0x50, 0x45, 0x43, 0x54,
	0x10, 0x02, 0x32, 0x66, 0x0a, 0x0a, 0x49, 0x6e, 0x74, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72,
	0x12, 0x38, 0x0a, 0x15, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x68, 0x69, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0b, 0x2e, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x10, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43,
	0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x22, 0x00, 0x12, 0x1e, 0x0a, 0x07, 0x53, 0x57,
	0x49, 0x4d, 0x63, 0x6d, 0x64, 0x12, 0x07, 0x2e, 0x53, 0x57, 0x49, 0x4d, 0x49, 0x6e, 0x1a, 0x08,
	0x2e, 0x53, 0x57, 0x49, 0x4d, 0x4f, 0x75, 0x74, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shared_pingack_proto_rawDescOnce sync.Once
	file_shared_pingack_proto_rawDescData = file_shared_pingack_proto_rawDesc
)

func file_shared_pingack_proto_rawDescGZIP() []byte {
	file_shared_pingack_proto_rawDescOnce.Do(func() {
		file_shared_pingack_proto_rawDescData = protoimpl.X.CompressGZIP(file_shared_pingack_proto_rawDescData)
	})
	return file_shared_pingack_proto_rawDescData
}

var file_shared_pingack_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_shared_pingack_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_shared_pingack_proto_goTypes = []any{
	(MessageType)(0),        // 0: MessageType
	(NodeState)(0),          // 1: NodeState
	(*SWIMIn)(nil),          // 2: SWIMIn
	(*SWIMOut)(nil),         // 3: SWIMOut
	(*PingAck)(nil),         // 4: PingAck
	(*Gossip)(nil),          // 5: Gossip
	(*MemberContainer)(nil), // 6: MemberContainer
	(*MemberInfo)(nil),      // 7: MemberInfo
	nil,                     // 8: PingAck.GossipBufferEntry
}
var file_shared_pingack_proto_depIdxs = []int32{
	0, // 0: PingAck.Type:type_name -> MessageType
	8, // 1: PingAck.GossipBuffer:type_name -> PingAck.GossipBufferEntry
	7, // 2: Gossip.Member:type_name -> MemberInfo
	7, // 3: MemberContainer.MemberList:type_name -> MemberInfo
	7, // 4: MemberContainer.FailList:type_name -> MemberInfo
	1, // 5: MemberInfo.State:type_name -> NodeState
	5, // 6: PingAck.GossipBufferEntry.value:type_name -> Gossip
	7, // 7: Introducer.RequestMembershipInfo:input_type -> MemberInfo
	2, // 8: Introducer.SWIMcmd:input_type -> SWIMIn
	6, // 9: Introducer.RequestMembershipInfo:output_type -> MemberContainer
	3, // 10: Introducer.SWIMcmd:output_type -> SWIMOut
	9, // [9:11] is the sub-list for method output_type
	7, // [7:9] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_shared_pingack_proto_init() }
func file_shared_pingack_proto_init() {
	if File_shared_pingack_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shared_pingack_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*SWIMIn); i {
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
		file_shared_pingack_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SWIMOut); i {
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
		file_shared_pingack_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*PingAck); i {
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
		file_shared_pingack_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Gossip); i {
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
		file_shared_pingack_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*MemberContainer); i {
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
		file_shared_pingack_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*MemberInfo); i {
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
			RawDescriptor: file_shared_pingack_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_shared_pingack_proto_goTypes,
		DependencyIndexes: file_shared_pingack_proto_depIdxs,
		EnumInfos:         file_shared_pingack_proto_enumTypes,
		MessageInfos:      file_shared_pingack_proto_msgTypes,
	}.Build()
	File_shared_pingack_proto = out.File
	file_shared_pingack_proto_rawDesc = nil
	file_shared_pingack_proto_goTypes = nil
	file_shared_pingack_proto_depIdxs = nil
}
