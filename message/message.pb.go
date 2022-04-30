// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: protobuf/message.proto

package message

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

type Type int32

const (
	// Ack响应
	Type_Acknowledge Type = 0
	// 好友上线通知
	Type_FRIEND_ONLINE Type = 1
	// 好友下线通知
	Type_FRIEND_OFFLINE Type = 2
	// 好友消息
	Type_FRIEND_TEXT  Type = 10
	Type_FRIEND_IMAGE Type = 11
	Type_FRIEND_FILE  Type = 12
	// 好友申请通知
	Type_FRIEND_REQUEST Type = 20
	// 好友接受申请通知
	Type_FRIEND_ACCEPT Type = 21
	// 解除好友关系
	Type_FRIEND_DISBAND Type = 22
	// 小组信息
	Type_GROUP_TEXT  Type = 50
	Type_GROUP_IMAGE Type = 51
	Type_GROUP_FILE  Type = 52
	// 小组人员变动通知
	Type_GROUP_USER_CHANGE Type = 60
	// 加入小组申请通知
	Type_GROUP_REQUEST Type = 61
	// 接受小组申请通知
	Type_GROUP_ACCEPT Type = 62
	// 小组解散通知
	Type_GROUP_DISBAND Type = 63
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0:  "Acknowledge",
		1:  "FRIEND_ONLINE",
		2:  "FRIEND_OFFLINE",
		10: "FRIEND_TEXT",
		11: "FRIEND_IMAGE",
		12: "FRIEND_FILE",
		20: "FRIEND_REQUEST",
		21: "FRIEND_ACCEPT",
		22: "FRIEND_DISBAND",
		50: "GROUP_TEXT",
		51: "GROUP_IMAGE",
		52: "GROUP_FILE",
		60: "GROUP_USER_CHANGE",
		61: "GROUP_REQUEST",
		62: "GROUP_ACCEPT",
		63: "GROUP_DISBAND",
	}
	Type_value = map[string]int32{
		"Acknowledge":       0,
		"FRIEND_ONLINE":     1,
		"FRIEND_OFFLINE":    2,
		"FRIEND_TEXT":       10,
		"FRIEND_IMAGE":      11,
		"FRIEND_FILE":       12,
		"FRIEND_REQUEST":    20,
		"FRIEND_ACCEPT":     21,
		"FRIEND_DISBAND":    22,
		"GROUP_TEXT":        50,
		"GROUP_IMAGE":       51,
		"GROUP_FILE":        52,
		"GROUP_USER_CHANGE": 60,
		"GROUP_REQUEST":     61,
		"GROUP_ACCEPT":      62,
		"GROUP_DISBAND":     63,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_protobuf_message_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_protobuf_message_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_protobuf_message_proto_rawDescGZIP(), []int{0}
}

type State int32

const (
	// 服务器收到消息
	State_SERVER_RECV State = 0
	// 服务器已转发消息
	State_WAIT_ACK State = 1
	// 客户端收到消息
	State_CLIENT_RECV State = 2
)

// Enum value maps for State.
var (
	State_name = map[int32]string{
		0: "SERVER_RECV",
		1: "WAIT_ACK",
		2: "CLIENT_RECV",
	}
	State_value = map[string]int32{
		"SERVER_RECV": 0,
		"WAIT_ACK":    1,
		"CLIENT_RECV": 2,
	}
)

func (x State) Enum() *State {
	p := new(State)
	*p = x
	return p
}

func (x State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (State) Descriptor() protoreflect.EnumDescriptor {
	return file_protobuf_message_proto_enumTypes[1].Descriptor()
}

func (State) Type() protoreflect.EnumType {
	return &file_protobuf_message_proto_enumTypes[1]
}

func (x State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use State.Descriptor instead.
func (State) EnumDescriptor() ([]byte, []int) {
	return file_protobuf_message_proto_rawDescGZIP(), []int{1}
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	LocalId int64  `protobuf:"varint,2,opt,name=local_id,json=localId,proto3" json:"local_id,omitempty"`
	Type    Type   `protobuf:"varint,3,opt,name=type,proto3,enum=Type" json:"type,omitempty"`
	State   State  `protobuf:"varint,4,opt,name=state,proto3,enum=State" json:"state,omitempty"`
	From    uint32 `protobuf:"varint,5,opt,name=from,proto3" json:"from,omitempty"`
	To      uint32 `protobuf:"varint,6,opt,name=to,proto3" json:"to,omitempty"`
	Content []byte `protobuf:"bytes,7,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_protobuf_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Message) GetLocalId() int64 {
	if x != nil {
		return x.LocalId
	}
	return 0
}

func (x *Message) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_Acknowledge
}

func (x *Message) GetState() State {
	if x != nil {
		return x.State
	}
	return State_SERVER_RECV
}

func (x *Message) GetFrom() uint32 {
	if x != nil {
		return x.From
	}
	return 0
}

func (x *Message) GetTo() uint32 {
	if x != nil {
		return x.To
	}
	return 0
}

func (x *Message) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type FileContent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key  []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *FileContent) Reset() {
	*x = FileContent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileContent) ProtoMessage() {}

func (x *FileContent) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileContent.ProtoReflect.Descriptor instead.
func (*FileContent) Descriptor() ([]byte, []int) {
	return file_protobuf_message_proto_rawDescGZIP(), []int{1}
}

func (x *FileContent) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *FileContent) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_protobuf_message_proto protoreflect.FileDescriptor

var file_protobuf_message_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xab, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x49, 0x64, 0x12,
	0x19, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x05, 0x2e,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x06, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02,
	0x74, 0x6f, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x18, 0x0a, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x33, 0x0a, 0x0b, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x2a, 0xad, 0x02, 0x0a, 0x04,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x41, 0x63, 0x6b, 0x6e, 0x6f, 0x77, 0x6c, 0x65,
	0x64, 0x67, 0x65, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x46, 0x52, 0x49, 0x45, 0x4e, 0x44, 0x5f,
	0x4f, 0x4e, 0x4c, 0x49, 0x4e, 0x45, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x46, 0x52, 0x49, 0x45,
	0x4e, 0x44, 0x5f, 0x4f, 0x46, 0x46, 0x4c, 0x49, 0x4e, 0x45, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b,
	0x46, 0x52, 0x49, 0x45, 0x4e, 0x44, 0x5f, 0x54, 0x45, 0x58, 0x54, 0x10, 0x0a, 0x12, 0x10, 0x0a,
	0x0c, 0x46, 0x52, 0x49, 0x45, 0x4e, 0x44, 0x5f, 0x49, 0x4d, 0x41, 0x47, 0x45, 0x10, 0x0b, 0x12,
	0x0f, 0x0a, 0x0b, 0x46, 0x52, 0x49, 0x45, 0x4e, 0x44, 0x5f, 0x46, 0x49, 0x4c, 0x45, 0x10, 0x0c,
	0x12, 0x12, 0x0a, 0x0e, 0x46, 0x52, 0x49, 0x45, 0x4e, 0x44, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45,
	0x53, 0x54, 0x10, 0x14, 0x12, 0x11, 0x0a, 0x0d, 0x46, 0x52, 0x49, 0x45, 0x4e, 0x44, 0x5f, 0x41,
	0x43, 0x43, 0x45, 0x50, 0x54, 0x10, 0x15, 0x12, 0x12, 0x0a, 0x0e, 0x46, 0x52, 0x49, 0x45, 0x4e,
	0x44, 0x5f, 0x44, 0x49, 0x53, 0x42, 0x41, 0x4e, 0x44, 0x10, 0x16, 0x12, 0x0e, 0x0a, 0x0a, 0x47,
	0x52, 0x4f, 0x55, 0x50, 0x5f, 0x54, 0x45, 0x58, 0x54, 0x10, 0x32, 0x12, 0x0f, 0x0a, 0x0b, 0x47,
	0x52, 0x4f, 0x55, 0x50, 0x5f, 0x49, 0x4d, 0x41, 0x47, 0x45, 0x10, 0x33, 0x12, 0x0e, 0x0a, 0x0a,
	0x47, 0x52, 0x4f, 0x55, 0x50, 0x5f, 0x46, 0x49, 0x4c, 0x45, 0x10, 0x34, 0x12, 0x15, 0x0a, 0x11,
	0x47, 0x52, 0x4f, 0x55, 0x50, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x43, 0x48, 0x41, 0x4e, 0x47,
	0x45, 0x10, 0x3c, 0x12, 0x11, 0x0a, 0x0d, 0x47, 0x52, 0x4f, 0x55, 0x50, 0x5f, 0x52, 0x45, 0x51,
	0x55, 0x45, 0x53, 0x54, 0x10, 0x3d, 0x12, 0x10, 0x0a, 0x0c, 0x47, 0x52, 0x4f, 0x55, 0x50, 0x5f,
	0x41, 0x43, 0x43, 0x45, 0x50, 0x54, 0x10, 0x3e, 0x12, 0x11, 0x0a, 0x0d, 0x47, 0x52, 0x4f, 0x55,
	0x50, 0x5f, 0x44, 0x49, 0x53, 0x42, 0x41, 0x4e, 0x44, 0x10, 0x3f, 0x2a, 0x37, 0x0a, 0x05, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52, 0x5f, 0x52,
	0x45, 0x43, 0x56, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x57, 0x41, 0x49, 0x54, 0x5f, 0x41, 0x43,
	0x4b, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x4c, 0x49, 0x45, 0x4e, 0x54, 0x5f, 0x52, 0x45,
	0x43, 0x56, 0x10, 0x02, 0x42, 0x0a, 0x5a, 0x08, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protobuf_message_proto_rawDescOnce sync.Once
	file_protobuf_message_proto_rawDescData = file_protobuf_message_proto_rawDesc
)

func file_protobuf_message_proto_rawDescGZIP() []byte {
	file_protobuf_message_proto_rawDescOnce.Do(func() {
		file_protobuf_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_message_proto_rawDescData)
	})
	return file_protobuf_message_proto_rawDescData
}

var file_protobuf_message_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_protobuf_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protobuf_message_proto_goTypes = []interface{}{
	(Type)(0),           // 0: Type
	(State)(0),          // 1: State
	(*Message)(nil),     // 2: Message
	(*FileContent)(nil), // 3: FileContent
}
var file_protobuf_message_proto_depIdxs = []int32{
	0, // 0: Message.type:type_name -> Type
	1, // 1: Message.state:type_name -> State
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_protobuf_message_proto_init() }
func file_protobuf_message_proto_init() {
	if File_protobuf_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protobuf_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_protobuf_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileContent); i {
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
			RawDescriptor: file_protobuf_message_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protobuf_message_proto_goTypes,
		DependencyIndexes: file_protobuf_message_proto_depIdxs,
		EnumInfos:         file_protobuf_message_proto_enumTypes,
		MessageInfos:      file_protobuf_message_proto_msgTypes,
	}.Build()
	File_protobuf_message_proto = out.File
	file_protobuf_message_proto_rawDesc = nil
	file_protobuf_message_proto_goTypes = nil
	file_protobuf_message_proto_depIdxs = nil
}
