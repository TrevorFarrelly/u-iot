// Code generated by protoc-gen-go. DO NOT EDIT.
// source: uiot.proto

package uiot

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
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

// Meta contains human-readable information about a device. Name, type of device,
// and room it is located in
type Meta struct {
	Type                 uint32   `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	Room                 uint32   `protobuf:"varint,2,opt,name=room,proto3" json:"room,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Meta) Reset()         { *m = Meta{} }
func (m *Meta) String() string { return proto.CompactTextString(m) }
func (*Meta) ProtoMessage()    {}
func (*Meta) Descriptor() ([]byte, []int) {
	return fileDescriptor_761f8d2c87cd1351, []int{0}
}

func (m *Meta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Meta.Unmarshal(m, b)
}
func (m *Meta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Meta.Marshal(b, m, deterministic)
}
func (m *Meta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta.Merge(m, src)
}
func (m *Meta) XXX_Size() int {
	return xxx_messageInfo_Meta.Size(m)
}
func (m *Meta) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta.DiscardUnknown(m)
}

var xxx_messageInfo_Meta proto.InternalMessageInfo

func (m *Meta) GetType() uint32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *Meta) GetRoom() uint32 {
	if m != nil {
		return m.Room
	}
	return 0
}

func (m *Meta) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Param messages represent valid parameters for a given function.
// Contains mininum and maximum values the parameter can take.
type Param struct {
	Min                  uint32   `protobuf:"varint,1,opt,name=min,proto3" json:"min,omitempty"`
	Max                  uint32   `protobuf:"varint,2,opt,name=max,proto3" json:"max,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Param) Reset()         { *m = Param{} }
func (m *Param) String() string { return proto.CompactTextString(m) }
func (*Param) ProtoMessage()    {}
func (*Param) Descriptor() ([]byte, []int) {
	return fileDescriptor_761f8d2c87cd1351, []int{1}
}

func (m *Param) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Param.Unmarshal(m, b)
}
func (m *Param) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Param.Marshal(b, m, deterministic)
}
func (m *Param) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Param.Merge(m, src)
}
func (m *Param) XXX_Size() int {
	return xxx_messageInfo_Param.Size(m)
}
func (m *Param) XXX_DiscardUnknown() {
	xxx_messageInfo_Param.DiscardUnknown(m)
}

var xxx_messageInfo_Param proto.InternalMessageInfo

func (m *Param) GetMin() uint32 {
	if m != nil {
		return m.Min
	}
	return 0
}

func (m *Param) GetMax() uint32 {
	if m != nil {
		return m.Max
	}
	return 0
}

// Func messages represent individual functions a device can perform.
// Contains an ID, human-readable name, and list of parameters.
type Func struct {
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Params               []*Param `protobuf:"bytes,3,rep,name=params,proto3" json:"params,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Func) Reset()         { *m = Func{} }
func (m *Func) String() string { return proto.CompactTextString(m) }
func (*Func) ProtoMessage()    {}
func (*Func) Descriptor() ([]byte, []int) {
	return fileDescriptor_761f8d2c87cd1351, []int{2}
}

func (m *Func) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Func.Unmarshal(m, b)
}
func (m *Func) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Func.Marshal(b, m, deterministic)
}
func (m *Func) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Func.Merge(m, src)
}
func (m *Func) XXX_Size() int {
	return xxx_messageInfo_Func.Size(m)
}
func (m *Func) XXX_DiscardUnknown() {
	xxx_messageInfo_Func.DiscardUnknown(m)
}

var xxx_messageInfo_Func proto.InternalMessageInfo

func (m *Func) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Func) GetParams() []*Param {
	if m != nil {
		return m.Params
	}
	return nil
}

// DevInfo message are sent when a new device is connecting to the network.
// Contains a device's identifying information and all of the functions it performs.
type DevInfo struct {
	Port                 uint32   `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Meta                 *Meta    `protobuf:"bytes,3,opt,name=meta,proto3" json:"meta,omitempty"`
	Funcs                []*Func  `protobuf:"bytes,4,rep,name=funcs,proto3" json:"funcs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DevInfo) Reset()         { *m = DevInfo{} }
func (m *DevInfo) String() string { return proto.CompactTextString(m) }
func (*DevInfo) ProtoMessage()    {}
func (*DevInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_761f8d2c87cd1351, []int{3}
}

func (m *DevInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DevInfo.Unmarshal(m, b)
}
func (m *DevInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DevInfo.Marshal(b, m, deterministic)
}
func (m *DevInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DevInfo.Merge(m, src)
}
func (m *DevInfo) XXX_Size() int {
	return xxx_messageInfo_DevInfo.Size(m)
}
func (m *DevInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_DevInfo.DiscardUnknown(m)
}

var xxx_messageInfo_DevInfo proto.InternalMessageInfo

func (m *DevInfo) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *DevInfo) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *DevInfo) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *DevInfo) GetFuncs() []*Func {
	if m != nil {
		return m.Funcs
	}
	return nil
}

func init() {
	proto.RegisterType((*Meta)(nil), "Meta")
	proto.RegisterType((*Param)(nil), "Param")
	proto.RegisterType((*Func)(nil), "Func")
	proto.RegisterType((*DevInfo)(nil), "DevInfo")
}

func init() { proto.RegisterFile("uiot.proto", fileDescriptor_761f8d2c87cd1351) }

var fileDescriptor_761f8d2c87cd1351 = []byte{
	// 241 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xc1, 0x4b, 0xc3, 0x30,
	0x14, 0x87, 0xed, 0x9a, 0x56, 0xf7, 0x86, 0x20, 0x39, 0x45, 0x05, 0xa9, 0x39, 0x15, 0x06, 0x39,
	0xd4, 0x9b, 0xc7, 0x31, 0x04, 0x0f, 0x82, 0xe4, 0x3f, 0x88, 0x5d, 0x06, 0x3d, 0x24, 0x2f, 0xa4,
	0xe9, 0xd0, 0xff, 0x5e, 0x5e, 0xd6, 0x75, 0xb7, 0x8f, 0x0f, 0xf2, 0xe5, 0x97, 0x00, 0x4c, 0x03,
	0x26, 0x15, 0x22, 0x26, 0x94, 0x3b, 0x60, 0x5f, 0x36, 0x19, 0xce, 0x81, 0xa5, 0xbf, 0x60, 0x45,
	0xd1, 0x14, 0xed, 0xbd, 0xce, 0x4c, 0x2e, 0x22, 0x3a, 0xb1, 0x3a, 0x3b, 0x62, 0x72, 0xde, 0x38,
	0x2b, 0xca, 0xa6, 0x68, 0xd7, 0x3a, 0xb3, 0xdc, 0x42, 0xf5, 0x6d, 0xa2, 0x71, 0xfc, 0x01, 0x4a,
	0x37, 0xf8, 0xb9, 0x41, 0x98, 0x8d, 0xf9, 0x9d, 0x0b, 0x84, 0xf2, 0x1d, 0xd8, 0xc7, 0xe4, 0xfb,
	0x25, 0xb4, 0xba, 0x86, 0xf8, 0x0b, 0xd4, 0x81, 0x42, 0xa3, 0x28, 0x9b, 0xb2, 0xdd, 0x74, 0xb5,
	0xca, 0x5d, 0x3d, 0x5b, 0x39, 0xc0, 0xed, 0xde, 0x9e, 0x3e, 0xfd, 0x11, 0xe9, 0x78, 0xc0, 0x98,
	0x2e, 0x7b, 0x89, 0xc9, 0x99, 0xc3, 0x21, 0x5e, 0x92, 0xc4, 0xfc, 0x11, 0x98, 0xb3, 0xc9, 0xe4,
	0xbd, 0x9b, 0xae, 0x52, 0xf4, 0x58, 0x9d, 0x15, 0x7f, 0x86, 0xea, 0x38, 0xf9, 0x7e, 0x14, 0x2c,
	0x5f, 0x56, 0x29, 0xda, 0xa5, 0xcf, 0xae, 0xdb, 0x42, 0xbd, 0xb7, 0xa7, 0xa1, 0xb7, 0xfc, 0x15,
	0xd6, 0x3b, 0xc4, 0x34, 0xa6, 0x68, 0x02, 0xbf, 0x53, 0xf3, 0x80, 0xa7, 0x85, 0xe4, 0xcd, 0x4f,
	0x9d, 0xff, 0xf2, 0xed, 0x3f, 0x00, 0x00, 0xff, 0xff, 0xcd, 0xc1, 0x52, 0x63, 0x59, 0x01, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DeviceClient is the client API for Device service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DeviceClient interface {
	// Send our information to a remote device, and get theirs in return
	Bootstrap(ctx context.Context, in *DevInfo, opts ...grpc.CallOption) (*DevInfo, error)
}

type deviceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceClient(cc grpc.ClientConnInterface) DeviceClient {
	return &deviceClient{cc}
}

func (c *deviceClient) Bootstrap(ctx context.Context, in *DevInfo, opts ...grpc.CallOption) (*DevInfo, error) {
	out := new(DevInfo)
	err := c.cc.Invoke(ctx, "/Device/Bootstrap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceServer is the server API for Device service.
type DeviceServer interface {
	// Send our information to a remote device, and get theirs in return
	Bootstrap(context.Context, *DevInfo) (*DevInfo, error)
}

// UnimplementedDeviceServer can be embedded to have forward compatible implementations.
type UnimplementedDeviceServer struct {
}

func (*UnimplementedDeviceServer) Bootstrap(ctx context.Context, req *DevInfo) (*DevInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bootstrap not implemented")
}

func RegisterDeviceServer(s *grpc.Server, srv DeviceServer) {
	s.RegisterService(&_Device_serviceDesc, srv)
}

func _Device_Bootstrap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DevInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServer).Bootstrap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Device/Bootstrap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServer).Bootstrap(ctx, req.(*DevInfo))
	}
	return interceptor(ctx, in, info, handler)
}

var _Device_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Device",
	HandlerType: (*DeviceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Bootstrap",
			Handler:    _Device_Bootstrap_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "uiot.proto",
}