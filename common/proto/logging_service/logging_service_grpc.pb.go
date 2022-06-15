// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: logging_service.proto

package logging_service

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

// LoggingServiceClient is the client API for LoggingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoggingServiceClient interface {
	LoggInfo(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResult, error)
	LoggError(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResult, error)
	LoggWarning(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResult, error)
	LoggSuccess(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResult, error)
}

type loggingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLoggingServiceClient(cc grpc.ClientConnInterface) LoggingServiceClient {
	return &loggingServiceClient{cc}
}

func (c *loggingServiceClient) LoggInfo(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResult, error) {
	out := new(LogResult)
	err := c.cc.Invoke(ctx, "/logging_service.LoggingService/LoggInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingServiceClient) LoggError(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResult, error) {
	out := new(LogResult)
	err := c.cc.Invoke(ctx, "/logging_service.LoggingService/LoggError", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingServiceClient) LoggWarning(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResult, error) {
	out := new(LogResult)
	err := c.cc.Invoke(ctx, "/logging_service.LoggingService/LoggWarning", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingServiceClient) LoggSuccess(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResult, error) {
	out := new(LogResult)
	err := c.cc.Invoke(ctx, "/logging_service.LoggingService/LoggSuccess", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoggingServiceServer is the server API for LoggingService service.
// All implementations must embed UnimplementedLoggingServiceServer
// for forward compatibility
type LoggingServiceServer interface {
	LoggInfo(context.Context, *LogRequest) (*LogResult, error)
	LoggError(context.Context, *LogRequest) (*LogResult, error)
	LoggWarning(context.Context, *LogRequest) (*LogResult, error)
	LoggSuccess(context.Context, *LogRequest) (*LogResult, error)
	mustEmbedUnimplementedLoggingServiceServer()
}

// UnimplementedLoggingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLoggingServiceServer struct {
}

func (UnimplementedLoggingServiceServer) LoggInfo(context.Context, *LogRequest) (*LogResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoggInfo not implemented")
}
func (UnimplementedLoggingServiceServer) LoggError(context.Context, *LogRequest) (*LogResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoggError not implemented")
}
func (UnimplementedLoggingServiceServer) LoggWarning(context.Context, *LogRequest) (*LogResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoggWarning not implemented")
}
func (UnimplementedLoggingServiceServer) LoggSuccess(context.Context, *LogRequest) (*LogResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoggSuccess not implemented")
}
func (UnimplementedLoggingServiceServer) mustEmbedUnimplementedLoggingServiceServer() {}

// UnsafeLoggingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LoggingServiceServer will
// result in compilation errors.
type UnsafeLoggingServiceServer interface {
	mustEmbedUnimplementedLoggingServiceServer()
}

func RegisterLoggingServiceServer(s grpc.ServiceRegistrar, srv LoggingServiceServer) {
	s.RegisterService(&LoggingService_ServiceDesc, srv)
}

func _LoggingService_LoggInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingServiceServer).LoggInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logging_service.LoggingService/LoggInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingServiceServer).LoggInfo(ctx, req.(*LogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingService_LoggError_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingServiceServer).LoggError(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logging_service.LoggingService/LoggError",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingServiceServer).LoggError(ctx, req.(*LogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingService_LoggWarning_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingServiceServer).LoggWarning(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logging_service.LoggingService/LoggWarning",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingServiceServer).LoggWarning(ctx, req.(*LogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingService_LoggSuccess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingServiceServer).LoggSuccess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logging_service.LoggingService/LoggSuccess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingServiceServer).LoggSuccess(ctx, req.(*LogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LoggingService_ServiceDesc is the grpc.ServiceDesc for LoggingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LoggingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "logging_service.LoggingService",
	HandlerType: (*LoggingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LoggInfo",
			Handler:    _LoggingService_LoggInfo_Handler,
		},
		{
			MethodName: "LoggError",
			Handler:    _LoggingService_LoggError_Handler,
		},
		{
			MethodName: "LoggWarning",
			Handler:    _LoggingService_LoggWarning_Handler,
		},
		{
			MethodName: "LoggSuccess",
			Handler:    _LoggingService_LoggSuccess_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "logging_service.proto",
}