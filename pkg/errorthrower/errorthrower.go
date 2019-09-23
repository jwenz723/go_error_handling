package errorthrower

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SomeError() error {
	return fmt.Errorf("SomeError: %w", LevelOne())
}

func LevelOne() error {
	return fmt.Errorf("LevelOne: %w", LevelTwo())
}

func LevelTwo() error {
	return fmt.Errorf("LevelTwo: %w", LevelThree())
}

func LevelThree() error {
	//return fmt.Errorf("an error has occurred")
	return fmt.Errorf("LevelThree: %w", &testError{})
}

type testError struct {
	err error
}

func (t testError) Error() string {
	return "a testError inside errorthrower"
}

func (t *testError) Unwrap() error {
	return t.err
}

func (t testError) GRPCStatus() *status.Status {
	return status.New(codes.Internal, "errorthrower status message")
}