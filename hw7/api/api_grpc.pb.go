package api

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

// MailServiceClient is the client API for MailService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MailServiceClient interface {
	MailSend(ctx context.Context, in *MailSendRequest, opts ...grpc.CallOption) (*Empty, error)
}

type mailServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMailServiceClient(cc grpc.ClientConnInterface) MailServiceClient {
	return &mailServiceClient{cc}
}

func (c *mailServiceClient) MailSend(ctx context.Context, in *MailSendRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/api.MailService/MailSend", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MailServiceServer is the server API for MailService service.
// All implementations must embed UnimplementedMailServiceServer
// for forward compatibility
type MailServiceServer interface {
	MailSend(context.Context, *MailSendRequest) (*Empty, error)
	mustEmbedUnimplementedMailServiceServer()
}

// UnimplementedMailServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMailServiceServer struct {
}

func (UnimplementedMailServiceServer) MailSend(context.Context, *MailSendRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MailSend not implemented")
}
func (UnimplementedMailServiceServer) mustEmbedUnimplementedMailServiceServer() {}

// UnsafeMailServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MailServiceServer will
// result in compilation errors.
type UnsafeMailServiceServer interface {
	mustEmbedUnimplementedMailServiceServer()
}

func RegisterMailServiceServer(s grpc.ServiceRegistrar, srv MailServiceServer) {
	s.RegisterService(&MailService_ServiceDesc, srv)
}

func _MailService_MailSend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MailSendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).MailSend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.MailService/MailSend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).MailSend(ctx, req.(*MailSendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MailService_ServiceDesc is the grpc.ServiceDesc for MailService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MailService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.MailService",
	HandlerType: (*MailServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MailSend",
			Handler:    _MailService_MailSend_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/api.proto",
}

// PersonalAccountServiceClient is the client API for PersonalAccountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PersonalAccountServiceClient interface {
	PersonalAccount(ctx context.Context, in *PersonalAccountRequest, opts ...grpc.CallOption) (*PersonalAccountResponse, error)
}

type personalAccountServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPersonalAccountServiceClient(cc grpc.ClientConnInterface) PersonalAccountServiceClient {
	return &personalAccountServiceClient{cc}
}

func (c *personalAccountServiceClient) PersonalAccount(ctx context.Context, in *PersonalAccountRequest, opts ...grpc.CallOption) (*PersonalAccountResponse, error) {
	out := new(PersonalAccountResponse)
	err := c.cc.Invoke(ctx, "/api.PersonalAccountService/PersonalAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PersonalAccountServiceServer is the server API for PersonalAccountService service.
// All implementations must embed UnimplementedPersonalAccountServiceServer
// for forward compatibility
type PersonalAccountServiceServer interface {
	PersonalAccount(context.Context, *PersonalAccountRequest) (*PersonalAccountResponse, error)
	mustEmbedUnimplementedPersonalAccountServiceServer()
}

// UnimplementedPersonalAccountServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPersonalAccountServiceServer struct {
}

func (UnimplementedPersonalAccountServiceServer) PersonalAccount(context.Context, *PersonalAccountRequest) (*PersonalAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PersonalAccount not implemented")
}
func (UnimplementedPersonalAccountServiceServer) mustEmbedUnimplementedPersonalAccountServiceServer() {
}

// UnsafePersonalAccountServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PersonalAccountServiceServer will
// result in compilation errors.
type UnsafePersonalAccountServiceServer interface {
	mustEmbedUnimplementedPersonalAccountServiceServer()
}

func RegisterPersonalAccountServiceServer(s grpc.ServiceRegistrar, srv PersonalAccountServiceServer) {
	s.RegisterService(&PersonalAccountService_ServiceDesc, srv)
}

func _PersonalAccountService_PersonalAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PersonalAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonalAccountServiceServer).PersonalAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.PersonalAccountService/PersonalAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonalAccountServiceServer).PersonalAccount(ctx, req.(*PersonalAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PersonalAccountService_ServiceDesc is the grpc.ServiceDesc for PersonalAccountService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PersonalAccountService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.PersonalAccountService",
	HandlerType: (*PersonalAccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PersonalAccount",
			Handler:    _PersonalAccountService_PersonalAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/api.proto",
}
