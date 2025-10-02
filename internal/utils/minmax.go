package utils

import "cmp"

func MinMax[T cmp.Ordered](minInt, n, maxInt T) T {
	return max(minInt, min(n, maxInt))
}
