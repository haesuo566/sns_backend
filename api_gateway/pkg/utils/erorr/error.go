package e

import (
	"fmt"
	"runtime"
)

func Wrap(msg string) error {
	_, file, line, ok := runtime.Caller(1)

	if !ok {
		return fmt.Errorf("unknown error: %s", msg)
	}

	return fmt.Errorf("%s : %d : %s", file, line, msg)
}

func WrapWithStack(msg string, stack int) error {
	_, file, line, ok := runtime.Caller(stack)
	if !ok {
		return fmt.Errorf("unknown error: %s", msg)
	}

	return fmt.Errorf("%s : %d : %s", file, line, msg)
}
