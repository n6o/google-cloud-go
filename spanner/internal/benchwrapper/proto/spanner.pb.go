package spanner_bench

import (
	reflect "reflect"
	sync "sync"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	"log"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Singer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	Id            int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	FirstName     string `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName      string `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	SingerInfo    string `protobuf:"bytes,4,opt,name=singer_info,json=singerInfo,proto3" json:"singer_info,omitempty"`
}

func (x *Singer) gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2() {
	*x = Singer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spanner_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}
func (x *Singer) gologoo__String_ddb42bbeb083878c2bfa75281902bae2() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Singer) gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2() {
}
func (x *Singer) gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2() protoreflect.Message {
	mi := &file_spanner_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Singer) gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2() ([]byte, []int) {
	return file_spanner_proto_rawDescGZIP(), []int{0}
}
func (x *Singer) gologoo__GetId_ddb42bbeb083878c2bfa75281902bae2() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}
func (x *Singer) gologoo__GetFirstName_ddb42bbeb083878c2bfa75281902bae2() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}
func (x *Singer) gologoo__GetLastName_ddb42bbeb083878c2bfa75281902bae2() string {
	if x != nil {
		return x.LastName
	}
	return ""
}
func (x *Singer) gologoo__GetSingerInfo_ddb42bbeb083878c2bfa75281902bae2() string {
	if x != nil {
		return x.SingerInfo
	}
	return ""
}

type Album struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	Id            int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	SingerId      int64  `protobuf:"varint,2,opt,name=singer_id,json=singerId,proto3" json:"singer_id,omitempty"`
	AlbumTitle    string `protobuf:"bytes,3,opt,name=album_title,json=albumTitle,proto3" json:"album_title,omitempty"`
}

func (x *Album) gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2() {
	*x = Album{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spanner_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}
func (x *Album) gologoo__String_ddb42bbeb083878c2bfa75281902bae2() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*Album) gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2() {
}
func (x *Album) gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2() protoreflect.Message {
	mi := &file_spanner_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*Album) gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2() ([]byte, []int) {
	return file_spanner_proto_rawDescGZIP(), []int{1}
}
func (x *Album) gologoo__GetId_ddb42bbeb083878c2bfa75281902bae2() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}
func (x *Album) gologoo__GetSingerId_ddb42bbeb083878c2bfa75281902bae2() int64 {
	if x != nil {
		return x.SingerId
	}
	return 0
}
func (x *Album) gologoo__GetAlbumTitle_ddb42bbeb083878c2bfa75281902bae2() string {
	if x != nil {
		return x.AlbumTitle
	}
	return ""
}

type ReadQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	Query         string `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
}

func (x *ReadQuery) gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2() {
	*x = ReadQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spanner_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}
func (x *ReadQuery) gologoo__String_ddb42bbeb083878c2bfa75281902bae2() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*ReadQuery) gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2() {
}
func (x *ReadQuery) gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2() protoreflect.Message {
	mi := &file_spanner_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*ReadQuery) gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2() ([]byte, []int) {
	return file_spanner_proto_rawDescGZIP(), []int{2}
}
func (x *ReadQuery) gologoo__GetQuery_ddb42bbeb083878c2bfa75281902bae2() string {
	if x != nil {
		return x.Query
	}
	return ""
}

type InsertQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	Singers       []*Singer `protobuf:"bytes,1,rep,name=singers,proto3" json:"singers,omitempty"`
	Albums        []*Album  `protobuf:"bytes,2,rep,name=albums,proto3" json:"albums,omitempty"`
}

func (x *InsertQuery) gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2() {
	*x = InsertQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spanner_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}
func (x *InsertQuery) gologoo__String_ddb42bbeb083878c2bfa75281902bae2() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*InsertQuery) gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2() {
}
func (x *InsertQuery) gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2() protoreflect.Message {
	mi := &file_spanner_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*InsertQuery) gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2() ([]byte, []int) {
	return file_spanner_proto_rawDescGZIP(), []int{3}
}
func (x *InsertQuery) gologoo__GetSingers_ddb42bbeb083878c2bfa75281902bae2() []*Singer {
	if x != nil {
		return x.Singers
	}
	return nil
}
func (x *InsertQuery) gologoo__GetAlbums_ddb42bbeb083878c2bfa75281902bae2() []*Album {
	if x != nil {
		return x.Albums
	}
	return nil
}

type UpdateQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	Queries       []string `protobuf:"bytes,1,rep,name=queries,proto3" json:"queries,omitempty"`
}

func (x *UpdateQuery) gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2() {
	*x = UpdateQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spanner_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}
func (x *UpdateQuery) gologoo__String_ddb42bbeb083878c2bfa75281902bae2() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*UpdateQuery) gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2() {
}
func (x *UpdateQuery) gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2() protoreflect.Message {
	mi := &file_spanner_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*UpdateQuery) gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2() ([]byte, []int) {
	return file_spanner_proto_rawDescGZIP(), []int{4}
}
func (x *UpdateQuery) gologoo__GetQueries_ddb42bbeb083878c2bfa75281902bae2() []string {
	if x != nil {
		return x.Queries
	}
	return nil
}

type EmptyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyResponse) gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2() {
	*x = EmptyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spanner_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}
func (x *EmptyResponse) gologoo__String_ddb42bbeb083878c2bfa75281902bae2() string {
	return protoimpl.X.MessageStringOf(x)
}
func (*EmptyResponse) gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2() {
}
func (x *EmptyResponse) gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2() protoreflect.Message {
	mi := &file_spanner_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}
func (*EmptyResponse) gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2() ([]byte, []int) {
	return file_spanner_proto_rawDescGZIP(), []int{5}
}

var File_spanner_proto protoreflect.FileDescriptor
var file_spanner_proto_rawDesc = []byte{0x0a, 0x0d, 0x73, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x73, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x22, 0x75, 0x0a, 0x06, 0x53, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x55, 0x0a, 0x05, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x5f, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x22, 0x21, 0x0a, 0x09, 0x52, 0x65, 0x61, 0x64, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x22, 0x6c, 0x0a, 0x0b, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x2f, 0x0a, 0x07, 0x73, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x52, 0x07, 0x73, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x73, 0x12, 0x2c, 0x0a, 0x06, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x2e, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x52, 0x06, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x22, 0x27, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x71, 0x75, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x71, 0x75, 0x65, 0x72, 0x69, 0x65, 0x73, 0x22, 0x0f, 0x0a, 0x0d, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xe3, 0x01, 0x0a, 0x13, 0x53, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x40, 0x0a, 0x04, 0x52, 0x65, 0x61, 0x64, 0x12, 0x18, 0x2e, 0x73, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x51, 0x75, 0x65, 0x72, 0x79, 0x1a, 0x1c, 0x2e, 0x73, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x06, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x1a, 0x2e, 0x73, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x2e, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x51, 0x75, 0x65, 0x72, 0x79, 0x1a, 0x1c, 0x2e, 0x73, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1a, 0x2e, 0x73, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x1a, 0x1c, 0x2e, 0x73, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x10, 0x5a, 0x0e, 0x73, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33}
var (
	file_spanner_proto_rawDescOnce sync.Once
	file_spanner_proto_rawDescData = file_spanner_proto_rawDesc
)

func gologoo__file_spanner_proto_rawDescGZIP_ddb42bbeb083878c2bfa75281902bae2() []byte {
	file_spanner_proto_rawDescOnce.Do(func() {
		file_spanner_proto_rawDescData = protoimpl.X.CompressGZIP(file_spanner_proto_rawDescData)
	})
	return file_spanner_proto_rawDescData
}

var file_spanner_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_spanner_proto_goTypes = []interface {
}{(*Singer)(nil), (*Album)(nil), (*ReadQuery)(nil), (*InsertQuery)(nil), (*UpdateQuery)(nil), (*EmptyResponse)(nil)}
var file_spanner_proto_depIdxs = []int32{0, 1, 2, 3, 4, 5, 5, 5, 5, 2, 2, 2, 0}

func gologoo__init_ddb42bbeb083878c2bfa75281902bae2() {
	file_spanner_proto_init()
}
func gologoo__file_spanner_proto_init_ddb42bbeb083878c2bfa75281902bae2() {
	if File_spanner_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_spanner_proto_msgTypes[0].Exporter = func(v interface {
		}, i int) interface {
		} {
			switch v := v.(*Singer); i {
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
		file_spanner_proto_msgTypes[1].Exporter = func(v interface {
		}, i int) interface {
		} {
			switch v := v.(*Album); i {
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
		file_spanner_proto_msgTypes[2].Exporter = func(v interface {
		}, i int) interface {
		} {
			switch v := v.(*ReadQuery); i {
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
		file_spanner_proto_msgTypes[3].Exporter = func(v interface {
		}, i int) interface {
		} {
			switch v := v.(*InsertQuery); i {
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
		file_spanner_proto_msgTypes[4].Exporter = func(v interface {
		}, i int) interface {
		} {
			switch v := v.(*UpdateQuery); i {
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
		file_spanner_proto_msgTypes[5].Exporter = func(v interface {
		}, i int) interface {
		} {
			switch v := v.(*EmptyResponse); i {
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
	type x struct {
	}
	out := protoimpl.TypeBuilder{File: protoimpl.DescBuilder{GoPackagePath: reflect.TypeOf(x{}).PkgPath(), RawDescriptor: file_spanner_proto_rawDesc, NumEnums: 0, NumMessages: 6, NumExtensions: 0, NumServices: 1}, GoTypes: file_spanner_proto_goTypes, DependencyIndexes: file_spanner_proto_depIdxs, MessageInfos: file_spanner_proto_msgTypes}.Build()
	File_spanner_proto = out.File
	file_spanner_proto_rawDesc = nil
	file_spanner_proto_goTypes = nil
	file_spanner_proto_depIdxs = nil
}
func (x *Singer) Reset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	x.gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (x *Singer) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__String_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *Singer) ProtoMessage() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	recv.gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (x *Singer) ProtoReflect() protoreflect.Message {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *Singer) Descriptor() ([]byte, []int) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0, r1 := recv.gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (x *Singer) GetId() int64 {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetId_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__GetId_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (x *Singer) GetFirstName() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetFirstName_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__GetFirstName_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (x *Singer) GetLastName() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetLastName_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__GetLastName_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (x *Singer) GetSingerInfo() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetSingerInfo_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__GetSingerInfo_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (x *Album) Reset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	x.gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (x *Album) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__String_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *Album) ProtoMessage() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	recv.gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (x *Album) ProtoReflect() protoreflect.Message {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *Album) Descriptor() ([]byte, []int) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0, r1 := recv.gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (x *Album) GetId() int64 {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetId_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__GetId_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (x *Album) GetSingerId() int64 {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetSingerId_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__GetSingerId_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (x *Album) GetAlbumTitle() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetAlbumTitle_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__GetAlbumTitle_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (x *ReadQuery) Reset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	x.gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (x *ReadQuery) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__String_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *ReadQuery) ProtoMessage() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	recv.gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (x *ReadQuery) ProtoReflect() protoreflect.Message {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *ReadQuery) Descriptor() ([]byte, []int) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0, r1 := recv.gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (x *ReadQuery) GetQuery() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetQuery_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__GetQuery_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (x *InsertQuery) Reset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	x.gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (x *InsertQuery) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__String_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *InsertQuery) ProtoMessage() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	recv.gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (x *InsertQuery) ProtoReflect() protoreflect.Message {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *InsertQuery) Descriptor() ([]byte, []int) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0, r1 := recv.gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (x *InsertQuery) GetSingers() []*Singer {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetSingers_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__GetSingers_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (x *InsertQuery) GetAlbums() []*Album {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetAlbums_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__GetAlbums_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (x *UpdateQuery) Reset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	x.gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (x *UpdateQuery) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__String_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *UpdateQuery) ProtoMessage() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	recv.gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (x *UpdateQuery) ProtoReflect() protoreflect.Message {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *UpdateQuery) Descriptor() ([]byte, []int) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0, r1 := recv.gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (x *UpdateQuery) GetQueries() []string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetQueries_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__GetQueries_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (x *EmptyResponse) Reset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	x.gologoo__Reset_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (x *EmptyResponse) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__String_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *EmptyResponse) ProtoMessage() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	recv.gologoo__ProtoMessage_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (x *EmptyResponse) ProtoReflect() protoreflect.Message {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := x.gologoo__ProtoReflect_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *EmptyResponse) Descriptor() ([]byte, []int) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0, r1 := recv.gologoo__Descriptor_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func file_spanner_proto_rawDescGZIP() []byte {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__file_spanner_proto_rawDescGZIP_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	r0 := gologoo__file_spanner_proto_rawDescGZIP_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func init() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__init_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	gologoo__init_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func file_spanner_proto_init() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__file_spanner_proto_init_ddb42bbeb083878c2bfa75281902bae2")
	log.Printf("Input : (none)\n")
	gologoo__file_spanner_proto_init_ddb42bbeb083878c2bfa75281902bae2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
