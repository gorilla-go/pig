package pig

import (
	"github.com/samber/lo"
)

type Pipeline[T any] struct {
	injector  T
	pipelines []func(T, func(T))
}

func NewPipeline[T any]() *Pipeline[T] {
	return &Pipeline[T]{}
}

func (h *Pipeline[T]) Send(i T) *Pipeline[T] {
	h.injector = i
	return h
}

func (h *Pipeline[T]) Through(unit func(T, func(T))) *Pipeline[T] {
	h.pipelines = append(h.pipelines, unit)
	return h
}

func (h *Pipeline[T]) Then(p func(T)) {
	f := lo.Reduce(
		h.pipelines,
		func(
			f func(T),
			f2 func(
				T,
				func(T),
			),
			i int,
		) func(T) {
			return func(i T) {
				f2(i, f)
			}
		},
		p,
	)
	f(h.injector)
}
