package pig

import (
	"bytes"
	"github.com/bwmarrin/snowflake"
	"github.com/samber/do"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Context struct {
	injector  *do.Injector
	paramVar  map[string]*ReqParamV
	paramOnce sync.Once
	postVar   map[string]*ReqParamV
	postOnce  sync.Once
	fileVar   map[string]*File
	fileOnce  sync.Once
}

func NewContext() *Context {
	return &Context{
		injector: do.New(),
	}
}

func (c *Context) Injector() *do.Injector {
	return c.injector
}

func (c *Context) routerParams() RouterParams {
	routerParams, err := do.Invoke[RouterParams](c.injector)
	if err != nil {
		return nil
	}

	return routerParams
}

func (c *Context) Request() *http.Request {
	return do.MustInvoke[*http.Request](c.injector)
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return do.MustInvoke[http.ResponseWriter](c.injector)
}

func (c *Context) ParamVar() map[string]*ReqParamV {
	c.paramOnce.Do(func() {
		c.paramVar = make(map[string]*ReqParamV)

		request, err := do.Invoke[*http.Request](c.Injector())
		if err == nil {
			for n, v := range request.URL.Query() {
				c.paramVar[n] = NewReqParamV(v)
			}
		}

		routerParams := c.routerParams()
		if routerParams != nil {
			for n, v := range routerParams {
				c.paramVar[n] = v
			}
		}
	})

	return c.paramVar
}

func (c *Context) PostVar() map[string]*ReqParamV {
	c.postOnce.Do(func() {
		c.postVar = make(map[string]*ReqParamV)
		request := do.MustInvoke[*http.Request](c.Injector())
		err := request.ParseForm()
		if err != nil {
			panic(err)
		}
		for n, v := range request.PostForm {
			c.postVar[n] = NewReqParamV(v)
		}

		if len(c.postVar) == 0 {
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
					c.postVar[formName] = NewReqParamV([]string{buf.String()})
				}
			}
		}
	})
	return c.postVar
}

func (c *Context) FileVar() map[string]*File {
	c.fileOnce.Do(func() {
		request := c.Request()
		c.fileVar = make(map[string]*File)

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
