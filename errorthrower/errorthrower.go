package errorthrower

import "fmt"

func SomeError() error {
	return fmt.Errorf("an error has occurred")
}
