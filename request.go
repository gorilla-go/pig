package pig

import (
	"bytes"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla-go/pig/foundation"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

type Request struct {
	request      *http.Request
	routerParams RouterParams
	paramVar     *foundation.ReqParamHelper
	paramOnce    sync.Once
	postVar      *foundation.ReqParamHelper
	postOnce     sync.Once
	fileVar      map[string]*File
	fileOnce     sync.Once
}

func NewRequest(req *http.Request, routerParams RouterParams) *Request {
	return &Request{
		request:      req,
		routerParams: routerParams,
	}
}

func (c *Request) Raw() *http.Request {
	return c.request
}

func (c *Request) ParamVar() *foundation.ReqParamHelper {
	c.paramOnce.Do(func() {
		paramVar := make(map[string]*foundation.ReqParamV)

		request := c.request
		rawQuery := request.URL.RawQuery
		if rawQuery != "" {
			kvGroup := strings.Split(rawQuery, "&")
			for _, kv := range kvGroup {
				kvArr := strings.Split(kv, "=")
				if len(kvArr) == 2 {
					k := strings.TrimSpace(kvArr[0])
					v := strings.TrimSpace(kvArr[1])
					if _, ok := paramVar[k]; ok {
						paramVar[k].SetReqParamAtoms(
							append(paramVar[k].ReqParamAtoms(), foundation.NewReqParamAtom(v)),
						)
						continue
					}
					paramVar[k] = foundation.NewReqParamV([]string{v})
				}
			}
		}

		routerParams := c.routerParams
		if routerParams != nil {
			for n, v := range routerParams {
				paramVar[n] = v
			}
		}

		c.paramVar = foundation.NewReqParamHelper(paramVar)
	})

	return c.paramVar
}

func (c *Request) PostVar() *foundation.ReqParamHelper {
	c.postOnce.Do(func() {
		postVar := make(map[string]*foundation.ReqParamV)
		request := c.request
		err := request.ParseForm()
		if err != nil {
			panic(err)
		}
		for n, v := range request.PostForm {
			postVar[n] = foundation.NewReqParamV(v)
		}

		contentType := request.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			panic("not support json post.")
		}

		if strings.Contains(contentType, "multipart/form-data") {
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
						postVar[formName] = foundation.NewReqParamV([]string{buf.String()})
					}
				}
			}
		}

		c.postVar = foundation.NewReqParamHelper(postVar)
	})
	return c.postVar
}

func (c *Request) FileVar() map[string]*File {
	c.fileOnce.Do(func() {
		request := c.request
		c.fileVar = make(map[string]*File)

		if !strings.Contains(request.Header.Get("Content-Type"), "multipart/form-data") {
			return
		}

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

func (c *Request) IsGet() bool {
	return c.request.Method == "GET"
}

func (c *Request) IsPost() bool {
	return c.request.Method == "POST"
}

func (c *Request) IsPut() bool {
	return c.request.Method == "PUT"
}

func (c *Request) IsDelete() bool {
	return c.request.Method == "DELETE"
}

func (c *Request) IsOption() bool {
	return c.request.Method == "OPTION"
}

func (c *Request) Bind(t any) {
	v := reflect.ValueOf(t)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic("bind target must be struct or struct pointer.")
	}

	// bind param
	for i := 0; i < v.NumField(); i++ {
		tag := ""
		field := v.Type().Field(i)

		tag = strings.TrimSpace(field.Tag.Get("form"))
		if tag != "" {
			foundation.RequestInjector(
				field.Type,
				(c.PostVar().Raw())[tag],
				unsafe.Pointer(v.Field(i).UnsafeAddr()),
			)
		}

		tag = strings.TrimSpace(field.Tag.Get("query"))
		if tag != "" {
			foundation.RequestInjector(
				field.Type,
				(c.ParamVar().Raw())[tag],
				unsafe.Pointer(v.Field(i).UnsafeAddr()),
			)
		}
	}
}

func (c *Request) JsonBind(t any) {
	ct := c.request.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		return
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.request.Body)
	if err != nil {
		panic(err)
	}

	if buf.Len() == 0 {
		return
	}

	err = json.Unmarshal(buf.Bytes(), t)
	if err != nil {
		panic(err)
	}
}
