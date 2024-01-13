package pig

import (
	"strconv"
	"strings"
)

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

func (r *ReqParamAtom) String() string {
	return r.v
}

func (r *ReqParamAtom) TrimString() string {
	return strings.TrimSpace(r.v)
}

func (r *ReqParamAtom) Int() int {
	i, err := strconv.Atoi(r.v)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *ReqParamAtom) Int64() int64 {
	i, err := strconv.ParseInt(r.v, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *ReqParamAtom) Float64() float64 {
	i, err := strconv.ParseFloat(r.v, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *ReqParamAtom) Bool() bool {
	i, err := strconv.ParseBool(r.v)
	if err != nil {
		panic(err)
	}
	return i
}

func (r *ReqParamAtom) Bytes() []byte {
	return []byte(r.v)
}

func (r *ReqParamV) String() string {
	return r.v[0].String()
}

func (r *ReqParamV) TrimString() string {
	return r.v[0].TrimString()
}

func (r *ReqParamV) Int() int {
	return r.v[0].Int()
}

func (r *ReqParamV) Int64() int64 {
	return r.v[0].Int64()
}

func (r *ReqParamV) Float64() float64 {
	return r.v[0].Float64()
}

func (r *ReqParamV) Bool() bool {
	return r.v[0].Bool()
}

func (r *ReqParamV) Bytes() []byte {
	return r.v[0].Bytes()
}

func (r *ReqParamV) ReqParamAtoms() []*ReqParamAtom {
	return r.v
}
