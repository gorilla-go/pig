package param

type RequestParamItems struct {
	v []*RequestParamItem
}

func NewRequestParamItems(value []string) *RequestParamItems {
	var paramsAtom []*RequestParamItem
	for _, vr := range value {
		paramsAtom = append(paramsAtom, NewRequestParamItem(vr))
	}
	return &RequestParamItems{
		v: paramsAtom,
	}
}

func (r *RequestParamItems) Slice() []*RequestParamItem {
	return r.v
}

func (r *RequestParamItems) String() string {
	return r.v[0].String()
}

func (r *RequestParamItems) Int() int {
	return r.v[0].Int()
}

func (r *RequestParamItems) Int64() int64 {
	return r.v[0].Int64()
}

func (r *RequestParamItems) Float64() float64 {
	return r.v[0].Float64()
}

func (r *RequestParamItems) Bool() bool {
	return r.v[0].Bool()
}

func (r *RequestParamItems) GetParams() []*RequestParamItem {
	return r.v
}

func (r *RequestParamItems) SetParams(rpa []*RequestParamItem) *RequestParamItems {
	r.v = rpa
	return r
}
