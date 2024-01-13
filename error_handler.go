package pig

import (
	"fmt"
	"runtime/debug"
)

type ErrorHandler struct {
}

func NewHttpErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (e ErrorHandler) Handle(err any, context *Context) {
	// set http code
	context.ResponseWriter().WriteHeader(500)
	errno := fmt.Sprintf("%v", err)
	context.Logger().Warning(errno)
	_, err = context.ResponseWriter().Write(
		[]byte(fmt.Sprintf("Error: %s\n\r%s", errno, debug.Stack())),
	)
	if err != nil {
		panic(err)
	}
}
