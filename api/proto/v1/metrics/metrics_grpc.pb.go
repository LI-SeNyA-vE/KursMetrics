// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: metrics.proto

package metrics

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MetricsService_SendMetric_FullMethodName       = "/metrics.v1.MetricsService/SendMetric"
	MetricsService_SendBatchMetrics_FullMethodName = "/metrics.v1.MetricsService/SendBatchMetrics"
	MetricsService_GetMetric_FullMethodName        = "/metrics.v1.MetricsService/GetMetric"
	MetricsService_GetAllMetrics_FullMethodName    = "/metrics.v1.MetricsService/GetAllMetrics"
)

// MetricsServiceClient is the client API for MetricsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Определение сервиса для работы с метриками
type MetricsServiceClient interface {
	// Метод для отправки одной метрики
	SendMetric(ctx context.Context, in *SendMetricRequest, opts ...grpc.CallOption) (*SendMetricResponse, error)
	// Метод для отправки нескольких метрик (batch)
	SendBatchMetrics(ctx context.Context, in *BatchMetricsRequest, opts ...grpc.CallOption) (*BatchMetricsResponse, error)
	// Метод для получения метрики по имени
	GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error)
	// Метод для получения всех метрик
	GetAllMetrics(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AllMetricsResponse, error)
}

type metricsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsServiceClient(cc grpc.ClientConnInterface) MetricsServiceClient {
	return &metricsServiceClient{cc}
}

func (c *metricsServiceClient) SendMetric(ctx context.Context, in *SendMetricRequest, opts ...grpc.CallOption) (*SendMetricResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendMetricResponse)
	err := c.cc.Invoke(ctx, MetricsService_SendMetric_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsServiceClient) SendBatchMetrics(ctx context.Context, in *BatchMetricsRequest, opts ...grpc.CallOption) (*BatchMetricsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchMetricsResponse)
	err := c.cc.Invoke(ctx, MetricsService_SendBatchMetrics_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsServiceClient) GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetMetricResponse)
	err := c.cc.Invoke(ctx, MetricsService_GetMetric_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsServiceClient) GetAllMetrics(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AllMetricsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AllMetricsResponse)
	err := c.cc.Invoke(ctx, MetricsService_GetAllMetrics_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServiceServer is the server API for MetricsService service.
// All implementations must embed UnimplementedMetricsServiceServer
// for forward compatibility.
//
// Определение сервиса для работы с метриками
type MetricsServiceServer interface {
	// Метод для отправки одной метрики
	SendMetric(context.Context, *SendMetricRequest) (*SendMetricResponse, error)
	// Метод для отправки нескольких метрик (batch)
	SendBatchMetrics(context.Context, *BatchMetricsRequest) (*BatchMetricsResponse, error)
	// Метод для получения метрики по имени
	GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error)
	// Метод для получения всех метрик
	GetAllMetrics(context.Context, *Empty) (*AllMetricsResponse, error)
	mustEmbedUnimplementedMetricsServiceServer()
}

// UnimplementedMetricsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMetricsServiceServer struct{}

func (UnimplementedMetricsServiceServer) SendMetric(context.Context, *SendMetricRequest) (*SendMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMetric not implemented")
}
func (UnimplementedMetricsServiceServer) SendBatchMetrics(context.Context, *BatchMetricsRequest) (*BatchMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendBatchMetrics not implemented")
}
func (UnimplementedMetricsServiceServer) GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetric not implemented")
}
func (UnimplementedMetricsServiceServer) GetAllMetrics(context.Context, *Empty) (*AllMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllMetrics not implemented")
}
func (UnimplementedMetricsServiceServer) mustEmbedUnimplementedMetricsServiceServer() {}
func (UnimplementedMetricsServiceServer) testEmbeddedByValue()                        {}

// UnsafeMetricsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServiceServer will
// result in compilation errors.
type UnsafeMetricsServiceServer interface {
	mustEmbedUnimplementedMetricsServiceServer()
}

func RegisterMetricsServiceServer(s grpc.ServiceRegistrar, srv MetricsServiceServer) {
	// If the following call pancis, it indicates UnimplementedMetricsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MetricsService_ServiceDesc, srv)
}

func _MetricsService_SendMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).SendMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_SendMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).SendMetric(ctx, req.(*SendMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_SendBatchMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).SendBatchMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_SendBatchMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).SendBatchMetrics(ctx, req.(*BatchMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_GetMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).GetMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_GetMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).GetMetric(ctx, req.(*GetMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_GetAllMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).GetAllMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_GetAllMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).GetAllMetrics(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// MetricsService_ServiceDesc is the grpc.ServiceDesc for MetricsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetricsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "metrics.v1.MetricsService",
	HandlerType: (*MetricsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMetric",
			Handler:    _MetricsService_SendMetric_Handler,
		},
		{
			MethodName: "SendBatchMetrics",
			Handler:    _MetricsService_SendBatchMetrics_Handler,
		},
		{
			MethodName: "GetMetric",
			Handler:    _MetricsService_GetMetric_Handler,
		},
		{
			MethodName: "GetAllMetrics",
			Handler:    _MetricsService_GetAllMetrics_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "metrics.proto",
}
