package mapping

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"sort"
)

type TreeMap[K constraints.Ordered, V any] struct {
	K []K
	M map[K]V
}

func NewTreeMap[K constraints.Ordered, V any]() *TreeMap[K, V] {
	return &TreeMap[K, V]{
		K: make([]K, 0),
		M: make(map[K]V),
	}
}

func (t *TreeMap[K, V]) Put(key K, value V) {
	index := sort.Search(len(t.K), func(i int) bool {
		return t.K[i] >= key
	})
	t.K = append(t.K, key)
	copy(t.K[index+1:], t.K[index:])
	t.K[index] = key
	t.M[key] = value
}

func (t *TreeMap[K, V]) Get(key K) (V, bool) {
	v, ok := t.M[key]
	return v, ok
}

func (t *TreeMap[K, V]) MustGet(key K) V {
	v, ok := t.M[key]
	if !ok {
		panic(fmt.Sprintf("key: %v not found", key))
	}
	return v
}

func (t *TreeMap[K, V]) Keys() []K {
	return t.K
}

func (t *TreeMap[K, V]) Values() []V {
	var values []V
	for _, k := range t.K {
		values = append(values, t.M[k])
	}
	return values
}

func (t *TreeMap[K, V]) Len() int {
	return len(t.K)
}

func (t *TreeMap[K, V]) Remove(key K) {
	index := sort.Search(len(t.K), func(i int) bool { return t.K[i] >= key })
	if index < len(t.K) && t.K[index] == key {
		copy(t.K[index:], t.K[index+1:])
		t.K = t.K[:len(t.K)-1]
	}
	delete(t.M, key)
}

func (t *TreeMap[K, V]) Clear() {
	t.M = make(map[K]V)
	t.K = t.K[:0]
}

func (t *TreeMap[K, V]) ContainsKey(key K) bool {
	_, ok := t.M[key]
	return ok
}

func (t *TreeMap[K, V]) IsEmpty() bool {
	return len(t.K) == 0
}

func (t *TreeMap[K, V]) ForEach(f func(key K, value V) bool) {
	for _, k := range t.K {
		r := f(k, t.M[k])
		if !r {
			break
		}
	}
}
