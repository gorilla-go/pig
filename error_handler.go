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
	_, err = context.ResponseWriter().Write([]byte(fmt.Sprintf("%v", err)))
	if err != nil {
		panic(err)
	}
}
