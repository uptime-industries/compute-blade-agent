// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: api/bladeapi/v1alpha1/blade.proto

package bladeapiv1alpha1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Event is an event the agent reacts to
type Event int32

const (
	Event_IDENTIFY         Event = 0
	Event_IDENTIFY_CONFIRM Event = 1
	Event_CRITICAL         Event = 2
	Event_CRITICAL_RESET   Event = 3
)

// Enum value maps for Event.
var (
	Event_name = map[int32]string{
		0: "IDENTIFY",
		1: "IDENTIFY_CONFIRM",
		2: "CRITICAL",
		3: "CRITICAL_RESET",
	}
	Event_value = map[string]int32{
		"IDENTIFY":         0,
		"IDENTIFY_CONFIRM": 1,
		"CRITICAL":         2,
		"CRITICAL_RESET":   3,
	}
)

func (x Event) Enum() *Event {
	p := new(Event)
	*p = x
	return p
}

func (x Event) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Event) Descriptor() protoreflect.EnumDescriptor {
	return file_api_bladeapi_v1alpha1_blade_proto_enumTypes[0].Descriptor()
}

func (Event) Type() protoreflect.EnumType {
	return &file_api_bladeapi_v1alpha1_blade_proto_enumTypes[0]
}

func (x Event) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Event.Descriptor instead.
func (Event) EnumDescriptor() ([]byte, []int) {
	return file_api_bladeapi_v1alpha1_blade_proto_rawDescGZIP(), []int{0}
}

// FanUnit defines the fan unit detected by the blade
type FanUnit int32

const (
	FanUnit_DEFAULT FanUnit = 0
	FanUnit_SMART   FanUnit = 1
)

// Enum value maps for FanUnit.
var (
	FanUnit_name = map[int32]string{
		0: "DEFAULT",
		1: "SMART",
	}
	FanUnit_value = map[string]int32{
		"DEFAULT": 0,
		"SMART":   1,
	}
)

func (x FanUnit) Enum() *FanUnit {
	p := new(FanUnit)
	*p = x
	return p
}

func (x FanUnit) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FanUnit) Descriptor() protoreflect.EnumDescriptor {
	return file_api_bladeapi_v1alpha1_blade_proto_enumTypes[1].Descriptor()
}

func (FanUnit) Type() protoreflect.EnumType {
	return &file_api_bladeapi_v1alpha1_blade_proto_enumTypes[1]
}

func (x FanUnit) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FanUnit.Descriptor instead.
func (FanUnit) EnumDescriptor() ([]byte, []int) {
	return file_api_bladeapi_v1alpha1_blade_proto_rawDescGZIP(), []int{1}
}

// PowerStatus defines the power status of the blade
type PowerStatus int32

const (
	PowerStatus_POE_OR_USBC PowerStatus = 0
	PowerStatus_POE_802_AT  PowerStatus = 1
)

// Enum value maps for PowerStatus.
var (
	PowerStatus_name = map[int32]string{
		0: "POE_OR_USBC",
		1: "POE_802_AT",
	}
	PowerStatus_value = map[string]int32{
		"POE_OR_USBC": 0,
		"POE_802_AT":  1,
	}
)

func (x PowerStatus) Enum() *PowerStatus {
	p := new(PowerStatus)
	*p = x
	return p
}

func (x PowerStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PowerStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_bladeapi_v1alpha1_blade_proto_enumTypes[2].Descriptor()
}

func (PowerStatus) Type() protoreflect.EnumType {
	return &file_api_bladeapi_v1alpha1_blade_proto_enumTypes[2]
}

func (x PowerStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PowerStatus.Descriptor instead.
func (PowerStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_bladeapi_v1alpha1_blade_proto_rawDescGZIP(), []int{2}
}

type StealthModeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enable bool `protobuf:"varint,1,opt,name=enable,proto3" json:"enable,omitempty"`
}

func (x *StealthModeRequest) Reset() {
	*x = StealthModeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_bladeapi_v1alpha1_blade_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StealthModeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StealthModeRequest) ProtoMessage() {}

func (x *StealthModeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_bladeapi_v1alpha1_blade_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StealthModeRequest.ProtoReflect.Descriptor instead.
func (*StealthModeRequest) Descriptor() ([]byte, []int) {
	return file_api_bladeapi_v1alpha1_blade_proto_rawDescGZIP(), []int{0}
}

func (x *StealthModeRequest) GetEnable() bool {
	if x != nil {
		return x.Enable
	}
	return false
}

type SetFanSpeedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Percent int64 `protobuf:"varint,1,opt,name=percent,proto3" json:"percent,omitempty"`
}

func (x *SetFanSpeedRequest) Reset() {
	*x = SetFanSpeedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_bladeapi_v1alpha1_blade_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetFanSpeedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetFanSpeedRequest) ProtoMessage() {}

func (x *SetFanSpeedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_bladeapi_v1alpha1_blade_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetFanSpeedRequest.ProtoReflect.Descriptor instead.
func (*SetFanSpeedRequest) Descriptor() ([]byte, []int) {
	return file_api_bladeapi_v1alpha1_blade_proto_rawDescGZIP(), []int{1}
}

func (x *SetFanSpeedRequest) GetPercent() int64 {
	if x != nil {
		return x.Percent
	}
	return 0
}

type EmitEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event Event `protobuf:"varint,1,opt,name=event,proto3,enum=api.bladeapi.v1alpha1.Event" json:"event,omitempty"`
}

func (x *EmitEventRequest) Reset() {
	*x = EmitEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_bladeapi_v1alpha1_blade_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmitEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmitEventRequest) ProtoMessage() {}

func (x *EmitEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_bladeapi_v1alpha1_blade_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmitEventRequest.ProtoReflect.Descriptor instead.
func (*EmitEventRequest) Descriptor() ([]byte, []int) {
	return file_api_bladeapi_v1alpha1_blade_proto_rawDescGZIP(), []int{2}
}

func (x *EmitEventRequest) GetEvent() Event {
	if x != nil {
		return x.Event
	}
	return Event_IDENTIFY
}

type StatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StealthMode    bool        `protobuf:"varint,1,opt,name=stealth_mode,json=stealthMode,proto3" json:"stealth_mode,omitempty"`
	IdentifyActive bool        `protobuf:"varint,2,opt,name=identify_active,json=identifyActive,proto3" json:"identify_active,omitempty"`
	CriticalActive bool        `protobuf:"varint,3,opt,name=critical_active,json=criticalActive,proto3" json:"critical_active,omitempty"`
	Temperature    int64       `protobuf:"varint,4,opt,name=temperature,proto3" json:"temperature,omitempty"`
	FanRpm         int64       `protobuf:"varint,5,opt,name=fan_rpm,json=fanRpm,proto3" json:"fan_rpm,omitempty"`
	PowerStatus    PowerStatus `protobuf:"varint,6,opt,name=power_status,json=powerStatus,proto3,enum=api.bladeapi.v1alpha1.PowerStatus" json:"power_status,omitempty"`
}

func (x *StatusResponse) Reset() {
	*x = StatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_bladeapi_v1alpha1_blade_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusResponse) ProtoMessage() {}

func (x *StatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_bladeapi_v1alpha1_blade_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusResponse.ProtoReflect.Descriptor instead.
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return file_api_bladeapi_v1alpha1_blade_proto_rawDescGZIP(), []int{3}
}

func (x *StatusResponse) GetStealthMode() bool {
	if x != nil {
		return x.StealthMode
	}
	return false
}

func (x *StatusResponse) GetIdentifyActive() bool {
	if x != nil {
		return x.IdentifyActive
	}
	return false
}

func (x *StatusResponse) GetCriticalActive() bool {
	if x != nil {
		return x.CriticalActive
	}
	return false
}

func (x *StatusResponse) GetTemperature() int64 {
	if x != nil {
		return x.Temperature
	}
	return 0
}

func (x *StatusResponse) GetFanRpm() int64 {
	if x != nil {
		return x.FanRpm
	}
	return 0
}

func (x *StatusResponse) GetPowerStatus() PowerStatus {
	if x != nil {
		return x.PowerStatus
	}
	return PowerStatus_POE_OR_USBC
}

var File_api_bladeapi_v1alpha1_blade_proto protoreflect.FileDescriptor

var file_api_bladeapi_v1alpha1_blade_proto_rawDesc = []byte{
	0x0a, 0x21, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x6c, 0x61, 0x64, 0x65, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x62, 0x6c, 0x61, 0x64, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x15, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x6c, 0x61, 0x64, 0x65, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2c, 0x0a, 0x12, 0x53, 0x74, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x4d, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x65,
	0x6e, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x2e, 0x0a, 0x12, 0x53, 0x65, 0x74, 0x46, 0x61, 0x6e, 0x53,
	0x70, 0x65, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70,
	0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x70, 0x65,
	0x72, 0x63, 0x65, 0x6e, 0x74, 0x22, 0x46, 0x0a, 0x10, 0x45, 0x6d, 0x69, 0x74, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x05, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x62,
	0x6c, 0x61, 0x64, 0x65, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x87, 0x02,
	0x0a, 0x0e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x21, 0x0a, 0x0c, 0x73, 0x74, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x5f, 0x6d, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x73, 0x74, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x4d,
	0x6f, 0x64, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x79, 0x5f,
	0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x69, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x66, 0x79, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x27, 0x0a, 0x0f,
	0x63, 0x72, 0x69, 0x74, 0x69, 0x63, 0x61, 0x6c, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x63, 0x72, 0x69, 0x74, 0x69, 0x63, 0x61, 0x6c, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x74, 0x65, 0x6d, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x61, 0x6e, 0x5f, 0x72,
	0x70, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x66, 0x61, 0x6e, 0x52, 0x70, 0x6d,
	0x12, 0x45, 0x0a, 0x0c, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x6c, 0x61,
	0x64, 0x65, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x50,
	0x6f, 0x77, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0b, 0x70, 0x6f, 0x77, 0x65,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2a, 0x4d, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x0c, 0x0a, 0x08, 0x49, 0x44, 0x45, 0x4e, 0x54, 0x49, 0x46, 0x59, 0x10, 0x00, 0x12, 0x14,
	0x0a, 0x10, 0x49, 0x44, 0x45, 0x4e, 0x54, 0x49, 0x46, 0x59, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49,
	0x52, 0x4d, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x52, 0x49, 0x54, 0x49, 0x43, 0x41, 0x4c,
	0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x52, 0x49, 0x54, 0x49, 0x43, 0x41, 0x4c, 0x5f, 0x52,
	0x45, 0x53, 0x45, 0x54, 0x10, 0x03, 0x2a, 0x21, 0x0a, 0x07, 0x46, 0x61, 0x6e, 0x55, 0x6e, 0x69,
	0x74, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x45, 0x46, 0x41, 0x55, 0x4c, 0x54, 0x10, 0x00, 0x12, 0x09,
	0x0a, 0x05, 0x53, 0x4d, 0x41, 0x52, 0x54, 0x10, 0x01, 0x2a, 0x2e, 0x0a, 0x0b, 0x50, 0x6f, 0x77,
	0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0f, 0x0a, 0x0b, 0x50, 0x4f, 0x45, 0x5f,
	0x4f, 0x52, 0x5f, 0x55, 0x53, 0x42, 0x43, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x4f, 0x45,
	0x5f, 0x38, 0x30, 0x32, 0x5f, 0x41, 0x54, 0x10, 0x01, 0x32, 0xa8, 0x03, 0x0a, 0x11, 0x42, 0x6c,
	0x61, 0x64, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x4e, 0x0a, 0x09, 0x45, 0x6d, 0x69, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x27, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x62, 0x6c, 0x61, 0x64, 0x65, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2e, 0x45, 0x6d, 0x69, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12,
	0x4a, 0x0a, 0x16, 0x57, 0x61, 0x69, 0x74, 0x46, 0x6f, 0x72, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x52, 0x0a, 0x0b, 0x53,
	0x65, 0x74, 0x46, 0x61, 0x6e, 0x53, 0x70, 0x65, 0x65, 0x64, 0x12, 0x29, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x62, 0x6c, 0x61, 0x64, 0x65, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x2e, 0x53, 0x65, 0x74, 0x46, 0x61, 0x6e, 0x53, 0x70, 0x65, 0x65, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12,
	0x55, 0x0a, 0x0e, 0x53, 0x65, 0x74, 0x53, 0x74, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x4d, 0x6f, 0x64,
	0x65, 0x12, 0x29, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x6c, 0x61, 0x64, 0x65, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x53, 0x74, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x4d, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x4c, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x25, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x62, 0x6c, 0x61, 0x64, 0x65, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x48, 0x5a, 0x46, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x78, 0x76, 0x7a, 0x66, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x62,
	0x6c, 0x61, 0x64, 0x65, 0x2d, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x62,
	0x6c, 0x61, 0x64, 0x65, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x62, 0x6c,
	0x61, 0x64, 0x65, 0x61, 0x70, 0x69, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_bladeapi_v1alpha1_blade_proto_rawDescOnce sync.Once
	file_api_bladeapi_v1alpha1_blade_proto_rawDescData = file_api_bladeapi_v1alpha1_blade_proto_rawDesc
)

func file_api_bladeapi_v1alpha1_blade_proto_rawDescGZIP() []byte {
	file_api_bladeapi_v1alpha1_blade_proto_rawDescOnce.Do(func() {
		file_api_bladeapi_v1alpha1_blade_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_bladeapi_v1alpha1_blade_proto_rawDescData)
	})
	return file_api_bladeapi_v1alpha1_blade_proto_rawDescData
}

var file_api_bladeapi_v1alpha1_blade_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_api_bladeapi_v1alpha1_blade_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_bladeapi_v1alpha1_blade_proto_goTypes = []interface{}{
	(Event)(0),                 // 0: api.bladeapi.v1alpha1.Event
	(FanUnit)(0),               // 1: api.bladeapi.v1alpha1.FanUnit
	(PowerStatus)(0),           // 2: api.bladeapi.v1alpha1.PowerStatus
	(*StealthModeRequest)(nil), // 3: api.bladeapi.v1alpha1.StealthModeRequest
	(*SetFanSpeedRequest)(nil), // 4: api.bladeapi.v1alpha1.SetFanSpeedRequest
	(*EmitEventRequest)(nil),   // 5: api.bladeapi.v1alpha1.EmitEventRequest
	(*StatusResponse)(nil),     // 6: api.bladeapi.v1alpha1.StatusResponse
	(*emptypb.Empty)(nil),      // 7: google.protobuf.Empty
}
var file_api_bladeapi_v1alpha1_blade_proto_depIdxs = []int32{
	0, // 0: api.bladeapi.v1alpha1.EmitEventRequest.event:type_name -> api.bladeapi.v1alpha1.Event
	2, // 1: api.bladeapi.v1alpha1.StatusResponse.power_status:type_name -> api.bladeapi.v1alpha1.PowerStatus
	5, // 2: api.bladeapi.v1alpha1.BladeAgentService.EmitEvent:input_type -> api.bladeapi.v1alpha1.EmitEventRequest
	7, // 3: api.bladeapi.v1alpha1.BladeAgentService.WaitForIdentifyConfirm:input_type -> google.protobuf.Empty
	4, // 4: api.bladeapi.v1alpha1.BladeAgentService.SetFanSpeed:input_type -> api.bladeapi.v1alpha1.SetFanSpeedRequest
	3, // 5: api.bladeapi.v1alpha1.BladeAgentService.SetStealthMode:input_type -> api.bladeapi.v1alpha1.StealthModeRequest
	7, // 6: api.bladeapi.v1alpha1.BladeAgentService.GetStatus:input_type -> google.protobuf.Empty
	7, // 7: api.bladeapi.v1alpha1.BladeAgentService.EmitEvent:output_type -> google.protobuf.Empty
	7, // 8: api.bladeapi.v1alpha1.BladeAgentService.WaitForIdentifyConfirm:output_type -> google.protobuf.Empty
	7, // 9: api.bladeapi.v1alpha1.BladeAgentService.SetFanSpeed:output_type -> google.protobuf.Empty
	7, // 10: api.bladeapi.v1alpha1.BladeAgentService.SetStealthMode:output_type -> google.protobuf.Empty
	6, // 11: api.bladeapi.v1alpha1.BladeAgentService.GetStatus:output_type -> api.bladeapi.v1alpha1.StatusResponse
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_bladeapi_v1alpha1_blade_proto_init() }
func file_api_bladeapi_v1alpha1_blade_proto_init() {
	if File_api_bladeapi_v1alpha1_blade_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_bladeapi_v1alpha1_blade_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StealthModeRequest); i {
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
		file_api_bladeapi_v1alpha1_blade_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetFanSpeedRequest); i {
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
		file_api_bladeapi_v1alpha1_blade_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmitEventRequest); i {
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
		file_api_bladeapi_v1alpha1_blade_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusResponse); i {
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
			RawDescriptor: file_api_bladeapi_v1alpha1_blade_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_bladeapi_v1alpha1_blade_proto_goTypes,
		DependencyIndexes: file_api_bladeapi_v1alpha1_blade_proto_depIdxs,
		EnumInfos:         file_api_bladeapi_v1alpha1_blade_proto_enumTypes,
		MessageInfos:      file_api_bladeapi_v1alpha1_blade_proto_msgTypes,
	}.Build()
	File_api_bladeapi_v1alpha1_blade_proto = out.File
	file_api_bladeapi_v1alpha1_blade_proto_rawDesc = nil
	file_api_bladeapi_v1alpha1_blade_proto_goTypes = nil
	file_api_bladeapi_v1alpha1_blade_proto_depIdxs = nil
}