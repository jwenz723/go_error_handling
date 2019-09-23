package errorthrower

import (
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SomeError() error {
	return errors.Wrap(LevelOne(), "SomeError")
}

func LevelOne() error {
	return errors.Wrap(LevelTwo(), "LevelOne")
}

func LevelTwo() error {
	return errors.Wrap(LevelThree(), "LevelTwo")
}

func LevelThree() error {
	return errors.Wrap(&testError{err: fmt.Errorf("my base error")}, "LevelThree")
}

type testError struct {
	err error
}

func (t testError) Error() string {
	return fmt.Sprintf("an error inside errorthrower: %s", t.err)
}

func (t *testError) Unwrap() error {
	return t.err
}

func (t testError) GRPCStatus() *status.Status {
	return status.New(codes.Internal, "my custom grpc status message")
}