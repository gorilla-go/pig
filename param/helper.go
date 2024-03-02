package param

import (
	"errors"
	"github.com/gorilla-go/pig/foundation"
	"github.com/gorilla-go/pig/foundation/mapping"
)

var ParamsUnSet = errors.New("param unset")

type IRequestInput interface {
	*RequestParamItems | *File
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

func (h *Helper[V]) Slice(s string) ([]*RequestParamItem, error) {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems).Slice(), nil
	}
	return nil, ParamsUnSet
}

func (h *Helper[V]) MustSlice(s string) []*RequestParamItem {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems).Slice()
	}
	panic(ParamsUnSet)
}

func (h *Helper[V]) Int(s string, def ...int) (int, error) {
	ret := foundation.DefaultParam(def, 0)
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems).Int(), nil
	}
	return ret, ParamsUnSet
}

func (h *Helper[V]) MustInt(s string) int {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems).Int()
	}
	panic(ParamsUnSet)
}

func (h *Helper[V]) Int64(s string, def ...int64) (int64, error) {
	ret := foundation.DefaultParam(def, int64(0))
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems).Int64(), nil
	}
	return ret, ParamsUnSet
}

func (h *Helper[V]) MustInt64(s string) int64 {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems).Int64()
	}
	panic(ParamsUnSet)
}

func (h *Helper[V]) Float64(s string, def ...float64) (float64, error) {
	ret := foundation.DefaultParam(def, float64(0))
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems).Float64(), nil
	}
	return ret, ParamsUnSet
}

func (h *Helper[V]) MustFloat64(s string) float64 {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems).Float64()
	}
	panic(ParamsUnSet)
}

func (h *Helper[V]) Bool(s string, def ...bool) (bool, error) {
	ret := foundation.DefaultParam(def, false)
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems).Bool(), nil
	}
	return ret, ParamsUnSet
}

func (h *Helper[V]) MustBool(s string) bool {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems).Bool()
	}
	panic(ParamsUnSet)
}

func (h *Helper[V]) String(s string, def ...string) (string, error) {
	ret := foundation.DefaultParam(def, "")
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems).String(), nil
	}
	return ret, ParamsUnSet
}

func (h *Helper[V]) MustString(s string) string {
	paramPairs := h.Raw()
	if paramPairs.ContainsKey(s) {
		return any(paramPairs.MustGet(s)).(*RequestParamItems).String()
	}
	panic(ParamsUnSet)
}
