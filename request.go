package pig

import (
	"bytes"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/samber/lo"
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
	paramVar     *ReqParamHelper
	paramOnce    sync.Once
	postVar      *ReqParamHelper
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

func (c *Request) ParamVar() *ReqParamHelper {
	c.paramOnce.Do(func() {
		paramVar := make(map[string]*ReqParamV)

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
						paramVar[k].v = append(paramVar[k].v, NewReqParamAtom(v))
						continue
					}
					paramVar[k] = NewReqParamV([]string{v})
				}
			}
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

		contentType := request.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(request.Body)
			if err != nil {
				panic(err)
			}

			j := make(map[string]string)
			err = json.Unmarshal(buf.Bytes(), &j)
			if err != nil {
				panic(err)
			}

			for n, v := range j {
				postVar[n] = NewReqParamV([]string{v})
			}
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
						postVar[formName] = NewReqParamV([]string{buf.String()})
					}
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

func (c *Request) IsPost() bool {
	return c.request.Method == "POST"
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
			injector(
				field.Type,
				(c.PostVar().Raw())[tag],
				unsafe.Pointer(v.Field(i).UnsafeAddr()),
			)
		}

		tag = strings.TrimSpace(field.Tag.Get("query"))
		if tag != "" {
			injector(
				field.Type,
				(c.ParamVar().Raw())[tag],
				unsafe.Pointer(v.Field(i).UnsafeAddr()),
			)
		}
	}
}

func injector(tp reflect.Type, val *ReqParamV, at unsafe.Pointer) {
	if val != nil && len(val.v) > 0 && canInjected(tp.Kind()) {
		reflect.NewAt(tp, at).Elem().Set(
			reflect.ValueOf(convertStringToKind(val, tp)),
		)
	}
}

func canInjected(k reflect.Kind) bool {
	return lo.IndexOf([]reflect.Kind{
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Interface,
		reflect.Slice,
		reflect.String,
	}, k) != -1
}

func convertStringToKind(s *ReqParamV, k reflect.Type) any {
	switch k.Kind() {
	case reflect.Bool:
		return s.Bool()
	case reflect.Int:
		return s.Int()
	case reflect.Int8:
		return int8(s.Int())
	case reflect.Int16:
		return int16(s.Int())
	case reflect.Int32:
		return int32(s.Int())
	case reflect.Int64:
		return s.Int64()
	case reflect.Uint:
		return uint(s.Int())
	case reflect.Uint8:
		return uint8(s.Int())
	case reflect.Uint16:
		return uint16(s.Int())
	case reflect.Uint32:
		return uint32(s.Int())
	case reflect.Uint64:
		return uint64(s.Int())
	case reflect.Float32:
		return float32(s.Float64())
	case reflect.Float64:
		return s.Float64()
	case reflect.Slice:
		sType := reflect.SliceOf(k.Elem())
		l := len(s.v)
		slice := reflect.MakeSlice(sType, l, l)
		for i := 0; i < l; i++ {
			itemType := slice.Index(i).Type()
			if !canInjected(slice.Index(i).Type().Kind()) {
				panic("unsupported inject type: []" + itemType.String())
			}
			slice.Index(i).Set(reflect.ValueOf(
				convertStringToKind(
					NewReqParamV([]string{s.v[i].String()}),
					itemType,
				),
			))
		}
		return slice.Interface()
	case reflect.Interface:
		return s.String()
	case reflect.String:
		return s.String()
	}

	return nil
}
