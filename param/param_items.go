package param

type RequestParamItems[T IRequestParamItemType] struct {
	v []*RequestParamItem[T]
}

func NewRequestParamItems[T IRequestParamItemType](value []T) *RequestParamItems[T] {
	var paramsAtom []*RequestParamItem[T]
	for _, vr := range value {
		paramsAtom = append(paramsAtom, NewRequestParamItem(vr))
	}
	return &RequestParamItems[T]{
		v: paramsAtom,
	}
}

func (r *RequestParamItems[T]) Slice() []*RequestParamItem[T] {
	return r.v
}

func (r *RequestParamItems[T]) String() string {
	return r.v[0].String()
}

func (r *RequestParamItems[T]) Int() int {
	return r.v[0].Int()
}

func (r *RequestParamItems[T]) Int64() int64 {
	return r.v[0].Int64()
}

func (r *RequestParamItems[T]) Float64() float64 {
	return r.v[0].Float64()
}

func (r *RequestParamItems[T]) Bool() bool {
	return r.v[0].Bool()
}

func (r *RequestParamItems[T]) File() *File {
	return r.v[0].File()
}

func (r *RequestParamItems[T]) GetParams() []*RequestParamItem[T] {
	return r.v
}

func (r *RequestParamItems[T]) SetParams(rpa []*RequestParamItem[T]) *RequestParamItems[T] {
	r.v = rpa
	return r
}
