package errorthrower

import (
	"github.com/jwenz723/errhandling/grpc/athens/errors"
	"google.golang.org/grpc/codes"
)

func SomeError() error {
	const op = errors.Op("SomeError")
	return errors.E(op, LevelOne())
}

func LevelOne() error {
	const op = errors.Op("LevelOne")
	return errors.E(op, LevelTwo())
}

func LevelTwo() error {
	const op = errors.Op("LevelTwo")
	return errors.E(op, LevelThree())
}

func LevelThree() error {
	const op = errors.Op("LevelThree")
	c := codes.Internal
	gm := errors.GM("grpc status message")

	// Can pass in a value of type string or error to be stored as an inner error
	//err := fmt.Errorf("my base error")
	//return errors.E(op, err, c, gm, errors.KindBadRequest)
	return errors.E(op, "my base error", c, gm, errors.KindBadRequest)
}