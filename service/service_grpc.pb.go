// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package service

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

// BotClient is the client API for Bot service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BotClient interface {
	// Connect is the request that is sent by the bot to the server to connect, and it establishes a stream to receive incomingmessages
	Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (Bot_ConnectClient, error)
	// SendMessage is used to send a message to a chat
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
	// RespondToCommand is used to respond to a command from the user
	// For commands sent in private chats, this just sends a regular message
	// In groups, this replies to a specific message
	RespondToCommand(ctx context.Context, in *RespondToCommandRequest, opts ...grpc.CallOption) (*RespondToCommandResponse, error)
	// EditTextMessage is used to edit a text message
	EditTextMessage(ctx context.Context, in *EditTextMessageRequest, opts ...grpc.CallOption) (*EditTextMessageResponse, error)
}

type botClient struct {
	cc grpc.ClientConnInterface
}

func NewBotClient(cc grpc.ClientConnInterface) BotClient {
	return &botClient{cc}
}

func (c *botClient) Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (Bot_ConnectClient, error) {
	stream, err := c.cc.NewStream(ctx, &Bot_ServiceDesc.Streams[0], "/bot.Bot/Connect", opts...)
	if err != nil {
		return nil, err
	}
	x := &botConnectClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Bot_ConnectClient interface {
	Recv() (*MessagesStream, error)
	grpc.ClientStream
}

type botConnectClient struct {
	grpc.ClientStream
}

func (x *botConnectClient) Recv() (*MessagesStream, error) {
	m := new(MessagesStream)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *botClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	out := new(SendMessageResponse)
	err := c.cc.Invoke(ctx, "/bot.Bot/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *botClient) RespondToCommand(ctx context.Context, in *RespondToCommandRequest, opts ...grpc.CallOption) (*RespondToCommandResponse, error) {
	out := new(RespondToCommandResponse)
	err := c.cc.Invoke(ctx, "/bot.Bot/RespondToCommand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *botClient) EditTextMessage(ctx context.Context, in *EditTextMessageRequest, opts ...grpc.CallOption) (*EditTextMessageResponse, error) {
	out := new(EditTextMessageResponse)
	err := c.cc.Invoke(ctx, "/bot.Bot/EditTextMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BotServer is the server API for Bot service.
// All implementations must embed UnimplementedBotServer
// for forward compatibility
type BotServer interface {
	// Connect is the request that is sent by the bot to the server to connect, and it establishes a stream to receive incomingmessages
	Connect(*ConnectRequest, Bot_ConnectServer) error
	// SendMessage is used to send a message to a chat
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error)
	// RespondToCommand is used to respond to a command from the user
	// For commands sent in private chats, this just sends a regular message
	// In groups, this replies to a specific message
	RespondToCommand(context.Context, *RespondToCommandRequest) (*RespondToCommandResponse, error)
	// EditTextMessage is used to edit a text message
	EditTextMessage(context.Context, *EditTextMessageRequest) (*EditTextMessageResponse, error)
	mustEmbedUnimplementedBotServer()
}

// UnimplementedBotServer must be embedded to have forward compatible implementations.
type UnimplementedBotServer struct {
}

func (UnimplementedBotServer) Connect(*ConnectRequest, Bot_ConnectServer) error {
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedBotServer) SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedBotServer) RespondToCommand(context.Context, *RespondToCommandRequest) (*RespondToCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RespondToCommand not implemented")
}
func (UnimplementedBotServer) EditTextMessage(context.Context, *EditTextMessageRequest) (*EditTextMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditTextMessage not implemented")
}
func (UnimplementedBotServer) mustEmbedUnimplementedBotServer() {}

// UnsafeBotServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BotServer will
// result in compilation errors.
type UnsafeBotServer interface {
	mustEmbedUnimplementedBotServer()
}

func RegisterBotServer(s grpc.ServiceRegistrar, srv BotServer) {
	s.RegisterService(&Bot_ServiceDesc, srv)
}

func _Bot_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConnectRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BotServer).Connect(m, &botConnectServer{stream})
}

type Bot_ConnectServer interface {
	Send(*MessagesStream) error
	grpc.ServerStream
}

type botConnectServer struct {
	grpc.ServerStream
}

func (x *botConnectServer) Send(m *MessagesStream) error {
	return x.ServerStream.SendMsg(m)
}

func _Bot_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BotServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bot.Bot/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BotServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bot_RespondToCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RespondToCommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BotServer).RespondToCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bot.Bot/RespondToCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BotServer).RespondToCommand(ctx, req.(*RespondToCommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bot_EditTextMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditTextMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BotServer).EditTextMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bot.Bot/EditTextMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BotServer).EditTextMessage(ctx, req.(*EditTextMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Bot_ServiceDesc is the grpc.ServiceDesc for Bot service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Bot_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bot.Bot",
	HandlerType: (*BotServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _Bot_SendMessage_Handler,
		},
		{
			MethodName: "RespondToCommand",
			Handler:    _Bot_RespondToCommand_Handler,
		},
		{
			MethodName: "EditTextMessage",
			Handler:    _Bot_EditTextMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Connect",
			Handler:       _Bot_Connect_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "service.proto",
}