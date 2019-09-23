package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	errors2 "github.com/jwenz723/errhandling/kit/athens/errors"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	NewOrderEndpoint endpoint.Endpoint
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
	const op = errors2.Op("endpoint.NewOrder")
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(NewOrderRequest)
		orderID, err := NewOrder(ctx, req.CustomerID)
		return NewOrderResponse{OrderID: orderID, Err: errors2.E(op, err)}, nil
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
		"NewOrderRequest.CustomerID", r.CustomerID)
}

// SumResponse collects the response values for the Sum method.
type NewOrderResponse struct {
	OrderID string `json:"order_id"`
	Err     error  `json:"-"` // should be intercepted by Failed/errorEncoder
}

// AppendKeyvals implements eplogger.AppendKeyvalser
func (r NewOrderResponse) AppendKeyvals(keyvals []interface{}) []interface{} {
	e, ok := r.Err.(errors2.Error)
	if !ok {
		return append(keyvals,
			"NewOrderResponse.OrderID", r.OrderID,
			"NewOrderResponse.Err", fmt.Sprintf("%+v", r.Err))
	}

	return append(keyvals,
		"NewOrderResponse.OrderID", r.OrderID,
		"NewOrderResponse.Err.Op", e.Op,
		"NewOrderResponse.Err.Kind", errors2.KindText(e),
		"NewOrderResponse.Err.CustomerID", errors2.CustomerID(e),
		"NewOrderResponse.Err.OrderID", errors2.OrderID(e),
		"NewOrderResponse.Err.Ops", errors2.OpsText(e))
}

// Failed implements endpoint.Failer.
func (r NewOrderResponse) Failed() error { return r.Err }
