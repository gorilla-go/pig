package pig

import "strconv"

type ReqParams map[string]*ReqParamV

type ReqParamV struct {
	v []*ReqParamAtom
}

type ReqParamAtom struct {
	v string
}

func NewReqParamAtom(v string) *ReqParamAtom {
	return &ReqParamAtom{
		v: v,
	}
}

func NewReqParamV(value []string) *ReqParamV {
	var paramsAtom []*ReqParamAtom
	for _, vr := range value {
		paramsAtom = append(paramsAtom, NewReqParamAtom(vr))
	}
	return &ReqParamV{
		v: paramsAtom,
	}
}

func (r *ReqParamAtom) ToString() string {
	return r.v
}

func (r *ReqParamAtom) ToInt() int {
	i, err := strconv.Atoi(r.v)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *ReqParamAtom) ToInt64() int64 {
	i, err := strconv.ParseInt(r.v, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *ReqParamAtom) ToFloat64() float64 {
	i, err := strconv.ParseFloat(r.v, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *ReqParamAtom) ToBool() bool {
	i, err := strconv.ParseBool(r.v)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *ReqParamAtom) ToBytes() []byte {
	return []byte(r.v)
}

func (r *ReqParamV) ToString() string {
	return r.v[0].ToString()
}

func (r *ReqParamV) ToInt() int {
	return r.v[0].ToInt()
}

func (r *ReqParamV) ToInt64() int64 {
	return r.v[0].ToInt64()
}

func (r *ReqParamV) ToFloat64() float64 {
	return r.v[0].ToFloat64()
}

func (r *ReqParamV) ToBool() bool {
	return r.v[0].ToBool()
}

func (r *ReqParamV) ToBytes() []byte {
	return r.v[0].ToBytes()
}

func (r *ReqParamV) ToReqParamAtoms() []*ReqParamAtom {
	return r.v
}
