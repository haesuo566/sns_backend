package e

import (
	"fmt"
	"runtime"
)

func Wrap(err error) error {
	_, file, line, ok := runtime.Caller(1)

	if !ok {
		return fmt.Errorf("unknown error: %s", err.Error())
	}

	return fmt.Errorf("%s : %d : %s", file, line, err.Error())
}

func WrapWithStack(err error, stack int) error {
	_, file, line, ok := runtime.Caller(stack)
	if !ok {
		return fmt.Errorf("unknown error: %s", err.Error())
	}

	return fmt.Errorf("%s : %d : %s", file, line, err.Error())
}
