// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: mlops/agent/agent.proto

package agent

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

// AgentServiceClient is the client API for AgentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentServiceClient interface {
	AgentEvent(ctx context.Context, in *ModelEventMessage, opts ...grpc.CallOption) (*ModelEventResponse, error)
	Subscribe(ctx context.Context, in *AgentSubscribeRequest, opts ...grpc.CallOption) (AgentService_SubscribeClient, error)
	ModelScalingTrigger(ctx context.Context, opts ...grpc.CallOption) (AgentService_ModelScalingTriggerClient, error)
	AgentDrain(ctx context.Context, in *AgentDrainRequest, opts ...grpc.CallOption) (*AgentDrainResponse, error)
}

type agentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentServiceClient(cc grpc.ClientConnInterface) AgentServiceClient {
	return &agentServiceClient{cc}
}

func (c *agentServiceClient) AgentEvent(ctx context.Context, in *ModelEventMessage, opts ...grpc.CallOption) (*ModelEventResponse, error) {
	out := new(ModelEventResponse)
	err := c.cc.Invoke(ctx, "/seldon.mlops.agent.AgentService/AgentEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentServiceClient) Subscribe(ctx context.Context, in *AgentSubscribeRequest, opts ...grpc.CallOption) (AgentService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &AgentService_ServiceDesc.Streams[0], "/seldon.mlops.agent.AgentService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &agentServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AgentService_SubscribeClient interface {
	Recv() (*ModelOperationMessage, error)
	grpc.ClientStream
}

type agentServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *agentServiceSubscribeClient) Recv() (*ModelOperationMessage, error) {
	m := new(ModelOperationMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *agentServiceClient) ModelScalingTrigger(ctx context.Context, opts ...grpc.CallOption) (AgentService_ModelScalingTriggerClient, error) {
	stream, err := c.cc.NewStream(ctx, &AgentService_ServiceDesc.Streams[1], "/seldon.mlops.agent.AgentService/ModelScalingTrigger", opts...)
	if err != nil {
		return nil, err
	}
	x := &agentServiceModelScalingTriggerClient{stream}
	return x, nil
}

type AgentService_ModelScalingTriggerClient interface {
	Send(*ModelScalingTriggerMessage) error
	CloseAndRecv() (*ModelScalingTriggerResponse, error)
	grpc.ClientStream
}

type agentServiceModelScalingTriggerClient struct {
	grpc.ClientStream
}

func (x *agentServiceModelScalingTriggerClient) Send(m *ModelScalingTriggerMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *agentServiceModelScalingTriggerClient) CloseAndRecv() (*ModelScalingTriggerResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(ModelScalingTriggerResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *agentServiceClient) AgentDrain(ctx context.Context, in *AgentDrainRequest, opts ...grpc.CallOption) (*AgentDrainResponse, error) {
	out := new(AgentDrainResponse)
	err := c.cc.Invoke(ctx, "/seldon.mlops.agent.AgentService/AgentDrain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentServiceServer is the server API for AgentService service.
// All implementations must embed UnimplementedAgentServiceServer
// for forward compatibility
type AgentServiceServer interface {
	AgentEvent(context.Context, *ModelEventMessage) (*ModelEventResponse, error)
	Subscribe(*AgentSubscribeRequest, AgentService_SubscribeServer) error
	ModelScalingTrigger(AgentService_ModelScalingTriggerServer) error
	AgentDrain(context.Context, *AgentDrainRequest) (*AgentDrainResponse, error)
	mustEmbedUnimplementedAgentServiceServer()
}

// UnimplementedAgentServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAgentServiceServer struct {
}

func (UnimplementedAgentServiceServer) AgentEvent(context.Context, *ModelEventMessage) (*ModelEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AgentEvent not implemented")
}
func (UnimplementedAgentServiceServer) Subscribe(*AgentSubscribeRequest, AgentService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedAgentServiceServer) ModelScalingTrigger(AgentService_ModelScalingTriggerServer) error {
	return status.Errorf(codes.Unimplemented, "method ModelScalingTrigger not implemented")
}
func (UnimplementedAgentServiceServer) AgentDrain(context.Context, *AgentDrainRequest) (*AgentDrainResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AgentDrain not implemented")
}
func (UnimplementedAgentServiceServer) mustEmbedUnimplementedAgentServiceServer() {}

// UnsafeAgentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentServiceServer will
// result in compilation errors.
type UnsafeAgentServiceServer interface {
	mustEmbedUnimplementedAgentServiceServer()
}

func RegisterAgentServiceServer(s grpc.ServiceRegistrar, srv AgentServiceServer) {
	s.RegisterService(&AgentService_ServiceDesc, srv)
}

func _AgentService_AgentEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModelEventMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).AgentEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seldon.mlops.agent.AgentService/AgentEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).AgentEvent(ctx, req.(*ModelEventMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AgentSubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AgentServiceServer).Subscribe(m, &agentServiceSubscribeServer{stream})
}

type AgentService_SubscribeServer interface {
	Send(*ModelOperationMessage) error
	grpc.ServerStream
}

type agentServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *agentServiceSubscribeServer) Send(m *ModelOperationMessage) error {
	return x.ServerStream.SendMsg(m)
}

func _AgentService_ModelScalingTrigger_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AgentServiceServer).ModelScalingTrigger(&agentServiceModelScalingTriggerServer{stream})
}

type AgentService_ModelScalingTriggerServer interface {
	SendAndClose(*ModelScalingTriggerResponse) error
	Recv() (*ModelScalingTriggerMessage, error)
	grpc.ServerStream
}

type agentServiceModelScalingTriggerServer struct {
	grpc.ServerStream
}

func (x *agentServiceModelScalingTriggerServer) SendAndClose(m *ModelScalingTriggerResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *agentServiceModelScalingTriggerServer) Recv() (*ModelScalingTriggerMessage, error) {
	m := new(ModelScalingTriggerMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _AgentService_AgentDrain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentDrainRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).AgentDrain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seldon.mlops.agent.AgentService/AgentDrain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).AgentDrain(ctx, req.(*AgentDrainRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AgentService_ServiceDesc is the grpc.ServiceDesc for AgentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AgentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "seldon.mlops.agent.AgentService",
	HandlerType: (*AgentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AgentEvent",
			Handler:    _AgentService_AgentEvent_Handler,
		},
		{
			MethodName: "AgentDrain",
			Handler:    _AgentService_AgentDrain_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _AgentService_Subscribe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ModelScalingTrigger",
			Handler:       _AgentService_ModelScalingTrigger_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "mlops/agent/agent.proto",
}
