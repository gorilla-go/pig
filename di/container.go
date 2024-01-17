package di

import "sync"

type typeName string

type Container struct {
	eager      map[typeName]any
	eagerNamed map[string]typeName

	lazy      map[typeName]any
	lazyNamed map[string]typeName

	rebuild      map[typeName]any
	rebuildNamed map[string]typeName

	locker *sync.Mutex
}

func New() *Container {
	return &Container{
		eager:        make(map[typeName]any),
		eagerNamed:   make(map[string]typeName),
		lazy:         make(map[typeName]any),
		lazyNamed:    make(map[string]typeName),
		rebuild:      make(map[typeName]any),
		rebuildNamed: make(map[string]typeName),
		locker:       &sync.Mutex{},
	}
}
