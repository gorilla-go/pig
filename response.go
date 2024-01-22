package pig

import (
	"encoding/json"
	"github.com/gorilla-go/pig/foundation"
	"io"
	"net/http"
	"os"
)

type Response struct {
	request *http.Request
	writer  http.ResponseWriter
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

func (c *Response) Download(file *File, basename string) {
	// download file
	_, err := os.Stat(file.FilePath)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(file.FilePath)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	c.writer.Header().Set("Content-Type", file.ContentType)
	c.writer.Header().Set(
		"Content-Disposition",
		"attachment; filename="+basename,
	)
	_, err = io.Copy(c.writer, f)
	if err != nil {
		panic(err)
	}
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
	httpCode := foundation.DefaultParam(code, http.StatusOK)
	c.writer.WriteHeader(httpCode)
	c.writer.Header().Set("Content-Type", "application/json")
	marshal, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	_, err = c.writer.Write(marshal)
	if err != nil {
		panic(err)
	}
}

func (c *Response) Text(s string, code ...int) {
	httpCode := foundation.DefaultParam(code, http.StatusOK)
	c.writer.WriteHeader(httpCode)
	c.writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := c.writer.Write([]byte(s))
	if err != nil {
		panic(err)
	}
}

func (c *Response) Html(s string, code ...int) {
	httpCode := foundation.DefaultParam(code, http.StatusOK)
	c.writer.WriteHeader(httpCode)
	c.writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := c.writer.Write([]byte(s))
	if err != nil {
		panic(err)
	}
}
