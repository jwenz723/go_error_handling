package main

import (
	"context"
	"fmt"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/jwenz723/errhandling/grpc/1.13-xerrors/errorthrower"
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
		e := fmt.Errorf("NewOrder: %w", err)
		return &pb.NewOrderReply{}, e
	}

	return &pb.NewOrderReply{OrderID: "my order id"}, nil
}