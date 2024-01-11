package pig

import (
	"github.com/samber/do"
	"github.com/samber/lo"
)

type Pipeline struct {
	injector  *do.Injector
	pipelines []func(*do.Injector, func(*do.Injector))
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (h *Pipeline) Send(i *do.Injector) *Pipeline {
	h.injector = i
	return h
}

func (h *Pipeline) Through(unit func(*do.Injector, func(*do.Injector))) *Pipeline {
	h.pipelines = append(h.pipelines, unit)
	return h
}

func (h *Pipeline) Then(p func(*do.Injector)) {
	f := lo.Reduce(
		h.pipelines,
		func(
			f func(*do.Injector),
			f2 func(
				*do.Injector,
				func(*do.Injector),
			),
			i int,
		) func(*do.Injector) {
			return func(i *do.Injector) {
				f2(i, f)
			}
		},
		p,
	)
	f(h.injector)
}
