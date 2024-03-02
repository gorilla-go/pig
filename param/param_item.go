package param

import (
	"strconv"
	"strings"
)

type RequestParamItem struct {
	v string
}

func NewRequestParamItem(v string) *RequestParamItem {
	return &RequestParamItem{
		v: v,
	}
}

func (r *RequestParamItem) String() string {
	return strings.TrimSpace(r.v)
}

func (r *RequestParamItem) Int() int {
	i, err := strconv.Atoi(r.v)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *RequestParamItem) Int64() int64 {
	i, err := strconv.ParseInt(r.v, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *RequestParamItem) Float64() float64 {
	i, err := strconv.ParseFloat(r.v, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *RequestParamItem) Bool() bool {
	i, err := strconv.ParseBool(r.v)
	if err != nil {
		panic(err)
	}
	return i
}
