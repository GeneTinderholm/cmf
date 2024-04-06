package stackerr

import (
	"fmt"
	"runtime/debug"
)

func New(message string) error {
	stack := debug.Stack()
	return StackErr{
		message: message,
		stack: string(stack),
	}
}

/*
 * StackErr is an error that also carries along a stack trace
 */
type StackErr struct {
	message string
	stack   string
}

func (se StackErr) Error() string {
	return fmt.Sprintf(`%s

%s`, se.message, se.stack)
}
