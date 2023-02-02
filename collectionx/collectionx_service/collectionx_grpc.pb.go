// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.15.8
// source: collectionx.proto

package collectionxservice

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

// ServiceCollectionClient is the client API for ServiceCollection service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceCollectionClient interface {
	Retrive(ctx context.Context, in *RetriveRequest, opts ...grpc.CallOption) (*RetriveResponse, error)
	Snapshots(ctx context.Context, in *SnapshotsRequest, opts ...grpc.CallOption) (ServiceCollection_SnapshotsClient, error)
}

type serviceCollectionClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceCollectionClient(cc grpc.ClientConnInterface) ServiceCollectionClient {
	return &serviceCollectionClient{cc}
}

func (c *serviceCollectionClient) Retrive(ctx context.Context, in *RetriveRequest, opts ...grpc.CallOption) (*RetriveResponse, error) {
	out := new(RetriveResponse)
	err := c.cc.Invoke(ctx, "/collectionx.ServiceCollection/Retrive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceCollectionClient) Snapshots(ctx context.Context, in *SnapshotsRequest, opts ...grpc.CallOption) (ServiceCollection_SnapshotsClient, error) {
	stream, err := c.cc.NewStream(ctx, &ServiceCollection_ServiceDesc.Streams[0], "/collectionx.ServiceCollection/Snapshots", opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceCollectionSnapshotsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ServiceCollection_SnapshotsClient interface {
	Recv() (*SnapshotsResponse, error)
	grpc.ClientStream
}

type serviceCollectionSnapshotsClient struct {
	grpc.ClientStream
}

func (x *serviceCollectionSnapshotsClient) Recv() (*SnapshotsResponse, error) {
	m := new(SnapshotsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ServiceCollectionServer is the server API for ServiceCollection service.
// All implementations must embed UnimplementedServiceCollectionServer
// for forward compatibility
type ServiceCollectionServer interface {
	Retrive(context.Context, *RetriveRequest) (*RetriveResponse, error)
	Snapshots(*SnapshotsRequest, ServiceCollection_SnapshotsServer) error
	mustEmbedUnimplementedServiceCollectionServer()
}

// UnimplementedServiceCollectionServer must be embedded to have forward compatible implementations.
type UnimplementedServiceCollectionServer struct {
}

func (UnimplementedServiceCollectionServer) Retrive(context.Context, *RetriveRequest) (*RetriveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Retrive not implemented")
}
func (UnimplementedServiceCollectionServer) Snapshots(*SnapshotsRequest, ServiceCollection_SnapshotsServer) error {
	return status.Errorf(codes.Unimplemented, "method Snapshots not implemented")
}
func (UnimplementedServiceCollectionServer) mustEmbedUnimplementedServiceCollectionServer() {}

// UnsafeServiceCollectionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceCollectionServer will
// result in compilation errors.
type UnsafeServiceCollectionServer interface {
	mustEmbedUnimplementedServiceCollectionServer()
}

func RegisterServiceCollectionServer(s grpc.ServiceRegistrar, srv ServiceCollectionServer) {
	s.RegisterService(&ServiceCollection_ServiceDesc, srv)
}

func _ServiceCollection_Retrive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetriveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceCollectionServer).Retrive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/collectionx.ServiceCollection/Retrive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceCollectionServer).Retrive(ctx, req.(*RetriveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceCollection_Snapshots_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SnapshotsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ServiceCollectionServer).Snapshots(m, &serviceCollectionSnapshotsServer{stream})
}

type ServiceCollection_SnapshotsServer interface {
	Send(*SnapshotsResponse) error
	grpc.ServerStream
}

type serviceCollectionSnapshotsServer struct {
	grpc.ServerStream
}

func (x *serviceCollectionSnapshotsServer) Send(m *SnapshotsResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ServiceCollection_ServiceDesc is the grpc.ServiceDesc for ServiceCollection service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServiceCollection_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "collectionx.ServiceCollection",
	HandlerType: (*ServiceCollectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Retrive",
			Handler:    _ServiceCollection_Retrive_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Snapshots",
			Handler:       _ServiceCollection_Snapshots_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "collectionx.proto",
}