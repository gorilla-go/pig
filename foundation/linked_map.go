package foundation

type LinkedHashMap[K comparable, T any] struct {
	K []K
	M map[K]T
}

func (l *LinkedHashMap[K, T]) Put(key K, value T) {
	l.K = append(l.K, key)
	l.M[key] = value
}

func (l *LinkedHashMap[K, T]) Get(key K) T {
	return l.M[key]
}

func (l *LinkedHashMap[K, T]) Keys() []K {
	return l.K
}

func (l *LinkedHashMap[K, T]) Values() []T {
	var values []T
	for _, k := range l.K {
		values = append(values, l.M[k])
	}
	return values
}

func (l *LinkedHashMap[K, T]) Size() int {
	return len(l.K)
}

func (l *LinkedHashMap[K, T]) Remove(key K) {
	delete(l.M, key)
	var newKeys []K
	for _, k := range l.K {
		if k != key {
			newKeys = append(newKeys, k)
		}
	}
	l.K = newKeys
}

func (l *LinkedHashMap[K, T]) Clear() {
	l.K = nil
	l.M = make(map[K]T)
}

func (l *LinkedHashMap[K, T]) ContainsKey(key K) bool {
	_, ok := l.M[key]
	return ok
}

func (l *LinkedHashMap[K, T]) IsEmpty() bool {
	return len(l.K) == 0
}

func (l *LinkedHashMap[K, T]) ForEach(f func(key K, value T) bool) {
	for _, k := range l.K {
		r := f(k, l.M[k])
		if r == false {
			break
		}
	}
}
