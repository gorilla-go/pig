package mapping

type Map[K comparable, V any] struct {
	raw map[K]V
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{raw: make(map[K]V)}
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	val, ok := m.raw[key]
	return val, ok
}

func (m *Map[K, V]) Set(key K, value V) {
	m.raw[key] = value
}

func (m *Map[K, V]) Delete(key K) {
	delete(m.raw, key)
}

func (m *Map[K, V]) ContainsKey(key K) bool {
	_, ok := m.raw[key]
	return ok
}

func (m *Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.raw))
	for key := range m.raw {
		keys = append(keys, key)
	}
	return keys
}

func (m *Map[K, V]) Values() []V {
	values := make([]V, 0, len(m.raw))
	for _, value := range m.raw {
		values = append(values, value)
	}
	return values
}

func (m *Map[K, V]) ForEach(f func(key K, value V) bool) {
	for k, v := range m.raw {
		r := f(k, v)
		if r == false {
			break
		}
	}
}
