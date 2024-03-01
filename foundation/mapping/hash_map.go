package mapping

import (
	"fmt"
)

type HashMap[K comparable, V any] struct {
	raw map[K]V
}

func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{raw: make(map[K]V)}
}

func (m *HashMap[K, V]) Get(key K) (V, bool) {
	val, ok := m.raw[key]
	return val, ok
}

func (m *HashMap[K, V]) MustGet(key K) V {
	val, ok := m.raw[key]
	if !ok {
		panic(fmt.Sprintf("key: %v not found", key))
	}
	return val
}

func (m *HashMap[K, V]) Set(key K, value V) {
	m.raw[key] = value
}

func (m *HashMap[K, V]) Delete(key K) {
	delete(m.raw, key)
}

func (m *HashMap[K, V]) ContainsKey(key K) bool {
	_, ok := m.raw[key]
	return ok
}

func (m *HashMap[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.raw))
	for key := range m.raw {
		keys = append(keys, key)
	}
	return keys
}

func (m *HashMap[K, V]) Values() []V {
	values := make([]V, 0, len(m.raw))
	for _, value := range m.raw {
		values = append(values, value)
	}
	return values
}

func (m *HashMap[K, V]) Len() int {
	return len(m.raw)
}

func (m *HashMap[K, V]) IsEmpty() bool {
	return len(m.raw) == 0
}

func (m *HashMap[K, V]) ForEach(f func(key K, value V) bool) {
	for k, v := range m.raw {
		r := f(k, v)
		if r == false {
			break
		}
	}
}
