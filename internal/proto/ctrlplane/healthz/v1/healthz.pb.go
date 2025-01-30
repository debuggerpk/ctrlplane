// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: ctrlplane/healthz/v1/healthz.proto

package healthzv1

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

type StatusResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Database      bool                   `protobuf:"varint,1,opt,name=database,proto3" json:"database,omitempty"`
	Temporal      bool                   `protobuf:"varint,2,opt,name=temporal,proto3" json:"temporal,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StatusResponse) Reset() {
	*x = StatusResponse{}
	mi := &file_ctrlplane_healthz_v1_healthz_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusResponse) ProtoMessage() {}

func (x *StatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_healthz_v1_healthz_proto_msgTypes[0]
	if x != nil {
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
	return file_ctrlplane_healthz_v1_healthz_proto_rawDescGZIP(), []int{0}
}

func (x *StatusResponse) GetDatabase() bool {
	if x != nil {
		return x.Database
	}
	return false
}

func (x *StatusResponse) GetTemporal() bool {
	if x != nil {
		return x.Temporal
	}
	return false
}

var File_ctrlplane_healthz_v1_healthz_proto protoreflect.FileDescriptor

var file_ctrlplane_healthz_v1_healthz_proto_rawDesc = string([]byte{
	0x0a, 0x22, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x68, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x7a, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x7a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e,
	0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x7a, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x48, 0x0a, 0x0e, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74,
	0x61, 0x62, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x61, 0x74,
	0x61, 0x62, 0x61, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61,
	0x6c, 0x32, 0x5c, 0x0a, 0x12, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x46, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x24, 0x2e, 0x63, 0x74, 0x72, 0x6c,
	0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x7a, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0xdb, 0x01, 0x0a, 0x18, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x7a, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x7a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3f, 0x67, 0x6f,
	0x2e, 0x62, 0x72, 0x65, 0x75, 0x2e, 0x69, 0x6f, 0x2f, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x6d, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63,
	0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x7a,
	0x2f, 0x76, 0x31, 0x3b, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x7a, 0x76, 0x31, 0xa2, 0x02, 0x03,
	0x43, 0x48, 0x58, 0xaa, 0x02, 0x14, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e,
	0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x7a, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x14, 0x43, 0x74, 0x72,
	0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x5c, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x7a, 0x5c, 0x56,
	0x31, 0xe2, 0x02, 0x20, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x5c, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x7a, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x16, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65,
	0x3a, 0x3a, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x7a, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_ctrlplane_healthz_v1_healthz_proto_rawDescOnce sync.Once
	file_ctrlplane_healthz_v1_healthz_proto_rawDescData []byte
)

func file_ctrlplane_healthz_v1_healthz_proto_rawDescGZIP() []byte {
	file_ctrlplane_healthz_v1_healthz_proto_rawDescOnce.Do(func() {
		file_ctrlplane_healthz_v1_healthz_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_ctrlplane_healthz_v1_healthz_proto_rawDesc), len(file_ctrlplane_healthz_v1_healthz_proto_rawDesc)))
	})
	return file_ctrlplane_healthz_v1_healthz_proto_rawDescData
}

var file_ctrlplane_healthz_v1_healthz_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ctrlplane_healthz_v1_healthz_proto_goTypes = []any{
	(*StatusResponse)(nil), // 0: ctrlplane.healthz.v1.StatusResponse
	(*emptypb.Empty)(nil),  // 1: google.protobuf.Empty
}
var file_ctrlplane_healthz_v1_healthz_proto_depIdxs = []int32{
	1, // 0: ctrlplane.healthz.v1.HealthCheckService.Status:input_type -> google.protobuf.Empty
	0, // 1: ctrlplane.healthz.v1.HealthCheckService.Status:output_type -> ctrlplane.healthz.v1.StatusResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ctrlplane_healthz_v1_healthz_proto_init() }
func file_ctrlplane_healthz_v1_healthz_proto_init() {
	if File_ctrlplane_healthz_v1_healthz_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_ctrlplane_healthz_v1_healthz_proto_rawDesc), len(file_ctrlplane_healthz_v1_healthz_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ctrlplane_healthz_v1_healthz_proto_goTypes,
		DependencyIndexes: file_ctrlplane_healthz_v1_healthz_proto_depIdxs,
		MessageInfos:      file_ctrlplane_healthz_v1_healthz_proto_msgTypes,
	}.Build()
	File_ctrlplane_healthz_v1_healthz_proto = out.File
	file_ctrlplane_healthz_v1_healthz_proto_goTypes = nil
	file_ctrlplane_healthz_v1_healthz_proto_depIdxs = nil
}
