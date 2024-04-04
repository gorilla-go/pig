package param

import (
	"errors"
	"fmt"
	"github.com/gorilla-go/pig/foundation"
	"github.com/gorilla-go/pig/foundation/mapping"
)

var ParamsUnSet = errors.New("param unset")

type IRequestInput interface {
	*RequestParamItems[string] | *RequestParamItems[*File]
}

type Helper[V IRequestInput] struct {
	pairs *RequestParamPairs[V]
}

func NewParamHelper[V IRequestInput](m *RequestParamPairs[V]) *Helper[V] {
	return &Helper[V]{pairs: m}
}

func (h *Helper[V]) Raw() *mapping.HashMap[string, V] {
	return h.pairs.Raw()
}

func (h *Helper[V]) Slice(s string) (r []*RequestParamItem[string], err error) {
	r = nil
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		defer func() {
			if errno := recover(); errno != nil {
				err = errors.New(fmt.Sprintf("%v", errno))
			}
		}()
		return any(paramPairs.MustGet(s)).(*RequestParamItems[string]).Slice(), nil
	}
	return r, ParamsUnSet
}

func (h *Helper[V]) MustSlice(s string) []*RequestParamItem[string] {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems[string]).Slice()
	}
	panic(ParamsUnSet)
}

func (h *Helper[V]) FileSlice(s string) (f []*RequestParamItem[*File], err error) {
	f = nil
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		defer func() {
			if errno := recover(); errno != nil {
				err = errors.New(fmt.Sprintf("%v", errno))
			}
		}()
		return any(paramPairs.MustGet(s)).(*RequestParamItems[*File]).Slice(), nil
	}
	return nil, ParamsUnSet
}

func (h *Helper[V]) MustFileSlice(s string) []*RequestParamItem[*File] {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems[*File]).Slice()
	}
	panic(ParamsUnSet)
}

func (h *Helper[V]) Int(s string, def ...int) (ret int, err error) {
	ret = foundation.DefaultParam(def, 0)
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		defer func() {
			if errno := recover(); errno != nil {
				err = errors.New(fmt.Sprintf("%v", errno))
			}
		}()
		return any(paramPairs.MustGet(s)).(*RequestParamItems[string]).Int(), nil
	}
	return ret, ParamsUnSet
}

func (h *Helper[V]) MustInt(s string) int {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems[string]).Int()
	}
	panic(ParamsUnSet)
}

func (h *Helper[V]) Int64(s string, def ...int64) (ret int64, err error) {
	ret = foundation.DefaultParam(def, int64(0))
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		defer func() {
			if errno := recover(); errno != nil {
				err = errors.New(fmt.Sprintf("%v", errno))
			}
		}()
		return any(paramPairs.MustGet(s)).(*RequestParamItems[string]).Int64(), nil
	}
	return ret, ParamsUnSet
}

func (h *Helper[V]) MustInt64(s string) int64 {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems[string]).Int64()
	}
	panic(ParamsUnSet)
}

func (h *Helper[V]) Float64(s string, def ...float64) (ret float64, err error) {
	ret = foundation.DefaultParam(def, float64(0))
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		defer func() {
			if errno := recover(); errno != nil {
				err = errors.New(fmt.Sprintf("%v", errno))
			}
		}()
		return any(paramPairs.MustGet(s)).(*RequestParamItems[string]).Float64(), nil
	}
	return ret, ParamsUnSet
}

func (h *Helper[V]) MustFloat64(s string) float64 {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems[string]).Float64()
	}
	panic(ParamsUnSet)
}

func (h *Helper[V]) Bool(s string, def ...bool) (ret bool, err error) {
	ret = foundation.DefaultParam(def, false)
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		defer func() {
			if errno := recover(); errno != nil {
				err = errors.New(fmt.Sprintf("%v", errno))
			}
		}()
		return any(paramPairs.MustGet(s)).(*RequestParamItems[string]).Bool(), nil
	}
	return ret, ParamsUnSet
}

func (h *Helper[V]) MustBool(s string) bool {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems[string]).Bool()
	}
	panic(ParamsUnSet)
}

func (h *Helper[V]) String(s string, def ...string) (ret string, err error) {
	ret = foundation.DefaultParam(def, "")
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		defer func() {
			if errno := recover(); errno != nil {
				err = errors.New(fmt.Sprintf("%v", errno))
			}
		}()
		return any(paramPairs.MustGet(s)).(*RequestParamItems[string]).String(), nil
	}
	return ret, ParamsUnSet
}

func (h *Helper[V]) MustString(s string) string {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems[string]).String()
	}
	panic(ParamsUnSet)
}

func (h *Helper[V]) File(s string) (f *File, err error) {
	f = nil
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		defer func() {
			if errno := recover(); errno != nil {
				err = errors.New(fmt.Sprintf("%v", errno))
			}
		}()
		return any(paramPairs.MustGet(s)).(*RequestParamItems[*File]).File(), nil
	}
	return nil, ParamsUnSet
}

func (h *Helper[V]) MustFile(s string) *File {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems[*File]).File()
	}
	panic(ParamsUnSet)
}
