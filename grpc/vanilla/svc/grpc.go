package main

import (
	"context"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/jwenz723/errhandling/grpc/vanilla/errorthrower"
	"github.com/jwenz723/errhandling/pb"
)

var _ pb.OrdersServer = &grpcServer{}

type grpcServer struct {}

// GrpcLoggingDecider specifies which methods should have their request/response parameters logged
// by the grpc logging interceptor. Returning false indicates logging should be suppressed.
func (s *grpcServer) GrpcLoggingDecider() grpc_logging.ServerPayloadLoggingDecider {
	return func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		switch fullMethodName {
		default:
			return true
		}
	}
}

func (s *grpcServer) NewOrder(ctx context.Context, req *pb.NewOrderRequest) (*pb.NewOrderReply, error) {
	err := errorthrower.SomeError()
	if err != nil {
		return &pb.NewOrderReply{}, err
	}

	return &pb.NewOrderReply{OrderID: "my order id"}, nil
}