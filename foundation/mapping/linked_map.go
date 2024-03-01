package mapping

type LinkedHashMap[K comparable, V any] struct {
	K []K
	M map[K]V
}

func NewLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{
		K: make([]K, 0),
		M: make(map[K]V),
	}
}

func (l *LinkedHashMap[K, V]) Put(key K, value V) {
	l.K = append(l.K, key)
	l.M[key] = value
}

func (l *LinkedHashMap[K, V]) Get(key K) V {
	return l.M[key]
}

func (l *LinkedHashMap[K, V]) Keys() []K {
	return l.K
}

func (l *LinkedHashMap[K, V]) Values() []V {
	var values []V
	for _, k := range l.K {
		values = append(values, l.M[k])
	}
	return values
}

func (l *LinkedHashMap[K, V]) Size() int {
	return len(l.K)
}

func (l *LinkedHashMap[K, V]) Remove(key K) {
	delete(l.M, key)
	for i, k := range l.K {
		if k == key {
			l.K = append(l.K[:i], l.K[i+1:]...)
			break
		}
	}
}

func (l *LinkedHashMap[K, V]) Clear() {
	l.K = nil
	l.M = make(map[K]V)
}

func (l *LinkedHashMap[K, V]) ContainsKey(key K) bool {
	_, ok := l.M[key]
	return ok
}

func (l *LinkedHashMap[K, V]) IsEmpty() bool {
	return len(l.K) == 0
}

func (l *LinkedHashMap[K, V]) ForEach(f func(key K, value V) bool) {
	for _, k := range l.K {
		r := f(k, l.M[k])
		if r == false {
			break
		}
	}
}
