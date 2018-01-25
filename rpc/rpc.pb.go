// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc.proto

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	rpc.proto

It has these top-level messages:
	Query
	OrderFragment
	ResultFragment
	Address
	MultiAddress
	MultiAddresses
	Nothing
*/
package rpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// A Query message contains the Address of a Node that needs to be found and
// the MultiAddress of the Node from which the Query originated.
type Query struct {
	// Network data.
	From *MultiAddress `protobuf:"bytes,1,opt,name=from" json:"from,omitempty"`
	// Public data.
	Query *Address `protobuf:"bytes,2,opt,name=query" json:"query,omitempty"`
	Deep  bool     `protobuf:"varint,3,opt,name=deep" json:"deep,omitempty"`
}

func (m *Query) Reset()                    { *m = Query{} }
func (m *Query) String() string            { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()               {}
func (*Query) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Query) GetFrom() *MultiAddress {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Query) GetQuery() *Address {
	if m != nil {
		return m.Query
	}
	return nil
}

func (m *Query) GetDeep() bool {
	if m != nil {
		return m.Deep
	}
	return false
}

// An OrderFragment is a message contains the details of an order fragment.
type OrderFragment struct {
	// Network data.
	To   *Address      `protobuf:"bytes,1,opt,name=to" json:"to,omitempty"`
	From *MultiAddress `protobuf:"bytes,2,opt,name=from" json:"from,omitempty"`
	// Public data.
	Id           []byte `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	OrderId      []byte `protobuf:"bytes,4,opt,name=orderId,proto3" json:"orderId,omitempty"`
	OrderType    int64  `protobuf:"varint,5,opt,name=orderType" json:"orderType,omitempty"`
	OrderBuySell int64  `protobuf:"varint,6,opt,name=orderBuySell" json:"orderBuySell,omitempty"`
	// Private data.
	FstCodeShare   []byte `protobuf:"bytes,7,opt,name=fstCodeShare,proto3" json:"fstCodeShare,omitempty"`
	SndCodeShare   []byte `protobuf:"bytes,8,opt,name=sndCodeShare,proto3" json:"sndCodeShare,omitempty"`
	PriceShare     []byte `protobuf:"bytes,9,opt,name=priceShare,proto3" json:"priceShare,omitempty"`
	MaxVolumeShare []byte `protobuf:"bytes,10,opt,name=maxVolumeShare,proto3" json:"maxVolumeShare,omitempty"`
	MinVolumeShare []byte `protobuf:"bytes,11,opt,name=minVolumeShare,proto3" json:"minVolumeShare,omitempty"`
}

func (m *OrderFragment) Reset()                    { *m = OrderFragment{} }
func (m *OrderFragment) String() string            { return proto.CompactTextString(m) }
func (*OrderFragment) ProtoMessage()               {}
func (*OrderFragment) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *OrderFragment) GetTo() *Address {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *OrderFragment) GetFrom() *MultiAddress {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *OrderFragment) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *OrderFragment) GetOrderId() []byte {
	if m != nil {
		return m.OrderId
	}
	return nil
}

func (m *OrderFragment) GetOrderType() int64 {
	if m != nil {
		return m.OrderType
	}
	return 0
}

func (m *OrderFragment) GetOrderBuySell() int64 {
	if m != nil {
		return m.OrderBuySell
	}
	return 0
}

func (m *OrderFragment) GetFstCodeShare() []byte {
	if m != nil {
		return m.FstCodeShare
	}
	return nil
}

func (m *OrderFragment) GetSndCodeShare() []byte {
	if m != nil {
		return m.SndCodeShare
	}
	return nil
}

func (m *OrderFragment) GetPriceShare() []byte {
	if m != nil {
		return m.PriceShare
	}
	return nil
}

func (m *OrderFragment) GetMaxVolumeShare() []byte {
	if m != nil {
		return m.MaxVolumeShare
	}
	return nil
}

func (m *OrderFragment) GetMinVolumeShare() []byte {
	if m != nil {
		return m.MinVolumeShare
	}
	return nil
}

// A ResultFragment message is the network representation of a
// compute.ResultFragment and the metadata needed to distribute it through the
// network.
type ResultFragment struct {
	// Network data.
	To   *Address      `protobuf:"bytes,1,opt,name=to" json:"to,omitempty"`
	From *MultiAddress `protobuf:"bytes,2,opt,name=from" json:"from,omitempty"`
	// Public data.
	Id                  []byte `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	BuyOrderId          []byte `protobuf:"bytes,4,opt,name=buyOrderId,proto3" json:"buyOrderId,omitempty"`
	SellOrderId         []byte `protobuf:"bytes,5,opt,name=sellOrderId,proto3" json:"sellOrderId,omitempty"`
	BuyOrderFragmentId  []byte `protobuf:"bytes,6,opt,name=buyOrderFragmentId,proto3" json:"buyOrderFragmentId,omitempty"`
	SellOrderFragmentId []byte `protobuf:"bytes,7,opt,name=sellOrderFragmentId,proto3" json:"sellOrderFragmentId,omitempty"`
	// Private data.
	FstCodeShare   []byte `protobuf:"bytes,8,opt,name=fstCodeShare,proto3" json:"fstCodeShare,omitempty"`
	SndCodeShare   []byte `protobuf:"bytes,9,opt,name=sndCodeShare,proto3" json:"sndCodeShare,omitempty"`
	PriceShare     []byte `protobuf:"bytes,10,opt,name=priceShare,proto3" json:"priceShare,omitempty"`
	MaxVolumeShare []byte `protobuf:"bytes,11,opt,name=maxVolumeShare,proto3" json:"maxVolumeShare,omitempty"`
	MinVolumeShare []byte `protobuf:"bytes,12,opt,name=minVolumeShare,proto3" json:"minVolumeShare,omitempty"`
}

func (m *ResultFragment) Reset()                    { *m = ResultFragment{} }
func (m *ResultFragment) String() string            { return proto.CompactTextString(m) }
func (*ResultFragment) ProtoMessage()               {}
func (*ResultFragment) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ResultFragment) GetTo() *Address {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *ResultFragment) GetFrom() *MultiAddress {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *ResultFragment) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *ResultFragment) GetBuyOrderId() []byte {
	if m != nil {
		return m.BuyOrderId
	}
	return nil
}

func (m *ResultFragment) GetSellOrderId() []byte {
	if m != nil {
		return m.SellOrderId
	}
	return nil
}

func (m *ResultFragment) GetBuyOrderFragmentId() []byte {
	if m != nil {
		return m.BuyOrderFragmentId
	}
	return nil
}

func (m *ResultFragment) GetSellOrderFragmentId() []byte {
	if m != nil {
		return m.SellOrderFragmentId
	}
	return nil
}

func (m *ResultFragment) GetFstCodeShare() []byte {
	if m != nil {
		return m.FstCodeShare
	}
	return nil
}

func (m *ResultFragment) GetSndCodeShare() []byte {
	if m != nil {
		return m.SndCodeShare
	}
	return nil
}

func (m *ResultFragment) GetPriceShare() []byte {
	if m != nil {
		return m.PriceShare
	}
	return nil
}

func (m *ResultFragment) GetMaxVolumeShare() []byte {
	if m != nil {
		return m.MaxVolumeShare
	}
	return nil
}

func (m *ResultFragment) GetMinVolumeShare() []byte {
	if m != nil {
		return m.MinVolumeShare
	}
	return nil
}

// An Address message is the network representation of an identity.Address.
type Address struct {
	Address string `protobuf:"bytes,1,opt,name=address" json:"address,omitempty"`
}

func (m *Address) Reset()                    { *m = Address{} }
func (m *Address) String() string            { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()               {}
func (*Address) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Address) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

// A MultiAddress is the public multiaddress of a Node in the overlay network.
// It provides the Republic address of the Node, as well as the network
// address.
type MultiAddress struct {
	Multi string `protobuf:"bytes,1,opt,name=multi" json:"multi,omitempty"`
}

func (m *MultiAddress) Reset()                    { *m = MultiAddress{} }
func (m *MultiAddress) String() string            { return proto.CompactTextString(m) }
func (*MultiAddress) ProtoMessage()               {}
func (*MultiAddress) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *MultiAddress) GetMulti() string {
	if m != nil {
		return m.Multi
	}
	return ""
}

// MultiAddresses are public multiaddress of multiple Nodes in the overlay
// network.
type MultiAddresses struct {
	Multis []*MultiAddress `protobuf:"bytes,1,rep,name=multis" json:"multis,omitempty"`
}

func (m *MultiAddresses) Reset()                    { *m = MultiAddresses{} }
func (m *MultiAddresses) String() string            { return proto.CompactTextString(m) }
func (*MultiAddresses) ProtoMessage()               {}
func (*MultiAddresses) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *MultiAddresses) GetMultis() []*MultiAddress {
	if m != nil {
		return m.Multis
	}
	return nil
}

// Nothing is in this message. It is used to send nothing, or signal a
// successful response.
type Nothing struct {
}

func (m *Nothing) Reset()                    { *m = Nothing{} }
func (m *Nothing) String() string            { return proto.CompactTextString(m) }
func (*Nothing) ProtoMessage()               {}
func (*Nothing) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func init() {
	proto.RegisterType((*Query)(nil), "rpc.Query")
	proto.RegisterType((*OrderFragment)(nil), "rpc.OrderFragment")
	proto.RegisterType((*ResultFragment)(nil), "rpc.ResultFragment")
	proto.RegisterType((*Address)(nil), "rpc.Address")
	proto.RegisterType((*MultiAddress)(nil), "rpc.MultiAddress")
	proto.RegisterType((*MultiAddresses)(nil), "rpc.MultiAddresses")
	proto.RegisterType((*Nothing)(nil), "rpc.Nothing")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Node service

type NodeClient interface {
	// Ping the connection and swap MultiAddresses.
	Ping(ctx context.Context, in *MultiAddress, opts ...grpc.CallOption) (*Nothing, error)
	// Get all peers connected to the Node.
	Peers(ctx context.Context, in *MultiAddress, opts ...grpc.CallOption) (*MultiAddresses, error)
	// Find the MultiAddresses of peers closer to some target Node.
	QueryCloserPeers(ctx context.Context, in *Query, opts ...grpc.CallOption) (*MultiAddresses, error)
	// Send an OrderFragment to some target Node.
	SendOrderFragment(ctx context.Context, in *OrderFragment, opts ...grpc.CallOption) (*Nothing, error)
	// Send a ResultFragment to some target Node, where the ResultFragment is the
	// result of a computation on two OrderFragments.
	SendResultFragment(ctx context.Context, in *ResultFragment, opts ...grpc.CallOption) (*Nothing, error)
}

type nodeClient struct {
	cc *grpc.ClientConn
}

func NewNodeClient(cc *grpc.ClientConn) NodeClient {
	return &nodeClient{cc}
}

func (c *nodeClient) Ping(ctx context.Context, in *MultiAddress, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := grpc.Invoke(ctx, "/rpc.Node/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Peers(ctx context.Context, in *MultiAddress, opts ...grpc.CallOption) (*MultiAddresses, error) {
	out := new(MultiAddresses)
	err := grpc.Invoke(ctx, "/rpc.Node/Peers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) QueryCloserPeers(ctx context.Context, in *Query, opts ...grpc.CallOption) (*MultiAddresses, error) {
	out := new(MultiAddresses)
	err := grpc.Invoke(ctx, "/rpc.Node/QueryCloserPeers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) SendOrderFragment(ctx context.Context, in *OrderFragment, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := grpc.Invoke(ctx, "/rpc.Node/SendOrderFragment", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) SendResultFragment(ctx context.Context, in *ResultFragment, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := grpc.Invoke(ctx, "/rpc.Node/SendResultFragment", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Node service

type NodeServer interface {
	// Ping the connection and swap MultiAddresses.
	Ping(context.Context, *MultiAddress) (*Nothing, error)
	// Get all peers connected to the Node.
	Peers(context.Context, *MultiAddress) (*MultiAddresses, error)
	// Find the MultiAddresses of peers closer to some target Node.
	QueryCloserPeers(context.Context, *Query) (*MultiAddresses, error)
	// Send an OrderFragment to some target Node.
	SendOrderFragment(context.Context, *OrderFragment) (*Nothing, error)
	// Send a ResultFragment to some target Node, where the ResultFragment is the
	// result of a computation on two OrderFragments.
	SendResultFragment(context.Context, *ResultFragment) (*Nothing, error)
}

func RegisterNodeServer(s *grpc.Server, srv NodeServer) {
	s.RegisterService(&_Node_serviceDesc, srv)
}

func _Node_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MultiAddress)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Node/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Ping(ctx, req.(*MultiAddress))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Peers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MultiAddress)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Peers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Node/Peers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Peers(ctx, req.(*MultiAddress))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_QueryCloserPeers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).QueryCloserPeers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Node/QueryCloserPeers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).QueryCloserPeers(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_SendOrderFragment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderFragment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).SendOrderFragment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Node/SendOrderFragment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).SendOrderFragment(ctx, req.(*OrderFragment))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_SendResultFragment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResultFragment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).SendResultFragment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Node/SendResultFragment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).SendResultFragment(ctx, req.(*ResultFragment))
	}
	return interceptor(ctx, in, info, handler)
}

var _Node_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Node",
	HandlerType: (*NodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Node_Ping_Handler,
		},
		{
			MethodName: "Peers",
			Handler:    _Node_Peers_Handler,
		},
		{
			MethodName: "QueryCloserPeers",
			Handler:    _Node_QueryCloserPeers_Handler,
		},
		{
			MethodName: "SendOrderFragment",
			Handler:    _Node_SendOrderFragment_Handler,
		},
		{
			MethodName: "SendResultFragment",
			Handler:    _Node_SendResultFragment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc.proto",
}

func init() { proto.RegisterFile("rpc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 504 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x94, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0x15, 0xff, 0x49, 0xe2, 0x89, 0x89, 0xe8, 0x94, 0x83, 0x85, 0xaa, 0x2a, 0x32, 0xff,
	0xc2, 0x25, 0xa0, 0x56, 0x88, 0x03, 0x27, 0xa8, 0x84, 0xd4, 0x03, 0x6d, 0x71, 0x10, 0xf7, 0x34,
	0x3b, 0x49, 0x2d, 0xd9, 0x5e, 0xb3, 0x6b, 0x4b, 0xe4, 0x01, 0x78, 0x30, 0xde, 0x8a, 0x23, 0xf2,
	0xd8, 0x6e, 0xec, 0xc4, 0x82, 0x5c, 0x7a, 0xdb, 0xf9, 0xe6, 0xf7, 0x79, 0x56, 0xfb, 0xad, 0x17,
	0x1c, 0x95, 0x2e, 0x67, 0xa9, 0x92, 0x99, 0x44, 0x53, 0xa5, 0x4b, 0x7f, 0x05, 0xf6, 0xd7, 0x9c,
	0xd4, 0x06, 0x5f, 0x80, 0xb5, 0x52, 0x32, 0xf6, 0x7a, 0x93, 0xde, 0x74, 0x74, 0x76, 0x34, 0x2b,
	0xb8, 0x2f, 0x79, 0x94, 0x85, 0x1f, 0x85, 0x50, 0xa4, 0x75, 0xc0, 0x6d, 0xf4, 0xc1, 0xfe, 0x51,
	0xf0, 0x9e, 0xc1, 0x9c, 0xcb, 0x5c, 0x8d, 0x94, 0x2d, 0x44, 0xb0, 0x04, 0x51, 0xea, 0x99, 0x93,
	0xde, 0x74, 0x18, 0xf0, 0xda, 0xff, 0x63, 0xc0, 0xa3, 0x6b, 0x25, 0x48, 0x7d, 0x56, 0x8b, 0x75,
	0x4c, 0x49, 0x86, 0x27, 0x60, 0x64, 0xb2, 0x1a, 0xd7, 0xfe, 0x8c, 0x91, 0xc9, 0xfb, 0xed, 0x18,
	0xff, 0xde, 0xce, 0x18, 0x8c, 0x50, 0xf0, 0x20, 0x37, 0x30, 0x42, 0x81, 0x1e, 0x0c, 0x64, 0x31,
	0xe5, 0x52, 0x78, 0x16, 0x8b, 0x75, 0x89, 0x27, 0xe0, 0xf0, 0xf2, 0xdb, 0x26, 0x25, 0xcf, 0x9e,
	0xf4, 0xa6, 0x66, 0xb0, 0x15, 0xd0, 0x07, 0x97, 0x8b, 0x4f, 0xf9, 0x66, 0x4e, 0x51, 0xe4, 0xf5,
	0x19, 0x68, 0x69, 0x05, 0xb3, 0xd2, 0xd9, 0x85, 0x14, 0x34, 0xbf, 0x5b, 0x28, 0xf2, 0x06, 0x3c,
	0xa0, 0xa5, 0x15, 0x8c, 0x4e, 0xc4, 0x96, 0x19, 0x96, 0x4c, 0x53, 0xc3, 0x53, 0x80, 0x54, 0x85,
	0xcb, 0x8a, 0x70, 0x98, 0x68, 0x28, 0xf8, 0x12, 0xc6, 0xf1, 0xe2, 0xe7, 0x77, 0x19, 0xe5, 0x71,
	0xc5, 0x00, 0x33, 0x3b, 0x2a, 0x73, 0x61, 0xd2, 0xe4, 0x46, 0x15, 0xd7, 0x52, 0xfd, 0xdf, 0x26,
	0x8c, 0x03, 0xd2, 0x79, 0x94, 0x3d, 0xec, 0xd9, 0x9f, 0x02, 0xdc, 0xe6, 0x9b, 0xeb, 0xd6, 0xf1,
	0x37, 0x14, 0x9c, 0xc0, 0x48, 0x53, 0x14, 0xd5, 0x80, 0xcd, 0x40, 0x53, 0xc2, 0x19, 0x60, 0xcd,
	0xd7, 0x5b, 0xbd, 0x14, 0x9c, 0x85, 0x1b, 0x74, 0x74, 0xf0, 0x2d, 0x1c, 0xdf, 0xdb, 0x1b, 0x86,
	0x32, 0x98, 0xae, 0xd6, 0x5e, 0x86, 0xc3, 0x03, 0x32, 0x74, 0xfe, 0x9b, 0x21, 0x1c, 0x90, 0xe1,
	0xe8, 0xc0, 0x0c, 0xdd, 0xce, 0x0c, 0x9f, 0xc1, 0xa0, 0x3a, 0xfc, 0xe2, 0x8a, 0x2f, 0xca, 0x25,
	0x07, 0xe8, 0x04, 0x75, 0xe9, 0x3f, 0x07, 0xb7, 0x19, 0x13, 0x3e, 0x01, 0x3b, 0x2e, 0xea, 0x8a,
	0x2b, 0x0b, 0xff, 0x03, 0x8c, 0x9b, 0x14, 0x69, 0x7c, 0x0d, 0x7d, 0x6e, 0x15, 0x1f, 0x34, 0xbb,
	0x13, 0xaf, 0x00, 0xdf, 0x81, 0xc1, 0x95, 0xcc, 0xee, 0xc2, 0x64, 0x7d, 0xf6, 0xcb, 0x00, 0xeb,
	0x4a, 0x0a, 0xc2, 0x57, 0x60, 0xdd, 0x84, 0xc9, 0x1a, 0xf7, 0x6d, 0x4f, 0xcb, 0xbb, 0x55, 0x39,
	0xf0, 0x0d, 0xd8, 0x37, 0x44, 0x4a, 0x77, 0x91, 0xc7, 0x7b, 0x12, 0x69, 0x3c, 0x87, 0xc7, 0xfc,
	0x38, 0x5d, 0x44, 0x52, 0x93, 0x2a, 0xbd, 0xc0, 0x20, 0xcb, 0xdd, 0xa6, 0x77, 0x70, 0x34, 0xa7,
	0x44, 0xb4, 0x1f, 0x1b, 0x64, 0xb2, 0xa5, 0xed, 0x6c, 0xee, 0x3d, 0x60, 0x61, 0xdb, 0xf9, 0x51,
	0xca, 0x09, 0x6d, 0xb1, 0x6d, 0xbc, 0xed, 0xf3, 0x6b, 0x7a, 0xfe, 0x37, 0x00, 0x00, 0xff, 0xff,
	0xd0, 0x18, 0x09, 0x0d, 0x5a, 0x05, 0x00, 0x00,
}
