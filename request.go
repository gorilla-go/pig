package pig

import (
	"bytes"
	"github.com/bwmarrin/snowflake"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Request struct {
	request      *http.Request
	routerParams RouterParams
	paramVar     *ReqParamHelper
	paramOnce    sync.Once
	postVar      *ReqParamHelper
	postOnce     sync.Once
	fileVar      map[string]*File
	fileOnce     sync.Once
}

func NewRequest(req *http.Request, routerParams RouterParams) *Request {
	return &Request{
		request: req,
	}
}

func (c *Request) ParamVar() *ReqParamHelper {
	c.paramOnce.Do(func() {
		paramVar := make(map[string]*ReqParamV)

		request := c.request
		for n, v := range request.URL.Query() {
			paramVar[n] = NewReqParamV(v)
		}

		routerParams := c.routerParams
		if routerParams != nil {
			for n, v := range routerParams {
				paramVar[n] = v
			}
		}

		c.paramVar = NewReqParamHelper(paramVar)
	})

	return c.paramVar
}

func (c *Request) PostVar() *ReqParamHelper {
	c.postOnce.Do(func() {
		postVar := make(map[string]*ReqParamV)
		request := c.request
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

func (c *Request) FileVar() map[string]*File {
	c.fileOnce.Do(func() {
		request := c.request
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

func (c *Request) IsAjax() bool {
	return c.request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (c *Request) IsPost() bool {
	return c.request.Method == "POST"
}
