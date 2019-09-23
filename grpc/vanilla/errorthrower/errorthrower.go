package errorthrower

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SomeError() error {
	return LevelOne()
}

func LevelOne() error {
	return LevelTwo()
}

func LevelTwo() error {
	return LevelThree()
}

func LevelThree() error {
	return &testError{err: fmt.Errorf("my base error")}
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