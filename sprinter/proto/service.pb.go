// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: proto/service.proto

package proto

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

type Status_StatusEnum int32

const (
	Status_ERROR     Status_StatusEnum = 0
	Status_RUNNING   Status_StatusEnum = 1
	Status_COMPLETED Status_StatusEnum = 2
	Status_ABORTED   Status_StatusEnum = 3
)

// Enum value maps for Status_StatusEnum.
var (
	Status_StatusEnum_name = map[int32]string{
		0: "ERROR",
		1: "RUNNING",
		2: "COMPLETED",
		3: "ABORTED",
	}
	Status_StatusEnum_value = map[string]int32{
		"ERROR":     0,
		"RUNNING":   1,
		"COMPLETED": 2,
		"ABORTED":   3,
	}
)

func (x Status_StatusEnum) Enum() *Status_StatusEnum {
	p := new(Status_StatusEnum)
	*p = x
	return p
}

func (x Status_StatusEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status_StatusEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_service_proto_enumTypes[0].Descriptor()
}

func (Status_StatusEnum) Type() protoreflect.EnumType {
	return &file_proto_service_proto_enumTypes[0]
}

func (x Status_StatusEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status_StatusEnum.Descriptor instead.
func (Status_StatusEnum) EnumDescriptor() ([]byte, []int) {
	return file_proto_service_proto_rawDescGZIP(), []int{1, 0}
}

type LaunchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd string   `protobuf:"bytes,1,opt,name=cmd,proto3" json:"cmd,omitempty"`
	Arg []string `protobuf:"bytes,2,rep,name=arg,proto3" json:"arg,omitempty"`
}

func (x *LaunchRequest) Reset() {
	*x = LaunchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LaunchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LaunchRequest) ProtoMessage() {}

func (x *LaunchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LaunchRequest.ProtoReflect.Descriptor instead.
func (*LaunchRequest) Descriptor() ([]byte, []int) {
	return file_proto_service_proto_rawDescGZIP(), []int{0}
}

func (x *LaunchRequest) GetCmd() string {
	if x != nil {
		return x.Cmd
	}
	return ""
}

func (x *LaunchRequest) GetArg() []string {
	if x != nil {
		return x.Arg
	}
	return nil
}

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status   Status_StatusEnum `protobuf:"varint,1,opt,name=status,proto3,enum=main.Status_StatusEnum" json:"status,omitempty"`
	ExitCode int32             `protobuf:"varint,2,opt,name=ExitCode,proto3" json:"ExitCode,omitempty"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_proto_service_proto_rawDescGZIP(), []int{1}
}

func (x *Status) GetStatus() Status_StatusEnum {
	if x != nil {
		return x.Status
	}
	return Status_ERROR
}

func (x *Status) GetExitCode() int32 {
	if x != nil {
		return x.ExitCode
	}
	return 0
}

type Output struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Output) Reset() {
	*x = Output{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Output) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Output) ProtoMessage() {}

func (x *Output) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Output.ProtoReflect.Descriptor instead.
func (*Output) Descriptor() ([]byte, []int) {
	return file_proto_service_proto_rawDescGZIP(), []int{2}
}

func (x *Output) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Process struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobId string `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (x *Process) Reset() {
	*x = Process{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Process) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Process) ProtoMessage() {}

func (x *Process) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Process.ProtoReflect.Descriptor instead.
func (*Process) Descriptor() ([]byte, []int) {
	return file_proto_service_proto_rawDescGZIP(), []int{3}
}

func (x *Process) GetJobId() string {
	if x != nil {
		return x.JobId
	}
	return ""
}

type GeneralRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobId string `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (x *GeneralRequest) Reset() {
	*x = GeneralRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GeneralRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeneralRequest) ProtoMessage() {}

func (x *GeneralRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GeneralRequest.ProtoReflect.Descriptor instead.
func (*GeneralRequest) Descriptor() ([]byte, []int) {
	return file_proto_service_proto_rawDescGZIP(), []int{4}
}

func (x *GeneralRequest) GetJobId() string {
	if x != nil {
		return x.JobId
	}
	return ""
}

var File_proto_service_proto protoreflect.FileDescriptor

var file_proto_service_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x33, 0x0a, 0x0d, 0x4c,
	0x61, 0x75, 0x6e, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x10,
	0x0a, 0x03, 0x61, 0x72, 0x67, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x61, 0x72, 0x67,
	0x22, 0x97, 0x01, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x45, 0x6e, 0x75, 0x6d, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1a, 0x0a, 0x08,
	0x45, 0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x45, 0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x40, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10,
	0x00, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0d,
	0x0a, 0x09, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0b, 0x0a,
	0x07, 0x41, 0x42, 0x4f, 0x52, 0x54, 0x45, 0x44, 0x10, 0x03, 0x22, 0x1c, 0x0a, 0x06, 0x4f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x20, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x15, 0x0a, 0x06, 0x6a, 0x6f, 0x62, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x6a, 0x6f, 0x62, 0x49, 0x64, 0x22, 0x27, 0x0a, 0x0e, 0x47, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06,
	0x6a, 0x6f, 0x62, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6a, 0x6f,
	0x62, 0x49, 0x64, 0x32, 0xcf, 0x01, 0x0a, 0x06, 0x52, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x12, 0x2b,
	0x0a, 0x03, 0x52, 0x75, 0x6e, 0x12, 0x13, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4c, 0x61, 0x75,
	0x6e, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x6d, 0x61, 0x69,
	0x6e, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0c, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x14, 0x2e, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0c, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22,
	0x00, 0x30, 0x01, 0x12, 0x31, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x14, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x2d, 0x0a, 0x05, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x12,
	0x14, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x00, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x53, 0x61, 0x6e, 0x64, 0x79, 0x57, 0x61, 0x6c, 0x73, 0x68, 0x2f, 0x72,
	0x75, 0x6e, 0x6e, 0x65, 0x72, 0x2f, 0x73, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_service_proto_rawDescOnce sync.Once
	file_proto_service_proto_rawDescData = file_proto_service_proto_rawDesc
)

func file_proto_service_proto_rawDescGZIP() []byte {
	file_proto_service_proto_rawDescOnce.Do(func() {
		file_proto_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_service_proto_rawDescData)
	})
	return file_proto_service_proto_rawDescData
}

var file_proto_service_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_service_proto_goTypes = []interface{}{
	(Status_StatusEnum)(0), // 0: main.Status.StatusEnum
	(*LaunchRequest)(nil),  // 1: main.LaunchRequest
	(*Status)(nil),         // 2: main.Status
	(*Output)(nil),         // 3: main.Output
	(*Process)(nil),        // 4: main.Process
	(*GeneralRequest)(nil), // 5: main.GeneralRequest
}
var file_proto_service_proto_depIdxs = []int32{
	0, // 0: main.Status.status:type_name -> main.Status.StatusEnum
	1, // 1: main.Runner.Run:input_type -> main.LaunchRequest
	5, // 2: main.Runner.StreamOutput:input_type -> main.GeneralRequest
	5, // 3: main.Runner.GetStatus:input_type -> main.GeneralRequest
	5, // 4: main.Runner.Abort:input_type -> main.GeneralRequest
	4, // 5: main.Runner.Run:output_type -> main.Process
	3, // 6: main.Runner.StreamOutput:output_type -> main.Output
	2, // 7: main.Runner.GetStatus:output_type -> main.Status
	2, // 8: main.Runner.Abort:output_type -> main.Status
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_service_proto_init() }
func file_proto_service_proto_init() {
	if File_proto_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LaunchRequest); i {
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
		file_proto_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
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
		file_proto_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Output); i {
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
		file_proto_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Process); i {
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
		file_proto_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GeneralRequest); i {
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
			RawDescriptor: file_proto_service_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_service_proto_goTypes,
		DependencyIndexes: file_proto_service_proto_depIdxs,
		EnumInfos:         file_proto_service_proto_enumTypes,
		MessageInfos:      file_proto_service_proto_msgTypes,
	}.Build()
	File_proto_service_proto = out.File
	file_proto_service_proto_rawDesc = nil
	file_proto_service_proto_goTypes = nil
	file_proto_service_proto_depIdxs = nil
}
