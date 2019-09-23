package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	errors2 "github.com/jwenz723/errhandling/kit/athens/errors"
	"github.com/jwenz723/errhandling/pb"
	"google.golang.org/grpc"
)

type grpcServer struct {
	newOrder grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(endpoints Set, logger log.Logger) pb.OrdersServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		newOrder: grpctransport.NewServer(
			endpoints.NewOrderEndpoint,
			decodeGRPCNewOrderRequest,
			encodeGRPCNewOrderResponse,
			options...,
		),
	}
}

func (s *grpcServer) NewOrder(ctx context.Context, req *pb.NewOrderRequest) (*pb.NewOrderReply, error) {
	_, rep, err := s.newOrder.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.NewOrderReply), nil
}

// decodeGRPCNewOrderRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC sum request to a user-domain sum request. Primarily useful in a server.
func decodeGRPCNewOrderRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.NewOrderRequest)
	return NewOrderRequest{CustomerID: req.CustomerID}, nil
}

// encodeGRPCNewOrderResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain sum response to a gRPC sum reply. Primarily useful in a server.
func encodeGRPCNewOrderResponse(_ context.Context, response interface{}) (interface{}, error) {
	const op = errors2.Op("grpc.NewOrder")
	resp := response.(NewOrderResponse)
	return &pb.NewOrderReply{OrderID: resp.OrderID, Err: errors2.E(op, resp.Err).Error()}, nil
}

// These annoying helper functions are required to translate Go error types to
// and from strings, which is the type we use in our IDLs to represent errors.
// There is special casing to treat empty strings as nil errors.

func str2err(op errors2.Op, s string) error {
	if s == "" {
		return nil
	}
	return errors2.E(op, s)
}

func NewGRPCClient(conn *grpc.ClientConn, logger log.Logger) OrderService {
	pbServiceName := "pb.Orders"

	var newOrderEndpoint endpoint.Endpoint
	{
		methodName := "NewOrder"
		newOrderEndpoint = grpctransport.NewClient(
			conn,
			pbServiceName,
			methodName,
			encodeGRPCNewOrderRequest,
			decodeGRPCNewOrderResponse,
			pb.NewOrderReply{},
		).Endpoint()
	}

	return Set{
		NewOrderEndpoint: newOrderEndpoint,
	}
}

func encodeGRPCNewOrderRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(NewOrderRequest)
	return &pb.NewOrderRequest{CustomerID: req.CustomerID}, nil
}

func decodeGRPCNewOrderResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.NewOrderReply)
	return NewOrderResponse{OrderID: reply.OrderID, Err: str2err(errors2.Op("grpc.NewOrder"), reply.Err)}, nil
}
