package main

import (
	"context"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jwenz723/errhandling/grpc/athens/errors"
	"github.com/jwenz723/errhandling/grpc/athens/errorthrower"
	"github.com/jwenz723/errhandling/pb"
	"go.uber.org/zap"
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
	const op = errors.Op("NewOrder")
	err := errorthrower.SomeError()
	if err != nil {
		e := errors.E(op, err, errors.C(req.CustomerID))
		ctxzap.AddFields(ctx,
			zap.String("Error.Ops", errors.OpsText(e)),
			zap.String("Error.Kind", errors.KindText(e)),
			zap.String("Error.CustomerID", string(errors.CustomerID(e))))
		return &pb.NewOrderReply{}, e
	}

	return &pb.NewOrderReply{OrderID: "my order id"}, nil
}