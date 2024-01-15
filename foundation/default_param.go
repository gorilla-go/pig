package foundation

func DefaultParam[T any](g []T, v T) T {
	if g == nil || len(g) == 0 {
		return v
	}

	return g[0]
}
