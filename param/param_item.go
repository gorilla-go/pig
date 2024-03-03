package param

import (
	"strconv"
	"strings"
)

type IRequestParamItemType interface {
	string | *File
}

type RequestParamItem[T IRequestParamItemType] struct {
	v T
}

func NewRequestParamItem[T IRequestParamItemType](v T) *RequestParamItem[T] {
	return &RequestParamItem[T]{
		v: v,
	}
}

func (r *RequestParamItem[T]) String() string {
	return strings.TrimSpace(any(r.v).(string))
}

func (r *RequestParamItem[T]) Int() int {
	i, err := strconv.Atoi(any(r.v).(string))
	if err != nil {
		panic(err)
	}
	return i
}

func (r *RequestParamItem[T]) Int64() int64 {
	i, err := strconv.ParseInt(any(r.v).(string), 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *RequestParamItem[T]) Float64() float64 {
	i, err := strconv.ParseFloat(any(r.v).(string), 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *RequestParamItem[T]) Bool() bool {
	i, err := strconv.ParseBool(any(r.v).(string))
	if err != nil {
		panic(err)
	}
	return i
}

func (r *RequestParamItem[T]) File() *File {
	return (any)(r.v).(*File)
}
