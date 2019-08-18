package main

import (
	"context"
	"github.com/inContact/errhandling/errorthrower"
	"github.com/pkg/errors"
)

// OrderService provides operations for Orders.
type OrderService interface {
	NewOrder(ctx context.Context, customerID string) (string, error)
}

// orderService is a concrete implementation of OrderService
type orderService struct{}

func NewService() orderService {
	return orderService{}
}

func (orderService) NewOrder(ctx context.Context, customerID string) (string, error) {
	if customerID == "" {
		return "", ErrEmpty
	}

	err := errorthrower.SomeError()
	if err != nil {
		return "", errors.Wrap(err, "service.NewOrder")
	}

	return "my order id", nil
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")
