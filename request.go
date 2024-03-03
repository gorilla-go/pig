package pig

import (
	"bytes"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla-go/pig/foundation/injector"
	"github.com/gorilla-go/pig/param"
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
	routerParams *param.RequestParamPairs[*param.RequestParamItems[string]]
	paramVar     *param.Helper[*param.RequestParamItems[string]]
	paramOnce    sync.Once
	postVar      *param.Helper[*param.RequestParamItems[string]]
	postOnce     sync.Once
	fileVar      *param.Helper[*param.RequestParamItems[*param.File]]
	fileOnce     sync.Once
}

func NewRequest(
	req *http.Request,
	routerParams *param.RequestParamPairs[*param.RequestParamItems[string]],
) *Request {
	return &Request{
		request:      req,
		routerParams: routerParams,
	}
}

func (c *Request) Raw() *http.Request {
	return c.request
}

func (c *Request) ParamVar() *param.Helper[*param.RequestParamItems[string]] {
	c.paramOnce.Do(func() {
		paramVarPairs := param.NewRequestParamPairs[*param.RequestParamItems[string]]()
		rawQuery := c.request.URL.RawQuery
		if rawQuery != "" {
			kvGroup := strings.Split(rawQuery, "&")
			for _, kv := range kvGroup {
				kvArr := strings.SplitN(kv, "=", 2)
				k := strings.TrimSpace(kvArr[0])
				v := strings.TrimSpace(kvArr[1])
				paramVarPairs.Raw().Set(k, param.NewRequestParamItems[string]([]string{v}))
			}
		}

		if c.routerParams != nil {
			c.routerParams.Raw().ForEach(func(k string, v *param.RequestParamItems[string]) bool {
				if paramVarPairs.Raw().ContainsKey(k) {
					paramVarPairs.Raw().MustGet(k).SetParams(
						append(
							paramVarPairs.Raw().MustGet(k).GetParams(),
							param.NewRequestParamItem(v.String()),
						),
					)
					return true
				}
				paramVarPairs.Raw().Set(k, v)
				return true
			})
		}

		c.paramVar = param.NewParamHelper[*param.RequestParamItems[string]](paramVarPairs)
	})

	return c.paramVar
}

func (c *Request) PostVar() *param.Helper[*param.RequestParamItems[string]] {
	c.postOnce.Do(func() {
		postVarPairs := param.NewRequestParamPairs[*param.RequestParamItems[string]]()
		request := c.request
		err := request.ParseForm()
		if err != nil {
			panic(err)
		}
		for n, v := range request.PostForm {
			postVarPairs.Raw().Set(n, param.NewRequestParamItems[string](v))
		}

		contentType := request.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			panic("not support json post.")
		}

		if postVarPairs.Raw().Len() == 0 &&
			strings.Contains(contentType, "multipart/form-data") {
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
					postVarPairs.Raw().Set(
						formName,
						param.NewRequestParamItems[string]([]string{buf.String()}),
					)
				}
			}
		}

		c.postVar = param.NewParamHelper(postVarPairs)
	})
	return c.postVar
}

func (c *Request) FileVar() *param.Helper[*param.RequestParamItems[*param.File]] {
	c.fileOnce.Do(func() {
		request := c.request

		if !strings.Contains(request.Header.Get("Content-Type"), "multipart/form-data") {
			return
		}

		multipartReader, err := request.MultipartReader()
		if err != nil {
			return
		}

		requestPairs := param.NewRequestParamPairs[*param.RequestParamItems[*param.File]]()
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
				tmp := os.TempDir()
				if !strings.HasSuffix(tmp, string(filepath.Separator)) {
					tmp += string(filepath.Separator)
				}
				filename := tmp + string(fileId) + ext

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

				f := &param.File{
					FilePath:    filename,
					ContentType: part.Header.Get("Content-Type"),
					Basename:    fileName,
					Ext:         ext,
				}
				if requestPairs.Raw().ContainsKey(formName) {
					var fileSlice []*param.File
					for _, r := range requestPairs.Raw().MustGet(formName).GetParams() {
						fileSlice = append(fileSlice, r.File())
					}
					requestPairs.Raw().Set(
						formName,
						param.NewRequestParamItems[*param.File](append(fileSlice, f)),
					)
					continue
				}

				requestPairs.Raw().Set(
					formName,
					param.NewRequestParamItems[*param.File]([]*param.File{f}),
				)
			}
		}

		c.fileVar = param.NewParamHelper(requestPairs)
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
			injector.RequestInjector(
				field.Type,
				c.postVar.Raw().MustGet(tag),
				unsafe.Pointer(v.Field(i).UnsafeAddr()),
			)
		}

		tag = strings.TrimSpace(field.Tag.Get("query"))
		if tag != "" {
			injector.RequestInjector(
				field.Type,
				c.paramVar.Raw().MustGet(tag),
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
