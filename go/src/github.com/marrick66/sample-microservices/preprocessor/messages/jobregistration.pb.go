// Code generated by protoc-gen-go. DO NOT EDIT.
// source: jobregistration.proto

package messages

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GetRegistrationReply_Status int32

const (
	GetRegistrationReply_NOTFOUND   GetRegistrationReply_Status = 0
	GetRegistrationReply_REGISTERED GetRegistrationReply_Status = 1
	GetRegistrationReply_RUNNING    GetRegistrationReply_Status = 2
	GetRegistrationReply_FAILED     GetRegistrationReply_Status = 3
	GetRegistrationReply_COMPLETED  GetRegistrationReply_Status = 4
)

var GetRegistrationReply_Status_name = map[int32]string{
	0: "NOTFOUND",
	1: "REGISTERED",
	2: "RUNNING",
	3: "FAILED",
	4: "COMPLETED",
}

var GetRegistrationReply_Status_value = map[string]int32{
	"NOTFOUND":   0,
	"REGISTERED": 1,
	"RUNNING":    2,
	"FAILED":     3,
	"COMPLETED":  4,
}

func (x GetRegistrationReply_Status) String() string {
	return proto.EnumName(GetRegistrationReply_Status_name, int32(x))
}

func (GetRegistrationReply_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_30037a16df18364b, []int{3, 0}
}

//The registration request message:
type RegistrationRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegistrationRequest) Reset()         { *m = RegistrationRequest{} }
func (m *RegistrationRequest) String() string { return proto.CompactTextString(m) }
func (*RegistrationRequest) ProtoMessage()    {}
func (*RegistrationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_30037a16df18364b, []int{0}
}

func (m *RegistrationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegistrationRequest.Unmarshal(m, b)
}
func (m *RegistrationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegistrationRequest.Marshal(b, m, deterministic)
}
func (m *RegistrationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegistrationRequest.Merge(m, src)
}
func (m *RegistrationRequest) XXX_Size() int {
	return xxx_messageInfo_RegistrationRequest.Size(m)
}
func (m *RegistrationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegistrationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegistrationRequest proto.InternalMessageInfo

func (m *RegistrationRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

//The registration request reply message:
type RegistrationReply struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegistrationReply) Reset()         { *m = RegistrationReply{} }
func (m *RegistrationReply) String() string { return proto.CompactTextString(m) }
func (*RegistrationReply) ProtoMessage()    {}
func (*RegistrationReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_30037a16df18364b, []int{1}
}

func (m *RegistrationReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegistrationReply.Unmarshal(m, b)
}
func (m *RegistrationReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegistrationReply.Marshal(b, m, deterministic)
}
func (m *RegistrationReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegistrationReply.Merge(m, src)
}
func (m *RegistrationReply) XXX_Size() int {
	return xxx_messageInfo_RegistrationReply.Size(m)
}
func (m *RegistrationReply) XXX_DiscardUnknown() {
	xxx_messageInfo_RegistrationReply.DiscardUnknown(m)
}

var xxx_messageInfo_RegistrationReply proto.InternalMessageInfo

func (m *RegistrationReply) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

//The get registration request message:
type GetRegistrationRequest struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRegistrationRequest) Reset()         { *m = GetRegistrationRequest{} }
func (m *GetRegistrationRequest) String() string { return proto.CompactTextString(m) }
func (*GetRegistrationRequest) ProtoMessage()    {}
func (*GetRegistrationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_30037a16df18364b, []int{2}
}

func (m *GetRegistrationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRegistrationRequest.Unmarshal(m, b)
}
func (m *GetRegistrationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRegistrationRequest.Marshal(b, m, deterministic)
}
func (m *GetRegistrationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRegistrationRequest.Merge(m, src)
}
func (m *GetRegistrationRequest) XXX_Size() int {
	return xxx_messageInfo_GetRegistrationRequest.Size(m)
}
func (m *GetRegistrationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRegistrationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRegistrationRequest proto.InternalMessageInfo

func (m *GetRegistrationRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

//The get registration reply message:
type GetRegistrationReply struct {
	Id                   int32                       `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Status               GetRegistrationReply_Status `protobuf:"varint,2,opt,name=status,proto3,enum=jobregistration.GetRegistrationReply_Status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *GetRegistrationReply) Reset()         { *m = GetRegistrationReply{} }
func (m *GetRegistrationReply) String() string { return proto.CompactTextString(m) }
func (*GetRegistrationReply) ProtoMessage()    {}
func (*GetRegistrationReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_30037a16df18364b, []int{3}
}

func (m *GetRegistrationReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRegistrationReply.Unmarshal(m, b)
}
func (m *GetRegistrationReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRegistrationReply.Marshal(b, m, deterministic)
}
func (m *GetRegistrationReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRegistrationReply.Merge(m, src)
}
func (m *GetRegistrationReply) XXX_Size() int {
	return xxx_messageInfo_GetRegistrationReply.Size(m)
}
func (m *GetRegistrationReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRegistrationReply.DiscardUnknown(m)
}

var xxx_messageInfo_GetRegistrationReply proto.InternalMessageInfo

func (m *GetRegistrationReply) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *GetRegistrationReply) GetStatus() GetRegistrationReply_Status {
	if m != nil {
		return m.Status
	}
	return GetRegistrationReply_NOTFOUND
}

func init() {
	proto.RegisterEnum("jobregistration.GetRegistrationReply_Status", GetRegistrationReply_Status_name, GetRegistrationReply_Status_value)
	proto.RegisterType((*RegistrationRequest)(nil), "jobregistration.RegistrationRequest")
	proto.RegisterType((*RegistrationReply)(nil), "jobregistration.RegistrationReply")
	proto.RegisterType((*GetRegistrationRequest)(nil), "jobregistration.GetRegistrationRequest")
	proto.RegisterType((*GetRegistrationReply)(nil), "jobregistration.GetRegistrationReply")
}

func init() { proto.RegisterFile("jobregistration.proto", fileDescriptor_30037a16df18364b) }

var fileDescriptor_30037a16df18364b = []byte{
	// 287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xcf, 0x4e, 0x83, 0x40,
	0x10, 0xc6, 0x59, 0xac, 0xd8, 0x8e, 0x0a, 0x38, 0xfe, 0x49, 0xe3, 0xa9, 0x59, 0x35, 0x62, 0x62,
	0x38, 0xd4, 0x27, 0x30, 0xb2, 0x25, 0x98, 0xba, 0x98, 0x2d, 0xf5, 0x0e, 0x76, 0x63, 0x30, 0xb5,
	0x20, 0x6c, 0x0f, 0x7d, 0x37, 0x9f, 0xc1, 0x67, 0x32, 0x01, 0x0f, 0x95, 0x92, 0xd4, 0xdb, 0xce,
	0xe4, 0x37, 0xdf, 0xec, 0xf7, 0x65, 0xe0, 0xf4, 0x3d, 0x4b, 0x0a, 0xf9, 0x96, 0x96, 0xaa, 0x88,
	0x55, 0x9a, 0x2d, 0xdc, 0xbc, 0xc8, 0x54, 0x86, 0x56, 0xa3, 0x4d, 0x6f, 0xe0, 0x58, 0xac, 0xd5,
	0x42, 0x7e, 0x2e, 0x65, 0xa9, 0x10, 0xa1, 0xb3, 0x88, 0x3f, 0x64, 0x9f, 0x0c, 0x88, 0xd3, 0x13,
	0xd5, 0x9b, 0x5e, 0xc0, 0xd1, 0x5f, 0x34, 0x9f, 0xaf, 0xd0, 0x04, 0x3d, 0x9d, 0x55, 0xd8, 0xae,
	0xd0, 0xd3, 0x19, 0x75, 0xe0, 0xcc, 0x97, 0xaa, 0x4d, 0xb2, 0x49, 0x7e, 0x11, 0x38, 0xd9, 0x40,
	0x5b, 0x24, 0xd1, 0x03, 0xa3, 0x54, 0xb1, 0x5a, 0x96, 0x7d, 0x7d, 0x40, 0x1c, 0x73, 0x78, 0xeb,
	0x36, 0xbd, 0xb5, 0xc9, 0xb8, 0x93, 0x6a, 0x46, 0xfc, 0xce, 0x52, 0x0e, 0x46, 0xdd, 0xc1, 0x03,
	0xe8, 0xf2, 0x30, 0x1a, 0x85, 0x53, 0xee, 0xd9, 0x1a, 0x9a, 0x00, 0x82, 0xf9, 0xc1, 0x24, 0x62,
	0x82, 0x79, 0x36, 0xc1, 0x7d, 0xd8, 0x13, 0x53, 0xce, 0x03, 0xee, 0xdb, 0x3a, 0x02, 0x18, 0xa3,
	0xfb, 0x60, 0xcc, 0x3c, 0x7b, 0x07, 0x0f, 0xa1, 0xf7, 0x10, 0x3e, 0x3d, 0x8f, 0x59, 0xc4, 0x3c,
	0xbb, 0x33, 0xfc, 0x26, 0x60, 0x3d, 0x66, 0xc9, 0xfa, 0x5e, 0x7c, 0x81, 0x6e, 0x5d, 0xcb, 0x02,
	0x2f, 0x37, 0x7e, 0xd9, 0x12, 0xca, 0x39, 0xdd, 0x42, 0xe5, 0xf3, 0x15, 0xd5, 0xf0, 0x15, 0xac,
	0x86, 0x45, 0xbc, 0xde, 0x1e, 0x42, 0xbd, 0xe1, 0xea, 0x5f, 0x69, 0x51, 0x2d, 0x31, 0xaa, 0x0b,
	0xb9, 0xfb, 0x09, 0x00, 0x00, 0xff, 0xff, 0x70, 0x15, 0x3a, 0x5c, 0x3a, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// JobRegistrationClient is the client API for JobRegistration service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type JobRegistrationClient interface {
	//Registers a job for execution.
	Register(ctx context.Context, in *RegistrationRequest, opts ...grpc.CallOption) (*RegistrationReply, error)
	GetRegistration(ctx context.Context, in *GetRegistrationRequest, opts ...grpc.CallOption) (*GetRegistrationReply, error)
}

type jobRegistrationClient struct {
	cc *grpc.ClientConn
}

func NewJobRegistrationClient(cc *grpc.ClientConn) JobRegistrationClient {
	return &jobRegistrationClient{cc}
}

func (c *jobRegistrationClient) Register(ctx context.Context, in *RegistrationRequest, opts ...grpc.CallOption) (*RegistrationReply, error) {
	out := new(RegistrationReply)
	err := c.cc.Invoke(ctx, "/jobregistration.JobRegistration/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobRegistrationClient) GetRegistration(ctx context.Context, in *GetRegistrationRequest, opts ...grpc.CallOption) (*GetRegistrationReply, error) {
	out := new(GetRegistrationReply)
	err := c.cc.Invoke(ctx, "/jobregistration.JobRegistration/GetRegistration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JobRegistrationServer is the server API for JobRegistration service.
type JobRegistrationServer interface {
	//Registers a job for execution.
	Register(context.Context, *RegistrationRequest) (*RegistrationReply, error)
	GetRegistration(context.Context, *GetRegistrationRequest) (*GetRegistrationReply, error)
}

// UnimplementedJobRegistrationServer can be embedded to have forward compatible implementations.
type UnimplementedJobRegistrationServer struct {
}

func (*UnimplementedJobRegistrationServer) Register(ctx context.Context, req *RegistrationRequest) (*RegistrationReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (*UnimplementedJobRegistrationServer) GetRegistration(ctx context.Context, req *GetRegistrationRequest) (*GetRegistrationReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRegistration not implemented")
}

func RegisterJobRegistrationServer(s *grpc.Server, srv JobRegistrationServer) {
	s.RegisterService(&_JobRegistration_serviceDesc, srv)
}

func _JobRegistration_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobRegistrationServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobregistration.JobRegistration/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobRegistrationServer).Register(ctx, req.(*RegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobRegistration_GetRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobRegistrationServer).GetRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobregistration.JobRegistration/GetRegistration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobRegistrationServer).GetRegistration(ctx, req.(*GetRegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _JobRegistration_serviceDesc = grpc.ServiceDesc{
	ServiceName: "jobregistration.JobRegistration",
	HandlerType: (*JobRegistrationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _JobRegistration_Register_Handler,
		},
		{
			MethodName: "GetRegistration",
			Handler:    _JobRegistration_GetRegistration_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "jobregistration.proto",
}