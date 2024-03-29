// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.6.1
// source: controls.proto

package rpc

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

const (
	GApi_CreateAccount_FullMethodName    = "/ghostnet.rpc.GApi/CreateAccount"
	GApi_CreateGenesis_FullMethodName    = "/ghostnet.rpc.GApi/CreateGenesis"
	GApi_GetPrivateKey_FullMethodName    = "/ghostnet.rpc.GApi/GetPrivateKey"
	GApi_LoginContainer_FullMethodName   = "/ghostnet.rpc.GApi/LoginContainer"
	GApi_ForkContainer_FullMethodName    = "/ghostnet.rpc.GApi/ForkContainer"
	GApi_CreateContainer_FullMethodName  = "/ghostnet.rpc.GApi/CreateContainer"
	GApi_ControlContainer_FullMethodName = "/ghostnet.rpc.GApi/ControlContainer"
	GApi_ReleaseContainer_FullMethodName = "/ghostnet.rpc.GApi/ReleaseContainer"
	GApi_GetContainerList_FullMethodName = "/ghostnet.rpc.GApi/GetContainerList"
	GApi_GetLog_FullMethodName           = "/ghostnet.rpc.GApi/GetLog"
	GApi_GetInfo_FullMethodName          = "/ghostnet.rpc.GApi/GetInfo"
	GApi_CheckStatus_FullMethodName      = "/ghostnet.rpc.GApi/CheckStatus"
	GApi_GetAccount_FullMethodName       = "/ghostnet.rpc.GApi/GetAccount"
	GApi_GetBlockInfo_FullMethodName     = "/ghostnet.rpc.GApi/GetBlockInfo"
)

// GApiClient is the client API for GApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GApiClient interface {
	CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error)
	CreateGenesis(ctx context.Context, in *CreateGenesisRequest, opts ...grpc.CallOption) (*CreateGenesisResponse, error)
	GetPrivateKey(ctx context.Context, in *PrivateKeyRequest, opts ...grpc.CallOption) (*PrivateKeyResponse, error)
	LoginContainer(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	ForkContainer(ctx context.Context, in *ForkContainerRequest, opts ...grpc.CallOption) (*ForkContainerResponse, error)
	CreateContainer(ctx context.Context, in *CreateContainerRequest, opts ...grpc.CallOption) (*CreateContainerResponse, error)
	ControlContainer(ctx context.Context, in *ControlContainerRequest, opts ...grpc.CallOption) (*ControlContainerResponse, error)
	ReleaseContainer(ctx context.Context, in *ReleaseRequest, opts ...grpc.CallOption) (*ReleaseResponse, error)
	GetContainerList(ctx context.Context, in *GetContainerListRequest, opts ...grpc.CallOption) (*GetContainerListResponse, error)
	GetLog(ctx context.Context, in *GetLogRequest, opts ...grpc.CallOption) (*GetLogResponse, error)
	GetInfo(ctx context.Context, in *GetInfoRequest, opts ...grpc.CallOption) (*GetInfoResponse, error)
	CheckStatus(ctx context.Context, in *CoreStateRequest, opts ...grpc.CallOption) (*CoreStateResponse, error)
	GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error)
	GetBlockInfo(ctx context.Context, in *GetBlockInfoRequest, opts ...grpc.CallOption) (*GetBlockInfoResponse, error)
}

type gApiClient struct {
	cc grpc.ClientConnInterface
}

func NewGApiClient(cc grpc.ClientConnInterface) GApiClient {
	return &gApiClient{cc}
}

func (c *gApiClient) CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error) {
	out := new(CreateAccountResponse)
	err := c.cc.Invoke(ctx, GApi_CreateAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) CreateGenesis(ctx context.Context, in *CreateGenesisRequest, opts ...grpc.CallOption) (*CreateGenesisResponse, error) {
	out := new(CreateGenesisResponse)
	err := c.cc.Invoke(ctx, GApi_CreateGenesis_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) GetPrivateKey(ctx context.Context, in *PrivateKeyRequest, opts ...grpc.CallOption) (*PrivateKeyResponse, error) {
	out := new(PrivateKeyResponse)
	err := c.cc.Invoke(ctx, GApi_GetPrivateKey_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) LoginContainer(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, GApi_LoginContainer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) ForkContainer(ctx context.Context, in *ForkContainerRequest, opts ...grpc.CallOption) (*ForkContainerResponse, error) {
	out := new(ForkContainerResponse)
	err := c.cc.Invoke(ctx, GApi_ForkContainer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) CreateContainer(ctx context.Context, in *CreateContainerRequest, opts ...grpc.CallOption) (*CreateContainerResponse, error) {
	out := new(CreateContainerResponse)
	err := c.cc.Invoke(ctx, GApi_CreateContainer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) ControlContainer(ctx context.Context, in *ControlContainerRequest, opts ...grpc.CallOption) (*ControlContainerResponse, error) {
	out := new(ControlContainerResponse)
	err := c.cc.Invoke(ctx, GApi_ControlContainer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) ReleaseContainer(ctx context.Context, in *ReleaseRequest, opts ...grpc.CallOption) (*ReleaseResponse, error) {
	out := new(ReleaseResponse)
	err := c.cc.Invoke(ctx, GApi_ReleaseContainer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) GetContainerList(ctx context.Context, in *GetContainerListRequest, opts ...grpc.CallOption) (*GetContainerListResponse, error) {
	out := new(GetContainerListResponse)
	err := c.cc.Invoke(ctx, GApi_GetContainerList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) GetLog(ctx context.Context, in *GetLogRequest, opts ...grpc.CallOption) (*GetLogResponse, error) {
	out := new(GetLogResponse)
	err := c.cc.Invoke(ctx, GApi_GetLog_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) GetInfo(ctx context.Context, in *GetInfoRequest, opts ...grpc.CallOption) (*GetInfoResponse, error) {
	out := new(GetInfoResponse)
	err := c.cc.Invoke(ctx, GApi_GetInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) CheckStatus(ctx context.Context, in *CoreStateRequest, opts ...grpc.CallOption) (*CoreStateResponse, error) {
	out := new(CoreStateResponse)
	err := c.cc.Invoke(ctx, GApi_CheckStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error) {
	out := new(GetAccountResponse)
	err := c.cc.Invoke(ctx, GApi_GetAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gApiClient) GetBlockInfo(ctx context.Context, in *GetBlockInfoRequest, opts ...grpc.CallOption) (*GetBlockInfoResponse, error) {
	out := new(GetBlockInfoResponse)
	err := c.cc.Invoke(ctx, GApi_GetBlockInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GApiServer is the server API for GApi service.
// All implementations must embed UnimplementedGApiServer
// for forward compatibility
type GApiServer interface {
	CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error)
	CreateGenesis(context.Context, *CreateGenesisRequest) (*CreateGenesisResponse, error)
	GetPrivateKey(context.Context, *PrivateKeyRequest) (*PrivateKeyResponse, error)
	LoginContainer(context.Context, *LoginRequest) (*LoginResponse, error)
	ForkContainer(context.Context, *ForkContainerRequest) (*ForkContainerResponse, error)
	CreateContainer(context.Context, *CreateContainerRequest) (*CreateContainerResponse, error)
	ControlContainer(context.Context, *ControlContainerRequest) (*ControlContainerResponse, error)
	ReleaseContainer(context.Context, *ReleaseRequest) (*ReleaseResponse, error)
	GetContainerList(context.Context, *GetContainerListRequest) (*GetContainerListResponse, error)
	GetLog(context.Context, *GetLogRequest) (*GetLogResponse, error)
	GetInfo(context.Context, *GetInfoRequest) (*GetInfoResponse, error)
	CheckStatus(context.Context, *CoreStateRequest) (*CoreStateResponse, error)
	GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error)
	GetBlockInfo(context.Context, *GetBlockInfoRequest) (*GetBlockInfoResponse, error)
	mustEmbedUnimplementedGApiServer()
}

// UnimplementedGApiServer must be embedded to have forward compatible implementations.
type UnimplementedGApiServer struct {
}

func (UnimplementedGApiServer) CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (UnimplementedGApiServer) CreateGenesis(context.Context, *CreateGenesisRequest) (*CreateGenesisResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGenesis not implemented")
}
func (UnimplementedGApiServer) GetPrivateKey(context.Context, *PrivateKeyRequest) (*PrivateKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPrivateKey not implemented")
}
func (UnimplementedGApiServer) LoginContainer(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginContainer not implemented")
}
func (UnimplementedGApiServer) ForkContainer(context.Context, *ForkContainerRequest) (*ForkContainerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForkContainer not implemented")
}
func (UnimplementedGApiServer) CreateContainer(context.Context, *CreateContainerRequest) (*CreateContainerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateContainer not implemented")
}
func (UnimplementedGApiServer) ControlContainer(context.Context, *ControlContainerRequest) (*ControlContainerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ControlContainer not implemented")
}
func (UnimplementedGApiServer) ReleaseContainer(context.Context, *ReleaseRequest) (*ReleaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseContainer not implemented")
}
func (UnimplementedGApiServer) GetContainerList(context.Context, *GetContainerListRequest) (*GetContainerListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetContainerList not implemented")
}
func (UnimplementedGApiServer) GetLog(context.Context, *GetLogRequest) (*GetLogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLog not implemented")
}
func (UnimplementedGApiServer) GetInfo(context.Context, *GetInfoRequest) (*GetInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (UnimplementedGApiServer) CheckStatus(context.Context, *CoreStateRequest) (*CoreStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckStatus not implemented")
}
func (UnimplementedGApiServer) GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (UnimplementedGApiServer) GetBlockInfo(context.Context, *GetBlockInfoRequest) (*GetBlockInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockInfo not implemented")
}
func (UnimplementedGApiServer) mustEmbedUnimplementedGApiServer() {}

// UnsafeGApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GApiServer will
// result in compilation errors.
type UnsafeGApiServer interface {
	mustEmbedUnimplementedGApiServer()
}

func RegisterGApiServer(s grpc.ServiceRegistrar, srv GApiServer) {
	s.RegisterService(&GApi_ServiceDesc, srv)
}

func _GApi_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_CreateAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).CreateAccount(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_CreateGenesis_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGenesisRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).CreateGenesis(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_CreateGenesis_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).CreateGenesis(ctx, req.(*CreateGenesisRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_GetPrivateKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrivateKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).GetPrivateKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_GetPrivateKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).GetPrivateKey(ctx, req.(*PrivateKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_LoginContainer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).LoginContainer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_LoginContainer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).LoginContainer(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_ForkContainer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForkContainerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).ForkContainer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_ForkContainer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).ForkContainer(ctx, req.(*ForkContainerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_CreateContainer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateContainerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).CreateContainer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_CreateContainer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).CreateContainer(ctx, req.(*CreateContainerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_ControlContainer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ControlContainerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).ControlContainer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_ControlContainer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).ControlContainer(ctx, req.(*ControlContainerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_ReleaseContainer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).ReleaseContainer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_ReleaseContainer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).ReleaseContainer(ctx, req.(*ReleaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_GetContainerList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetContainerListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).GetContainerList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_GetContainerList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).GetContainerList(ctx, req.(*GetContainerListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_GetLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).GetLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_GetLog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).GetLog(ctx, req.(*GetLogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_GetInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).GetInfo(ctx, req.(*GetInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_CheckStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CoreStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).CheckStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_CheckStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).CheckStatus(ctx, req.(*CoreStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_GetAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).GetAccount(ctx, req.(*GetAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GApi_GetBlockInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GApiServer).GetBlockInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GApi_GetBlockInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GApiServer).GetBlockInfo(ctx, req.(*GetBlockInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GApi_ServiceDesc is the grpc.ServiceDesc for GApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ghostnet.rpc.GApi",
	HandlerType: (*GApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAccount",
			Handler:    _GApi_CreateAccount_Handler,
		},
		{
			MethodName: "CreateGenesis",
			Handler:    _GApi_CreateGenesis_Handler,
		},
		{
			MethodName: "GetPrivateKey",
			Handler:    _GApi_GetPrivateKey_Handler,
		},
		{
			MethodName: "LoginContainer",
			Handler:    _GApi_LoginContainer_Handler,
		},
		{
			MethodName: "ForkContainer",
			Handler:    _GApi_ForkContainer_Handler,
		},
		{
			MethodName: "CreateContainer",
			Handler:    _GApi_CreateContainer_Handler,
		},
		{
			MethodName: "ControlContainer",
			Handler:    _GApi_ControlContainer_Handler,
		},
		{
			MethodName: "ReleaseContainer",
			Handler:    _GApi_ReleaseContainer_Handler,
		},
		{
			MethodName: "GetContainerList",
			Handler:    _GApi_GetContainerList_Handler,
		},
		{
			MethodName: "GetLog",
			Handler:    _GApi_GetLog_Handler,
		},
		{
			MethodName: "GetInfo",
			Handler:    _GApi_GetInfo_Handler,
		},
		{
			MethodName: "CheckStatus",
			Handler:    _GApi_CheckStatus_Handler,
		},
		{
			MethodName: "GetAccount",
			Handler:    _GApi_GetAccount_Handler,
		},
		{
			MethodName: "GetBlockInfo",
			Handler:    _GApi_GetBlockInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "controls.proto",
}
