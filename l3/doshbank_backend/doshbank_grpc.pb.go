// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: doshbank_backend/doshbank.proto

package doshbank_backend

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DoshBankClient is the client API for DoshBank service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DoshBankClient interface {
	GetCurrentReward(ctx context.Context, in *GetCurrentRewardRequest, opts ...grpc.CallOption) (*GetCurrentRewardResponse, error)
}

type doshBankClient struct {
	cc grpc.ClientConnInterface
}

func NewDoshBankClient(cc grpc.ClientConnInterface) DoshBankClient {
	return &doshBankClient{cc}
}

func (c *doshBankClient) GetCurrentReward(ctx context.Context, in *GetCurrentRewardRequest, opts ...grpc.CallOption) (*GetCurrentRewardResponse, error) {
	out := new(GetCurrentRewardResponse)
	err := c.cc.Invoke(ctx, "/doshbank_backend.DoshBank/GetCurrentReward", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DoshBankServer is the server API for DoshBank service.
// All implementations must embed UnimplementedDoshBankServer
// for forward compatibility
type DoshBankServer interface {
	GetCurrentReward(context.Context, *GetCurrentRewardRequest) (*GetCurrentRewardResponse, error)
	mustEmbedUnimplementedDoshBankServer()
}

// UnimplementedDoshBankServer must be embedded to have forward compatible implementations.
type UnimplementedDoshBankServer struct {
}

func (UnimplementedDoshBankServer) GetCurrentReward(context.Context, *GetCurrentRewardRequest) (*GetCurrentRewardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCurrentReward not implemented")
}
func (UnimplementedDoshBankServer) mustEmbedUnimplementedDoshBankServer() {}

// UnsafeDoshBankServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DoshBankServer will
// result in compilation errors.
type UnsafeDoshBankServer interface {
	mustEmbedUnimplementedDoshBankServer()
}

func RegisterDoshBankServer(s grpc.ServiceRegistrar, srv DoshBankServer) {
	s.RegisterService(&DoshBank_ServiceDesc, srv)
}

func _DoshBank_GetCurrentReward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCurrentRewardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DoshBankServer).GetCurrentReward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/doshbank_backend.DoshBank/GetCurrentReward",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DoshBankServer).GetCurrentReward(ctx, req.(*GetCurrentRewardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DoshBank_ServiceDesc is the grpc.ServiceDesc for DoshBank service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DoshBank_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "doshbank_backend.DoshBank",
	HandlerType: (*DoshBankServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCurrentReward",
			Handler:    _DoshBank_GetCurrentReward_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "doshbank_backend/doshbank.proto",
}
