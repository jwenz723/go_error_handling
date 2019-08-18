package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/xerrors"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	NewOrderEndpoint    endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func NewSet(svc OrderService, logger log.Logger) Set {
	var newOrderEndpoint endpoint.Endpoint
	{
		methodLogger := log.With(logger, "method", "NewOrder")
		infoLogger := level.Info(methodLogger)
		errorLogger := level.Error(methodLogger)
		newOrderEndpoint = MakeNewOrderEndpoint(svc)
		newOrderEndpoint = LoggingMiddleware(infoLogger, errorLogger)(newOrderEndpoint)
	}
	return Set{
		NewOrderEndpoint: newOrderEndpoint,
	}
}

// Sum implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) NewOrder(ctx context.Context, customerID string) (string, error) {
	resp, err := s.NewOrderEndpoint(ctx, NewOrderRequest{CustomerID: customerID})
	if err != nil {
		return "", err
	}
	response := resp.(NewOrderResponse)
	return response.OrderID, response.Err
}

// MakeSumEndpoint constructs a Sum endpoint wrapping the service.
func MakeNewOrderEndpoint(s OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(NewOrderRequest)
		orderID, err := s.NewOrder(ctx, req.CustomerID)
		return NewOrderResponse{OrderID: orderID, Err: xerrors.Errorf("endpoint.NewOrder: %w", err)}, nil
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = NewOrderResponse{}
)

type NewOrderRequest struct {
	CustomerID string
}

// AppendKeyvals implements eplogger.AppendKeyvalser
func (r NewOrderRequest) AppendKeyvals(keyvals []interface{}) []interface{} {
	return append(keyvals,
		"NewOrderRequest.OrderID", r.CustomerID)
}

// SumResponse collects the response values for the Sum method.
type NewOrderResponse struct {
	OrderID   string   `json:"order_id"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// AppendKeyvals implements eplogger.AppendKeyvalser
func (r NewOrderResponse) AppendKeyvals(keyvals []interface{}) []interface{} {
	return append(keyvals,
		"NewOrderResponse.OrderID", r.OrderID,
		"NewOrderResponse.Err", fmt.Sprintf("%+v",r.Err))
}

// Failed implements endpoint.Failer.
func (r NewOrderResponse) Failed() error { return r.Err }