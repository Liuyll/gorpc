// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.3.0
// source: responseHeader.proto

package protocol

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type RPCResponseHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  string `protobuf:"bytes,1,opt,name=Type,proto3" json:"Type,omitempty"`
	Seq   int64  `protobuf:"varint,2,opt,name=Seq,proto3" json:"Seq,omitempty"`
	Error string `protobuf:"bytes,3,opt,name=Error,proto3" json:"Error,omitempty"`
	Reply []byte `protobuf:"bytes,4,opt,name=Reply,proto3" json:"Reply,omitempty"`
}

func (x *RPCResponseHeader) Reset() {
	*x = RPCResponseHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_responseHeader_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RPCResponseHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RPCResponseHeader) ProtoMessage() {}

func (x *RPCResponseHeader) ProtoReflect() protoreflect.Message {
	mi := &file_responseHeader_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RPCResponseHeader.ProtoReflect.Descriptor instead.
func (*RPCResponseHeader) Descriptor() ([]byte, []int) {
	return file_responseHeader_proto_rawDescGZIP(), []int{0}
}

func (x *RPCResponseHeader) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *RPCResponseHeader) GetSeq() int64 {
	if x != nil {
		return x.Seq
	}
	return 0
}

func (x *RPCResponseHeader) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *RPCResponseHeader) GetReply() []byte {
	if x != nil {
		return x.Reply
	}
	return nil
}

var File_responseHeader_proto protoreflect.FileDescriptor

var file_responseHeader_proto_rawDesc = []byte{
	0x0a, 0x14, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x65, 0x0a, 0x11, 0x52, 0x50, 0x43, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x53, 0x65, 0x71, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x53, 0x65,
	0x71, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x0b, 0x5a,
	0x09, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_responseHeader_proto_rawDescOnce sync.Once
	file_responseHeader_proto_rawDescData = file_responseHeader_proto_rawDesc
)

func file_responseHeader_proto_rawDescGZIP() []byte {
	file_responseHeader_proto_rawDescOnce.Do(func() {
		file_responseHeader_proto_rawDescData = protoimpl.X.CompressGZIP(file_responseHeader_proto_rawDescData)
	})
	return file_responseHeader_proto_rawDescData
}

var file_responseHeader_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_responseHeader_proto_goTypes = []interface{}{
	(*RPCResponseHeader)(nil), // 0: RPCResponseHeader
}
var file_responseHeader_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_responseHeader_proto_init() }
func file_responseHeader_proto_init() {
	if File_responseHeader_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_responseHeader_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RPCResponseHeader); i {
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
			RawDescriptor: file_responseHeader_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_responseHeader_proto_goTypes,
		DependencyIndexes: file_responseHeader_proto_depIdxs,
		MessageInfos:      file_responseHeader_proto_msgTypes,
	}.Build()
	File_responseHeader_proto = out.File
	file_responseHeader_proto_rawDesc = nil
	file_responseHeader_proto_goTypes = nil
	file_responseHeader_proto_depIdxs = nil
}