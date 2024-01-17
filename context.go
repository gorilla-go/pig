package pig

import (
	"bytes"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla-go/pig/di"
	"github.com/gorilla-go/pig/foundation"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Context struct {
	container *di.Container
	paramVar  *ReqParamHelper
	paramOnce sync.Once
	postVar   *ReqParamHelper
	postOnce  sync.Once
	fileVar   map[string]*File
	fileOnce  sync.Once
	config    IConfig
}

func NewContext() *Context {
	return &Context{
		container: di.New(),
	}
}

func (c *Context) Container() *di.Container {
	return c.container
}

func (c *Context) routerParams() RouterParams {
	routerParams, err := di.Invoke[RouterParams](c.container)
	if err != nil {
		return nil
	}

	return routerParams
}

func (c *Context) Request() *http.Request {
	return di.MustInvoke[*http.Request](c.container)
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return di.MustInvoke[http.ResponseWriter](c.container)
}

func (c *Context) ParamVar() *ReqParamHelper {
	c.paramOnce.Do(func() {
		paramVar := make(map[string]*ReqParamV)

		request := di.MustInvoke[*http.Request](c.container)
		for n, v := range request.URL.Query() {
			paramVar[n] = NewReqParamV(v)
		}

		routerParams := c.routerParams()
		if routerParams != nil {
			for n, v := range routerParams {
				paramVar[n] = v
			}
		}

		c.paramVar = NewReqParamHelper(paramVar)
	})

	return c.paramVar
}

func (c *Context) PostVar() *ReqParamHelper {
	c.postOnce.Do(func() {
		postVar := make(map[string]*ReqParamV)
		request := di.MustInvoke[*http.Request](c.container)
		err := request.ParseForm()
		if err != nil {
			panic(err)
		}
		for n, v := range request.PostForm {
			postVar[n] = NewReqParamV(v)
		}

		if len(postVar) == 0 {
			multipartReader, err := request.MultipartReader()
			if err != nil {
				panic(err)
			}
			for true {
				part, err := multipartReader.NextPart()
				if err != nil {
					if err == io.EOF {
						break
					}
					panic(err)
				}
				fileName := part.FileName()
				formName := part.FormName()
				if len(formName) > 0 && fileName == "" {
					buf := new(bytes.Buffer)
					_, err := buf.ReadFrom(part)
					if err != nil {
						panic(err)
					}
					postVar[formName] = NewReqParamV([]string{buf.String()})
				}
			}
		}

		c.postVar = NewReqParamHelper(postVar)
	})
	return c.postVar
}

func (c *Context) FileVar() map[string]*File {
	c.fileOnce.Do(func() {
		request := c.Request()
		c.fileVar = make(map[string]*File)

		multipartReader, err := request.MultipartReader()
		if err != nil {
			return
		}
		for true {
			part, err := multipartReader.NextPart()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}

			fileName := part.FileName()
			formName := part.FormName()
			if len(formName) > 0 && len(fileName) > 0 {
				ext := filepath.Ext(fileName)
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(part)
				if err != nil {
					panic(err)
				}

				// save file at tmp dir
				node, err := snowflake.NewNode(int64(rand.Intn(100)))
				if err != nil {
					panic(err)
				}
				fileId := node.Generate().Bytes()
				filename := os.TempDir() + string(fileId) + ext

				file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					panic(err)
				}

				_, err = file.Write(buf.Bytes())
				if err != nil {
					panic(err)
				}

				err = file.Close()
				if err != nil {
					panic(err)
				}

				c.fileVar[formName] = &File{
					FilePath:    filename,
					ContentType: part.Header.Get("Content-Type"),
					Basename:    fileName,
					Ext:         ext,
				}
			}
		}
	})
	return c.fileVar
}

func (c *Context) Download(file *File, basename string) {
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

	c.ResponseWriter().Header().Set("Content-Type", file.ContentType)
	c.ResponseWriter().Header().Set(
		"Content-Disposition",
		"attachment; filename="+basename,
	)
	_, err = io.Copy(c.ResponseWriter(), f)
	if err != nil {
		panic(err)
	}
}

func (c *Context) Redirect(uri string, code ...int) {
	http.Redirect(
		c.ResponseWriter(),
		c.Request(),
		uri,
		foundation.DefaultParam(code, http.StatusFound),
	)
}

func (c *Context) Json(o any, code ...int) {
	httpCode := foundation.DefaultParam(code, http.StatusOK)
	c.ResponseWriter().WriteHeader(httpCode)
	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	marshal, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	_, err = c.ResponseWriter().Write(marshal)
	if err != nil {
		panic(err)
	}
}

func (c *Context) Echo(s string, code ...int) {
	httpCode := foundation.DefaultParam(code, http.StatusOK)
	c.ResponseWriter().WriteHeader(httpCode)
	c.ResponseWriter().Header().Set("Content-Type", "text/plain")
	_, err := c.ResponseWriter().Write([]byte(s))
	if err != nil {
		panic(err)
	}
}

func (c *Context) Logger() ILogger {
	return di.MustInvoke[ILogger](c.container)
}

func (c *Context) Config(s string) any {
	if c.config == nil {
		c.config = di.MustInvoke[IConfig](c.container)
	}
	v, err := c.config.Get(s)
	if err != nil {
		panic(err)
	}
	return v
}
