package foundation

import (
	"golang.org/x/exp/constraints"
)

type TreeMap[T constraints.Ordered, V any] struct {
	K []T
	M map[T]V
}

func (t *TreeMap[T, V]) Put(key T, value V) {
	t.K = append(t.K, key)
	t.M[key] = value

	for i := 0; i < len(t.K)-1; i++ {
		for j := 0; j < len(t.K)-i-1; j++ {
			if t.K[j] > t.K[j+1] {
				t.K[j], t.K[j+1] = t.K[j+1], t.K[j]
			}
		}
	}
}

func (t *TreeMap[T, V]) Get(key T) V {
	return t.M[key]
}

func (t *TreeMap[T, V]) Keys() []T {
	return t.K
}

func (t *TreeMap[T, V]) Values() []V {
	var values []V
	for _, k := range t.K {
		values = append(values, t.M[k])
	}
	return values
}

func (t *TreeMap[T, V]) Size() int {
	return len(t.K)
}

func (t *TreeMap[T, V]) Remove(key T) {
	delete(t.M, key)
	var newKeys []T
	for _, k := range t.K {
		if k != key {
			newKeys = append(newKeys, k)
		}
	}
	t.K = newKeys
}

func (t *TreeMap[T, V]) Clear() {
	t.K = nil
	t.M = make(map[T]V)
}

func (t *TreeMap[T, V]) ContainsKey(key T) bool {
	_, ok := t.M[key]
	return ok
}

func (t *TreeMap[T, V]) IsEmpty() bool {
	return len(t.K) == 0
}

func (t *TreeMap[T, V]) ForEach(f func(key T, value V) bool) {
	for _, k := range t.K {
		r := f(k, t.M[k])
		if !r {
			break
		}
	}
}
