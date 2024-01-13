package pig

import "fmt"

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
	_, err = context.ResponseWriter().Write([]byte(errno))
	if err != nil {
		panic(err)
	}
}
