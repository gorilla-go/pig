package param

import "github.com/gorilla-go/pig/foundation/mapping"

type RequestParamPairs[V any] mapping.HashMap[string, V]

func NewRequestParamPairs[V any]() *RequestParamPairs[V] {
	p := mapping.NewHashMap[string, V]()
	return (*RequestParamPairs[V])(p)
}

func (p *RequestParamPairs[V]) Raw() *mapping.HashMap[string, V] {
	return (*mapping.HashMap[string, V])(p)
}
