// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: node.proto

package proto

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

// NodeServiceClient is the client API for NodeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeServiceClient interface {
	WriteForward(ctx context.Context, in *ForwardRequest, opts ...grpc.CallOption) (*Empty, error)
	Invalidate(ctx context.Context, in *InvalidateRequest, opts ...grpc.CallOption) (*Empty, error)
	Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*Empty, error)
	InitWrite(ctx context.Context, in *InitWriteRequest, opts ...grpc.CallOption) (*Empty, error)
	// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	InitRead(ctx context.Context, in *InitReadRequest, opts ...grpc.CallOption) (*InitReadResponse, error)
	ReadForward(ctx context.Context, in *ReadForwardRequest, opts ...grpc.CallOption) (*Empty, error)
	SendContent(ctx context.Context, in *SendContentRequest, opts ...grpc.CallOption) (*Empty, error)
}

type nodeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeServiceClient(cc grpc.ClientConnInterface) NodeServiceClient {
	return &nodeServiceClient{cc}
}

func (c *nodeServiceClient) WriteForward(ctx context.Context, in *ForwardRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.NodeService/WriteForward", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServiceClient) Invalidate(ctx context.Context, in *InvalidateRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.NodeService/Invalidate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServiceClient) Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.NodeService/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServiceClient) InitWrite(ctx context.Context, in *InitWriteRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.NodeService/InitWrite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServiceClient) InitRead(ctx context.Context, in *InitReadRequest, opts ...grpc.CallOption) (*InitReadResponse, error) {
	out := new(InitReadResponse)
	err := c.cc.Invoke(ctx, "/proto.NodeService/InitRead", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServiceClient) ReadForward(ctx context.Context, in *ReadForwardRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.NodeService/ReadForward", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServiceClient) SendContent(ctx context.Context, in *SendContentRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.NodeService/SendContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServiceServer is the server API for NodeService service.
// All implementations must embed UnimplementedNodeServiceServer
// for forward compatibility
type NodeServiceServer interface {
	WriteForward(context.Context, *ForwardRequest) (*Empty, error)
	Invalidate(context.Context, *InvalidateRequest) (*Empty, error)
	Send(context.Context, *SendRequest) (*Empty, error)
	InitWrite(context.Context, *InitWriteRequest) (*Empty, error)
	// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	InitRead(context.Context, *InitReadRequest) (*InitReadResponse, error)
	ReadForward(context.Context, *ReadForwardRequest) (*Empty, error)
	SendContent(context.Context, *SendContentRequest) (*Empty, error)
	mustEmbedUnimplementedNodeServiceServer()
}

// UnimplementedNodeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNodeServiceServer struct {
}

func (UnimplementedNodeServiceServer) WriteForward(context.Context, *ForwardRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WriteForward not implemented")
}
func (UnimplementedNodeServiceServer) Invalidate(context.Context, *InvalidateRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Invalidate not implemented")
}
func (UnimplementedNodeServiceServer) Send(context.Context, *SendRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedNodeServiceServer) InitWrite(context.Context, *InitWriteRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitWrite not implemented")
}
func (UnimplementedNodeServiceServer) InitRead(context.Context, *InitReadRequest) (*InitReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitRead not implemented")
}
func (UnimplementedNodeServiceServer) ReadForward(context.Context, *ReadForwardRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadForward not implemented")
}
func (UnimplementedNodeServiceServer) SendContent(context.Context, *SendContentRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendContent not implemented")
}
func (UnimplementedNodeServiceServer) mustEmbedUnimplementedNodeServiceServer() {}

// UnsafeNodeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeServiceServer will
// result in compilation errors.
type UnsafeNodeServiceServer interface {
	mustEmbedUnimplementedNodeServiceServer()
}

func RegisterNodeServiceServer(s grpc.ServiceRegistrar, srv NodeServiceServer) {
	s.RegisterService(&NodeService_ServiceDesc, srv)
}

func _NodeService_WriteForward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForwardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).WriteForward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NodeService/WriteForward",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).WriteForward(ctx, req.(*ForwardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeService_Invalidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InvalidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).Invalidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NodeService/Invalidate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).Invalidate(ctx, req.(*InvalidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeService_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NodeService/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).Send(ctx, req.(*SendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeService_InitWrite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitWriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).InitWrite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NodeService/InitWrite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).InitWrite(ctx, req.(*InitWriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeService_InitRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).InitRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NodeService/InitRead",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).InitRead(ctx, req.(*InitReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeService_ReadForward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadForwardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).ReadForward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NodeService/ReadForward",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).ReadForward(ctx, req.(*ReadForwardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeService_SendContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendContentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).SendContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NodeService/SendContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).SendContent(ctx, req.(*SendContentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NodeService_ServiceDesc is the grpc.ServiceDesc for NodeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NodeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.NodeService",
	HandlerType: (*NodeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WriteForward",
			Handler:    _NodeService_WriteForward_Handler,
		},
		{
			MethodName: "Invalidate",
			Handler:    _NodeService_Invalidate_Handler,
		},
		{
			MethodName: "Send",
			Handler:    _NodeService_Send_Handler,
		},
		{
			MethodName: "InitWrite",
			Handler:    _NodeService_InitWrite_Handler,
		},
		{
			MethodName: "InitRead",
			Handler:    _NodeService_InitRead_Handler,
		},
		{
			MethodName: "ReadForward",
			Handler:    _NodeService_ReadForward_Handler,
		},
		{
			MethodName: "SendContent",
			Handler:    _NodeService_SendContent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node.proto",
}
