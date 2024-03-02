package pig

import (
	"encoding/json"
	"github.com/gorilla-go/pig/foundation"
	"github.com/gorilla-go/pig/param"
	"io"
	"net/http"
	"os"
)

type Response struct {
	request      *http.Request
	writer       http.ResponseWriter
	responseCode int
}

func NewResponse(w http.ResponseWriter, request *http.Request) *Response {
	return &Response{
		writer:  w,
		request: request,
	}
}

func (c *Response) Raw() http.ResponseWriter {
	return c.writer
}

func (c *Response) Download(file *param.File, basename string) error {
	f, err := os.Open(file.FilePath)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	writer := c.Raw()
	writer.Header().Set("Content-Type", file.ContentType)
	writer.Header().Set("Content-Disposition", "attachment; filename="+basename)

	_, err = io.Copy(writer, f)
	if err != nil {
		return err
	}
	return nil
}

func (c *Response) Redirect(uri string, code ...int) {
	http.Redirect(
		c.writer,
		c.request,
		uri,
		foundation.DefaultParam(code, http.StatusFound),
	)
}

func (c *Response) Json(o any, code ...int) {
	writer := c.Raw()
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	marshal, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	httpCode := foundation.DefaultParam(code, http.StatusOK)
	if c.responseCode == 0 {
		c.Code(httpCode)
	}
	_, err = writer.Write(marshal)
	if err != nil {
		panic(err)
	}
}

func (c *Response) Text(s string, code ...int) {
	writer := c.Raw()
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")

	httpCode := foundation.DefaultParam(code, http.StatusOK)
	if c.responseCode == 0 {
		c.Code(httpCode)
	}
	_, err := writer.Write([]byte(s))
	if err != nil {
		panic(err)
	}
}

func (c *Response) Html(s string, code ...int) {
	writer := c.Raw()
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	httpCode := foundation.DefaultParam(code, http.StatusOK)
	if c.responseCode == 0 {
		c.Code(httpCode)
	}
	_, err := writer.Write([]byte(s))
	if err != nil {
		panic(err)
	}
}

func (c *Response) Header(key, value string) *Response {
	c.Raw().Header().Set(key, value)
	return c
}

func (c *Response) Code(code int) {
	c.responseCode = code
	c.writer.WriteHeader(code)
}

func (c *Response) GetCode() int {
	return c.responseCode
}
