package constant

type RequestMethod string

const (
	GET     RequestMethod = "GET"
	POST    RequestMethod = "POST"
	PUT     RequestMethod = "PUT"
	DELETE  RequestMethod = "DELETE"
	OPTION  RequestMethod = "OPTION"
	PATCH   RequestMethod = "PATCH"
	HEAD    RequestMethod = "HEAD"
	CONNECT RequestMethod = "CONNECT"
	TRACE   RequestMethod = "TRACE"
	ANY     RequestMethod = "ANY"
)
