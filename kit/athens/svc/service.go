package main

import (
	"context"
	errors2 "github.com/jwenz723/errhandling/kit/athens/errors"
	"github.com/jwenz723/errhandling/pkg/errorthrower"
	"golang.org/x/xerrors"
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
	const op = errors2.Op("service.NewOrder")
	if customerID == "" {
		return "", ErrEmpty
	}

	err := errorthrower.SomeError()
	if err != nil {
		return "", errors2.E(op, err, errors2.C(customerID), errors2.KindBadRequest)
	}

	return "my order id", nil
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = xerrors.New("empty string")
